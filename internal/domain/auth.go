package domain

import "github.com/golang-jwt/jwt/v5"

// JWT Claims
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}
