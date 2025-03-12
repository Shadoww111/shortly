package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/shortly/internal/models"
)

type AuthService struct {
	db     *pgxpool.Pool
	secret string
}

func NewAuthService(db *pgxpool.Pool, secret string) *AuthService {
	return &AuthService{db: db, secret: secret}
}

func (s *AuthService) Register(ctx context.Context, req models.CreateUserRequest) (*models.TokenResponse, error) {
	// check existing
	var exists bool
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE email=$1 OR username=$2)",
		req.Email, req.Username,
	).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email or username already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = s.db.QueryRow(ctx,
		`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)
		 RETURNING id, username, email, is_active, created_at`,
		req.Username, req.Email, string(hash),
	).Scan(&user.ID, &user.Username, &user.Email, &user.IsActive, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		Token: token,
		User:  models.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email, IsActive: user.IsActive, CreatedAt: user.CreatedAt},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.TokenResponse, error) {
	var user models.User
	err := s.db.QueryRow(ctx,
		"SELECT id, username, email, password_hash, is_active, created_at FROM users WHERE email=$1",
		req.Email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsActive, &user.CreatedAt)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		Token: token,
		User:  models.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email, IsActive: user.IsActive, CreatedAt: user.CreatedAt},
	}, nil
}

func (s *AuthService) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}
