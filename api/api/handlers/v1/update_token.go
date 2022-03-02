package v1

import (
	"context"
	"time"
	pbuser "two_services/api/genproto/user"
	"two_services/api/api/token"
)

func (h *HandlerV1) UpdateToken(id string) (*pbuser.Tokens, error) {
	h.jwtHandler = token.JWTHandler{
		SigninKey: h.cfg.SigninKey,
		Sub:       id,
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
		return nil, err
	}
	newToken := pbuser.TokensReq{
		Id:           id,
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.UserService().UpdateToken(ctx, &newToken)
	if err != nil {
		return nil, err
	}
	return res, nil
}
