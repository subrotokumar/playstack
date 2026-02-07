package server

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gitlab.com/subrotokumar/playstack/libs/db"
)

const (
	ErrInvalidToken = "invalid token"
)

type (
	Profile struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Sub   string `json:"sub"`
	}

	AuthResponse struct {
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
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			name		formData	string	true	"Name"
//	@Param			email		formData	string	true	"Email"
//	@Param			password	formData	string	true	"Password"
//	@Success		200		{object}	AuthResponse
//	@Failure		400		{object} 	AuthResponse
//	@Failure		500		{object}	AuthResponse
//	@Router			/auth/users [post]
func (s *Server) SignupHandler(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if err := validator.Var(name, "required"); err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: err.Error()})
	}
	if err := validator.Var(email, "required,email"); err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: err.Error()})
	}
	if err := validator.Var(password, "required,min=8"); err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: "password must be at least 8 characters"})
	}
	if err := validator.Var(password, "containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"); err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: "password must contain at least 1 uppercase letter"})
	}
	if err := validator.Var(password, "containsany=abcdefghijklmnopqrstuvwxyz"); err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: "password must contain at least 1 lowercase letter"})
	}
	if err := validator.Var(password, "containsany=0123456789"); err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: "password must contain at least 1 number"})
	}

	_, userSub, err := s.idp.SignUp(
		c.Request().Context(),
		name,
		email,
		password,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthResponse{
			Message: "failed to signup",
			Error:   err.Error(),
		})
	}

	s.store.CreateUser(c.Request().Context(), db.CreateUserParams{
		ID:    uuid.MustParse(userSub),
		Email: email,
		Name:  name,
	})

	return c.JSON(http.StatusOK, AuthResponse{
		Message: "User signed up successfully. Please confirm your email.",
	})
}

// LoginHandler godoc
//
//	@Summary		Login user
//	@Description	Authenticate a user and set access/refresh cookies
//	@Tags			auth
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			email		formData	string	true	"User email"
//	@Param			password	formData	string	true	"User password"
//	@Success		200		{object}	AuthResponse
//	@Failure		400		{object}	AuthResponse
//	@Failure		401		{object}	AuthResponse
//	@Router			/auth/sessions [post]
func (s *Server) LoginHandler(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if err := validator.Var(email, "required,email"); err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: err.Error()})
	}
	if err := validator.Var(password, "required"); err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: err.Error()})
	}

	tokens, err := s.idp.Login(c.Request().Context(), email, password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, AuthResponse{Error: err.Error()})
	}
	_ = tokens

	c.SetCookie(&http.Cookie{Name: "access_token", Value: tokens.AccessToken, HttpOnly: true, Secure: true, Path: "/"})
	c.SetCookie(&http.Cookie{Name: "refresh_token", Value: tokens.RefreshToken, HttpOnly: true, Secure: true, Path: "/"})
	c.SetCookie(&http.Cookie{Name: "id_token", Value: tokens.IdToken, HttpOnly: true, Secure: true, Path: "/"})

	return c.JSON(http.StatusOK, AuthResponse{Message: "User logged in successfully"})
}

// ResentOTP godoc
//
//	@Summary		Resend OTP
//	@Description	Resend confirmation code (OTP) to user's email
//	@Tags			auth
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			email	formData	string	true	"User email"
//	@Success		200		{object}	AuthResponse
//	@Failure		400		{object}	AuthResponse
//	@Failure		500		{object}	AuthResponse
//	@Router			/auth/verifications [post]
func (s *Server) ResentOTP(c echo.Context) error {
	email := c.FormValue("email")
	err := validator.Var(email, "required,email")
	if err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: err.Error()})
	}

	if err := s.idp.ResendOTP(c.Request().Context(), email); err != nil {
		return c.JSON(http.StatusInternalServerError, AuthResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, AuthResponse{Message: "OTP send successfully"})
}

// ConfirmSignupHandler godoc
//
//	@Summary		Confirm signup
//	@Description	Confirm a user's signup using OTP
//	@Tags			auth
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			email	formData	string	true	"User email"
//	@Param			otp		formData	string	true	"OTP code"
//	@Success		200		{object}	AuthResponse
//	@Failure		400		{object}	AuthResponse
//	@Failure		500		{object}	AuthResponse
//	@Router			/auth/verifications/confirm [post]
func (s *Server) ConfirmSignupHandler(c echo.Context) error {
	email := c.FormValue("email")
	otp := c.FormValue("otp")

	if err := validator.Var(email, "required,email"); err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: err.Error()})
	}
	if err := validator.Var(otp, "required"); err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: err.Error()})
	}

	if err := s.idp.ConfirmSignUp(c.Request().Context(), email, otp); err != nil {
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
//	@Router			/auth/tokens [post]
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
//	@Success		200	{object}	AuthResponse
//	@Failure		400	{object}	AuthResponse
//	@Failure		500	{object}	AuthResponse
//	@Router			/auth/profile [post]
func (s *Server) ProfileHandler(c echo.Context) error {
	idTokenCookie, err := c.Request().Cookie("id_token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, AuthResponse{Error: err.Error()})
	}

	claims := jwt.MapClaims{}
	_, _, err = new(jwt.Parser).ParseUnverified(idTokenCookie.Value, claims)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, AuthResponse{Error: ErrInvalidToken})
	}

	sub, ok := claims["username"].(string)
	if !ok {
		sub = claims["sub"].(string)
	}

	email, ok := claims["email"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, AuthResponse{Error: ErrInvalidToken})
	}
	name, ok := claims["name"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, AuthResponse{Error: ErrInvalidToken})
	}
	return c.JSON(http.StatusOK, AuthResponse{
		Data: &Profile{
			Name:  name,
			Email: email,
			Sub:   sub,
		},
	})
}
