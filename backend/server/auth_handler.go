package server

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gitlab.com/subrotokumar/glitchr/libs/db"
)

const (
	ErrInvalidToken = "invalid token"
)

type (
	SignUpRequest struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"required"`
	}

	ConfirmSignupRequest struct {
		Email string `json:"email" validate:"email,required"`
		Otp   string `json:"otp" validate:"required"`
	}

	RefreshTokenRequest struct {
		Email        string `json:"email" validate:"email,required"`
		RefreshToken string `json:"refresh_token"`
	}

	AuthResponse struct {
		Message string `json:"message,omitempty"`
		Error   any    `json:"error,omitempty"`
	}

	Profile struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Sub   string `json:"sub"`
	}

	ProfileResponse struct {
		Data    *Profile `json:"data,omitempty"`
		Message string   `json:"message,omitempty"`
		Error   any      `json:"error,omitempty"`
	}
)

// SignupHandler godoc
//
//	@Summary		Sign up a new user
//	@Description	Create a new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		SignUpRequest	true	"Sign up payload"
//	@Success		200		{object}	AuthResponse
//	@Failure		400		{object} 	AuthResponse
//	@Failure		500		{object}	AuthResponse
//	@Router			/auth/signup [post]
func (s *Server) SignupHandler(c echo.Context) error {
	var body SignUpRequest
	if err := RequestBody(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: err.Error()})
	}

	confirmed, userSub, err := s.idp.SignUp(
		c.Request().Context(),
		body.Name,
		body.Email,
		body.Password,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthResponse{Error: err.Error()})
	}

	if confirmed {
		s.store.CreateUser(c.Request().Context(), db.CreateUserParams{
			ID:    uuid.MustParse(userSub),
			Email: body.Email,
		})
	}
	return c.JSON(http.StatusOK, AuthResponse{
		Message: "User signed up successfully. Please confirm your email.",
	})
}

// LoginHandler godoc
//
//	@Summary		Login user
//	@Description	Authenticate a user and set access/refresh cookies
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginRequest	true	"Login payload"
//	@Success		200		{object}	AuthResponse
//	@Failure		400		{object}	AuthResponse
//	@Failure		401		{object}	AuthResponse
//	@Router			/auth/login [post]
func (s *Server) LoginHandler(c echo.Context) error {
	var body LoginRequest
	if err := RequestBody(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: err.Error()})
	}

	tokens, err := s.idp.Login(c.Request().Context(), body.Email, body.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, AuthResponse{Error: err.Error()})
	}
	_ = tokens

	c.SetCookie(&http.Cookie{Name: "access_token", Value: tokens.AccessToken, HttpOnly: true, Secure: true, Path: "/"})
	c.SetCookie(&http.Cookie{Name: "refresh_token", Value: tokens.RefreshToken, HttpOnly: true, Secure: true, Path: "/"})
	c.SetCookie(&http.Cookie{Name: "id_token", Value: tokens.IdToken, HttpOnly: true, Secure: true, Path: "/"})

	return c.JSON(http.StatusOK, AuthResponse{Message: "User logged in successfully"})
}

// ConfirmSignupHandler godoc
//
//	@Summary		Confirm signup
//	@Description	Confirm a user's signup using OTP
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		ConfirmSignupRequest	true	"Confirm signup payload"
//	@Success		200		{object}	AuthResponse
//	@Failure		400		{object}	AuthResponse
//	@Failure		500		{object}	AuthResponse
//	@Router			/auth/confirm-signup [post]
func (s *Server) ConfirmSignupHandler(c echo.Context) error {
	var body ConfirmSignupRequest
	if err := RequestBody(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: err.Error()})
	}

	if err := s.idp.ConfirmSignUp(c.Request().Context(), body.Email, body.Otp); err != nil {
		return c.JSON(http.StatusInternalServerError, AuthResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, AuthResponse{Message: "User confirmed successfully"})
}

// RefreshTokenHandler godoc
//
//	@Summary		Refresh access token
//	@Description	Refresh access token using refresh token cookie
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	AuthResponse
//	@Failure		400	{object}	AuthResponse
//	@Failure		500	{object}	AuthResponse
//	@Router			/auth/refresh [post]
func (s *Server) RefreshTokenHandler(c echo.Context) error {
	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: err.Error()})
	}

	idTokenCookie, err := c.Request().Cookie("id_token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: err.Error()})
	}

	claims := jwt.MapClaims{}
	_, _, err = new(jwt.Parser).ParseUnverified(idTokenCookie.Value, claims)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, AuthResponse{Error: ErrInvalidToken})
	}

	username, ok := claims["username"].(string)
	if !ok {
		username, ok = claims["sub"].(string) // fallback
	}

	if !ok || username == "" {
		return c.JSON(http.StatusUnauthorized, AuthResponse{Error: "user identity missing"})
	}

	accessToken, err := s.idp.RefreshAccessToken(
		c.Request().Context(),
		username,
		refreshCookie.Value,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthResponse{Error: err.Error()})
	}

	c.SetCookie(&http.Cookie{Name: "access_token", Value: accessToken, HttpOnly: true, Secure: true, Path: "/"})

	return c.JSON(http.StatusOK, AuthResponse{Message: "Access token refreshed"})
}

// ProfileHandler godoc
//
//	@Summary		Profile
//	@Description	Get Profile Detail
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ProfileResponse
//	@Failure		400	{object}	ProfileResponse
//	@Failure		500	{object}	ProfileResponse
//	@Router			/auth/profile [post]
func (s *Server) ProfileHandler(c echo.Context) error {
	idTokenCookie, err := c.Request().Cookie("id_token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ProfileResponse{Error: err.Error()})
	}

	claims := jwt.MapClaims{}
	_, _, err = new(jwt.Parser).ParseUnverified(idTokenCookie.Value, claims)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ProfileResponse{Error: ErrInvalidToken})
	}

	sub, ok := claims["username"].(string)
	if !ok {
		sub = claims["sub"].(string)
	}

	email, ok := claims["email"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ProfileResponse{Error: ErrInvalidToken})
	}
	name, ok := claims["name"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ProfileResponse{Error: ErrInvalidToken})
	}
	return c.JSON(http.StatusOK, ProfileResponse{
		Data: &Profile{
			Name:  name,
			Email: email,
			Sub:   sub,
		},
	})
}
