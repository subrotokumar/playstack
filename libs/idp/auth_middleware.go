package idp

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gitlab.com/subrotokumar/glitchr/libs/core"
)

type AuthMiddleware struct {
	region     string
	userPoolID string
	clientID   string
	issuer     string
	keyFunc    keyfunc.Keyfunc
}

func NewAuthMiddleware(region, userPoolID, clientId string) *AuthMiddleware {
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(2*time.Second),
	)
	defer cancel()
	jwksURL := fmt.Sprintf(
		"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
		region, userPoolID,
	)
	jwksKeyFunc, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		core.LogFatal("Failed to create a keyfunc.Keyfunc from the server's URL.", "err", err.Error())
	}
	return &AuthMiddleware{
		region:     region,
		userPoolID: userPoolID,
		clientID:   clientId,
		issuer: fmt.Sprintf(
			"https://cognito-idp.%s.amazonaws.com/%s",
			region, userPoolID,
		),
		keyFunc: jwksKeyFunc,
	}
}

type AuthResponse struct {
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
}

func (m *AuthMiddleware) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, echoHTTPError := m.extractJwtToken(c)
			if echoHTTPError != nil {
				return echoHTTPError
			}
			claims := token.Claims.(jwt.MapClaims)

			if claims["token_use"] != "access" {
				return echo.NewHTTPError(http.StatusUnauthorized, AuthResponse{Error: "not access token"})
			}
			if claims["iss"] != m.issuer {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid issuer")
			}
			if claims["client_id"] != m.clientID {
				return echo.NewHTTPError(http.StatusUnauthorized, AuthResponse{Error: "invalid client"})
			}
			c.Set("sub", uuid.MustParse(claims["sub"].(string)))
			return next(c)
		}
	}
}

func (m *AuthMiddleware) extractJwtToken(c echo.Context) (*jwt.Token, *echo.HTTPError) {
	var tokenStr string
	accessTokenCookie, err := c.Request().Cookie("access_token")
	if err != nil {
		auth := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, AuthResponse{Error: "missing token"})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		if tokenStr == "" {
			return nil, echo.NewHTTPError(http.StatusBadRequest, AuthResponse{Error: err.Error()})
		}
	}
	tokenStr = accessTokenCookie.Value
	token, err := jwt.Parse(tokenStr, m.keyFunc.Keyfunc)
	if err != nil || !token.Valid {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, AuthResponse{Error: "invalid token"})
	}
	return token, nil
}
