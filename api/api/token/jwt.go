package token

import (
	"fmt"
	"time"
	"two_services/api/pkg/logger"

	jwt"github.com/dgrijalva/jwt-go"
)

type JWTHandler struct {
	Sub       string
	Iss       string
	Exp       string
	Iat       string
	Aud       []string
	Role      string
	SigninKey string
	Log       logger.Logger
	Token     string
}

func (jwtHandler *JWTHandler) GenerateAuthJWT() (access, refresh string, err error) {
	var (
		accessToken  *jwt.Token
		refreshToken *jwt.Token
		claims       jwt.MapClaims
	)
	accessToken = jwt.New(jwt.SigningMethodHS256)
	refreshToken = jwt.New(jwt.SigningMethodHS256)

	claims = accessToken.Claims.(jwt.MapClaims)
	claims["iss"] = jwtHandler.Iss
	claims["sub"] = jwtHandler.Sub
	claims["exp"] = time.Now().Add(time.Hour * 500).Unix()
	claims["iat"] = time.Now().Unix()
	claims["role"] = jwtHandler.Role
	claims["aud"] = jwtHandler.Aud
	fmt.Println(jwtHandler.SigninKey)
	access, err = accessToken.SignedString([]byte(jwtHandler.SigninKey))
	if err != nil {
		jwtHandler.Log.Error("error generate access token", logger.Error(err))
		return
	}
	refresh, err = refreshToken.SignedString([]byte(jwtHandler.SigninKey))
	if err != nil {
		jwtHandler.Log.Error("error generate refresh token", logger.Error(err))
		return
	}

	return

}

func (jwtHandler *JWTHandler) ExtractClaims() (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = jwt.Parse(jwtHandler.Token, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtHandler.SigninKey), nil
	})
	if err != nil {
		fmt.Println(err, "err1")
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		jwtHandler.Log.Error("invalid jwt token")
		return nil, err
	}
	return claims, err
}

func ExtractClaim(tokenStr string, signingKey []byte) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)
	token, err = jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, err
	}
	return claims, nil
}
