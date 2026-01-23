package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/bert1727/ChatApp/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/bert1727/ChatApp/internal/config"
	"github.com/bert1727/ChatApp/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
	valid    *validator.Validate
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
		valid:    validator.New(),
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
	refreshClaims := domain.JWTClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

func (s *authService) VerifyToken(refreshToken string) (*domain.JWTClaims, error) {
	claims := &domain.JWTClaims{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Info().Msg("unexpected signing method")
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
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		log.Info().Msg("invalid email")
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := CheckPassword(user.Password, password); err != nil {
		fmt.Println("Пароли:\n", user.Password, password)
		log.Info().Msg("invalid password")
		return nil, fmt.Errorf("invalid credentials")
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
	// Validation
	err := s.valid.Struct(*req)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				fmt.Println(e.Namespace())
			}
		}
		return nil, err
		// FIX: change return types
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    req.Email,
		Password: hashedPassword,
		Username: req.Username,
	}

	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		log.Err(err).Msg("failed to create a user")
		return nil, err
	}

	log.Info().Msg("user was added")

	accessToken, refreshToken, err := s.GenerateTokens(createdUser)
	if err != nil {
		log.Info().Err(err).Msg("Cannot generate refreshToken")
		return nil, err
	}

	log.Info().Msg("accessToken is generated")
	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *createdUser,
	}, nil
}

// Update accessToken via refresh token
func (s *authService) RefreshToken(refreshToken string) (string, error) {
	claims, err := s.VerifyToken(refreshToken)
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", err
	}

	accessToken, _, err := s.GenerateTokens(user)
	log.Info().Uint("user id", user.ID).Msg("accessToken generated for user")
	return accessToken, err
}

func HashPassword(plain string) (string, error) {
	const cost = 10

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
