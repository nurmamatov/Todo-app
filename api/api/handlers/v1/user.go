package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"two_services/api/api/handlers/model"
	"two_services/api/etc"

	"two_services/api/api/token"
	"two_services/api/genproto/email"
	pbuser "two_services/api/genproto/user"
	l "two_services/api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

// RegistrUserAs godoc
// @Summary Registr new user
// @Description This method registr new user
// @Tags User
// @Accept json
// @Produce json
// @Param todo body user.CreateUserReq true "Registr User"
// @Success 201 {object} email.Status
// @Failure 400 {object} user.ErrOrStatus
// @Failure 404 {object} user.ErrOrStatus
// @Failure 500 {object} user.ErrOrStatus
// @Router /users [POST]
func (h *HandlerV1) RegistrUser(c *gin.Context) {
	var (
		body        pbuser.CreateUserReqWithCode
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true

	body.VerfCode = etc.GenerateCode(7)
	fmt.Println(body.VerfCode)
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.UserService().ChekUser(ctx, &pbuser.EmailWithUsername{Email: body.Email, Username: body.Username})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": res,
		})
	}
	_, err = h.serviceManager.UserService().Registr(ctx, &body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to registr task", l.Error(err))
		return
	}

	/*  // ourWiFi not work, send to phone
	res1, err := h.serviceManager.EmailService().SendSms(ctx, &email.Sms{
		Id:    "",
		Body:  body.VerfCode,
		Phone: "998337172004",
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Varify code cant't send to phone",
		})
	} else {
		c.JSON(http.StatusCreated, *res1)
	} */

	// this func sent to email
	_, err = h.serviceManager.EmailService().SendEmail(ctx, &email.Email{Subject: "Verify Code", Body: body.VerfCode, Email: body.Email})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Varify code cant't send to email",
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"send": "send",
		})
	}

}

// Verify godoc
// @Summary Verify new user
// @Description This method Verify new user
// @Tags User
// @Accept json
// @Produce json
// @Param todo body user.Check true "Verify User"
// @Success 201 {object} user.CreateUserReq
// @Failure 400 {object} user.ErrOrStatus
// @Failure 404 {object} user.ErrOrStatus
// @Failure 500 {object} user.ErrOrStatus
// @Router /users/verify [POST]
func (h *HandlerV1) Verify(c *gin.Context) {
	var (
		body        pbuser.Check
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.UserService().Verfy(ctx, &pbuser.Check{Username: body.Username, Code: body.Code})
	if err != nil {
		c.JSON(http.StatusConflict, model.ResponseError{
			Error: model.InternalServerError{
				Message: "Error while verify user",
			},
		})
		return
	}
	id, err := uuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusConflict, model.ResponseError{
			Error: model.InternalServerError{
				Message: "Error while generating new uuid for user",
			},
		})
		return
	}

	h.jwtHandler = token.JWTHandler{
		SigninKey: h.cfg.SigninKey,
		Sub:       id.String(),
		Iss:       "user",
		Role:      "authorized",
		Aud: []string{
			"khusniddin",
		},
		Log: h.log,
	}

	// Creating access and refresh tokens
	accessTokenString, refreshTokenString, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusConflict, model.ResponseError{
			Error: model.InternalServerError{
				Message: "Error while generating tokens",
			},
		})
		return
	}

	hashedPassword, err := etc.GeneratePasswordHash(res.Password)
	if err != nil {
		c.JSON(http.StatusConflict, model.ResponseError{
			Error: model.InternalServerError{
				Message: "Error while generating hash for password",
			},
		})
		return
	}

	checkEmail, err := h.serviceManager.UserService().ChekUser(
		context.Background(), &pbuser.EmailWithUsername{
			Email:    res.Email,
			Username: res.Username,
		},
	)
	if err != nil {
		c.JSON(http.StatusConflict, model.ResponseError{
			Error: model.InternalServerError{
				Message: "Error while checking field",
			},
		})
		return
	}

	if checkEmail.Chekfild {
		c.JSON(http.StatusConflict, model.ResponseError{
			Error: model.InternalServerError{
				Message: "User already exists",
			},
		})
		return
	}

	resUser, err := h.serviceManager.UserService().Create(context.Background(), &pbuser.CreateUserReqWithCode{
		FirstName:    res.FirstName,
		LastName:     res.LastName,
		Username:     res.Username,
		ProfilePhoto: res.ProfilePhoto,
		Bio:          res.Bio,
		Email:        res.Email,
		Gender:       res.Gender,
		Address:      res.Address,
		Phone:        res.Phone,
		Password:     string(hashedPassword),
		VerfCode:     body.Code,
		AccsesToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		Id:           id.String(),
	})
	if err != nil {
		c.JSON(http.StatusConflict, model.ResponseError{
			Error: model.InternalServerError{
				Message: "Error while create user",
			},
		})
		return
	} else {
		c.JSON(http.StatusCreated, resUser)
		return
	}
}

