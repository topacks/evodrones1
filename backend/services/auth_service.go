package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/root/evodrones/backend/config"
	"github.com/root/evodrones/backend/models"
	"github.com/root/evodrones/backend/repositories"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	userRepo *repositories.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo *repositories.UserRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *AuthService) Register(input models.CreateUserInput) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Create new user
	user, err := models.NewUser(input)
	if err != nil {
		return nil, err
	}

	// Save user to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(input models.LoginInput) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check password
	if err := user.CheckPassword(input.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(s.config.JWT.AccessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.SecretKey))
}

func (s *AuthService) generateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(s.config.JWT.RefreshTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.SecretKey))
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWT.SecretKey), nil
	})
}

func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	return s.userRepo.GetByID(userID)
}
