package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	SignUpRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	SignUpResponse struct {
		Message string `json:"message"`
		UserSub string `json:"user_sub,omitempty"`
	}

	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		Message string `json:"message"`
	}

	ConfirmSignupRequest struct {
		Email string `json:"email"`
		Otp   string `json:"otp"`
	}

	ConfirmSignupResponse struct {
		Message string `json:"message"`
	}

	RefreshTokenRequest struct {
		Email        string `json:"email"`
		RefreshToken string `json:"refresh_token"`
	}

	RefreshTokenResponse struct {
		Message string `json:"message"`
	}
)

func (s *Server) SignupHandler(c echo.Context) error {
	var body SignUpRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseEntity{Error: err.Error()})
	}
	if err := validator.Struct(&body); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseEntity{Error: err.Error()})
	}

	confirmed, userSub, err := s.idp.SignUp(
		c.Request().Context(),
		body.Name,
		body.Email,
		body.Password,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseEntity{Error: err.Error()})
	}

	if confirmed {
		return c.JSON(http.StatusOK, ResponseEntity{Data: SignUpResponse{
			Message: "User signed up successfully",
			UserSub: userSub,
		}})
	}

	return c.JSON(http.StatusOK, ResponseEntity{Data: SignUpResponse{
		Message: "User signed up successfully. Please confirm your email.",
	}})
}

func (s *Server) LoginHandler(c echo.Context) error {
	var body LoginRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseEntity{Error: err.Error()})
	}
	if err := validator.Struct(&body); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseEntity{Error: err.Error()})
	}

	tokens, err := s.idp.Login(c.Request().Context(), body.Email, body.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseEntity{Error: err.Error()})
	}
	_ = tokens

	c.SetCookie(&http.Cookie{Name: "access_token", Value: tokens.AccessToken, HttpOnly: true, Secure: true, Path: "/"})
	c.SetCookie(&http.Cookie{Name: "refresh_token", Value: tokens.RefreshToken, HttpOnly: true, Secure: true, Path: "/"})

	return c.JSON(http.StatusOK, ResponseEntity{Data: LoginResponse{Message: "User logged in successfully"}})
}

func (s *Server) ConfirmSignupHandler(c echo.Context) error {
	var body ConfirmSignupRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseEntity{Error: err.Error()})
	}
	if err := validator.Struct(&body); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseEntity{Error: err.Error()})
	}

	if err := s.idp.ConfirmSignUp(c.Request().Context(), body.Email, body.Otp); err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseEntity{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, ResponseEntity{Data: ConfirmSignupResponse{Message: "User confirmed successfully"}})
}

func (s *Server) RefreshTokenHandler(c echo.Context) error {
	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseEntity{Error: err.Error()})
	}

	emailCookie, err := c.Request().Cookie("email")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseEntity{Error: err.Error()})
	}

	accessToken, err := s.idp.RefreshAccessToken(
		c.Request().Context(),
		emailCookie.Value,
		refreshCookie.Value,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseEntity{Error: err.Error()})
	}

	c.SetCookie(&http.Cookie{Name: "access_token", Value: accessToken, HttpOnly: true, Secure: true, Path: "/"})

	return c.JSON(http.StatusOK, ResponseEntity{Data: RefreshTokenResponse{Message: "Access token refreshed"}})
}
