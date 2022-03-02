package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"two_services/api/api/models"
	"two_services/api/api/token"
	"two_services/api/config"

	"github.com/casbin/casbin/v2"
	jwtg "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JWTRoleAuthorizer struct {
	enforcer   *casbin.Enforcer
	cfg        config.Config
	jwtHandler token.JWTHandler
}

func NewAuthorizer(e *casbin.Enforcer, jwtHandler token.JWTHandler, cfg config.Config) gin.HandlerFunc {
	a := &JWTRoleAuthorizer{
		enforcer:   e,
		cfg:        cfg,
		jwtHandler: jwtHandler,
	}
	return func(c *gin.Context) {
		allow, err := a.ChekPermission(c.Request)
		if err != nil {
			v, _ := err.(*jwtg.ValidationError)
			if v.Errors == jwtg.ValidationErrorExpired {
				a.RequireRefresh(c)
			} else {
				a.RequirePermission(c)
			}
		} else if !allow {
			fmt.Println("salom222")
			a.RequirePermission(c)
		}
	}
}

func (a *JWTRoleAuthorizer) GetRole(r *http.Request) (string, error) {
	var (
		role   string
		claims jwtg.MapClaims
		err    error
	)

	jwtToken := r.Header.Get("Authorization")
	if jwtToken == "" {
		return "unauthorized", nil
	} else if strings.Contains(jwtToken, "Basic") {
		return "unauthorized", nil
	}

	a.jwtHandler.Token = jwtToken
	claims, err = a.jwtHandler.ExtractClaims()
	if err != nil {
		return "", err
	}

	if claims["role"].(string) == "authorized" {
		role = "authorized"
	} else {
		role = "unknown"
	}
	return role, nil
}

func (a *JWTRoleAuthorizer) ChekPermission(r *http.Request) (bool, error) {
	user, err := a.GetRole(r)
	if err != nil {
		return false, err
	}
	method := r.Method
	path := r.URL.Path
	fmt.Println(method,path,user)
	allowed, err := a.enforcer.Enforce(user, path, method)
	if err != nil {
		panic(err)
	}
	return allowed, nil
}

func (a *JWTRoleAuthorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(403)
}
func (a *JWTRoleAuthorizer) RequireRefresh(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, models.ResponseError{
		Error: models.ServerError{
			Status:  "UNAUTHORIZED",
			Message: "Token is expired",
		},
	})
	c.AbortWithStatus(401)
}
