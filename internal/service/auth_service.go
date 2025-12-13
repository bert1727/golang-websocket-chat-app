package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/bert1727/ChatApp/internal/config"
	"github.com/bert1727/ChatApp/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

type AuthService interface {
	Login(email, password string) (*domain.LoginResponse, error)
	Register(req *domain.RegisterRequest) (*domain.LoginResponse, error)
	RefreshToken(refreshToken string) (string, error)
	VerifyToken(token string) (*domain.JWTClaims, error)
	GenerateTokens(user *domain.User) (accessToken, refreshToken string, err error)
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) GenerateTokens(user *domain.User) (string, string, error) {
	// Access Token
	accessClaims := domain.JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

func (s *authService) VerifyToken(tokenString string) (*domain.JWTClaims, error) {
	claims := &domain.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (s *authService) Login(email, password string) (*domain.LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		log.Info("invalid email")
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := CheckPassword(user.Password, password); err != nil {
		log.Info("invalid password")
		return nil, errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := s.GenerateTokens(user)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *authService) Register(req *domain.RegisterRequest) (*domain.LoginResponse, error) {
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    req.Email,
		Password: hashedPassword,
		Username: req.Username,
	}

	createdUser, err := s.userRepo.CreateNewUser(user)
	if err != nil {
		log.Info(err)
		return nil, err
	}

	log.Info("user was added")

	accessToken, refreshToken, err := s.GenerateTokens(createdUser)
	if err != nil {
		log.Info("Cannot generate refreshToken: ", err)
		return nil, err
	}

	log.Info("accessToken is generated")
	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *createdUser,
	}, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, error) {
	claims, err := s.VerifyToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Get user
	user, err := s.userRepo.FindUserByID(claims.UserID)
	if err != nil {
		return "", err
	}

	accessToken, _, err := s.GenerateTokens(user)
	return accessToken, err
}

func HashPassword(plain string) (string, error) {
	const cost = 12

	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), cost) // [web:41]
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func CheckPassword(hashed string, plain string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)); err != nil {
		return errors.New("invalid credentials")
	}
	return nil
}