// GetTask godoc
// @Summary Get User
// @Description This method Get User
// @Security BearerAuth
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 201 {object} user.UserRes
// @Failure 400 {object} user.ErrOrStatus
// @Failure 404 {object} user.ErrOrStatus
// @Failure 500 {object} user.ErrOrStatus
// @Router /users/{id} [GET]
func (h *HandlerV1) GetUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	fmt.Println(guid)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().Get(ctx, &pbuser.GetOrDeleteUser{Id: guid})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get task", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUser godoc
// @Summary Update User
// @Description This method Update user
// @Security BearerAuth
// @Tags User
// @Accept json
// @Produce json
// @Param todo body user.UpdateUserReq true "Update User"
// @Success 201 {object} user.UserRes
// @Failure 400 {object} user.ErrOrStatus
// @Failure 404 {object} user.ErrOrStatus
// @Failure 500 {object} user.ErrOrStatus
// @Router /users [PUT]
func (h *HandlerV1) UpdateUser(c *gin.Context) {
	var (
		body        pbuser.UpdateUserReq
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Id = c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().Update(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUser godoc
// @Summary Delete User
// @Description This method Delete user
// @Security BearerAuth
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} user.ErrOrStatus
// @Failure 400 {object} user.ErrOrStatus
// @Failure 404 {object} user.ErrOrStatus
// @Failure 500 {object} user.ErrOrStatus
// @Router /users/{id} [PUT]
func (h *HandlerV1) DeleteUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().Delete(
		ctx, &pbuser.GetOrDeleteUser{
			Id: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete task", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Login godoc
// @Summary Login
// @Description This method Login
// @Tags User
// @Accept json
// @Produce json
// @Param todo body user.EmailWithPassword true "Login"
// @Success 200 {object} user.Mess
// @Failure 400 {object} user.ErrOrStatus
// @Failure 404 {object} user.ErrOrStatus
// @Failure 500 {object} user.ErrOrStatus
// @Router /users/login [POST]
func (h *HandlerV1) Login(c *gin.Context) {
	var (
		body        pbuser.EmailWithPassword
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().Login(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error(response.Username, l.Error(err))
		return
	}
	res, err := h.UpdateToken(response.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	response.AccesToken = res.AccessToken
	response.RefreshToken = res.RefreshToken
	c.JSON(http.StatusOK, response)
}

// Filtr godoc
// @Summary Filtr
// @Description This method Filtr
// @Tags User
// @Accept json
// @Produce json
// @Param todo body user.FiltrReq true "Filtr"
// @Success 200 {object} user.UsersList
// @Failure 400 {object} user.ErrOrStatus
// @Failure 404 {object} user.ErrOrStatus
// @Failure 500 {object} user.ErrOrStatus
// @Router /users/filtr [POST]
func (h *HandlerV1) Filtr(c *gin.Context) {
	var (
		body        pbuser.FiltrReq
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().Filtr(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error(response.String(), l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
