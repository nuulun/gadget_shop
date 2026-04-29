package service

import (
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidToken = errors.New("invalid token")

type Repo interface {
	CreateUser(context.Context, model.User) error
	GetUserByLoginOrEmail(context.Context, string) (model.User, error)
	GetUserByID(context.Context, uint64) (model.User, error)
	SaveRefreshToken(context.Context, model.RefreshToken) error
	GetRefreshToken(context.Context, string) (model.RefreshToken, error)
	RevokeRefreshToken(context.Context, string) error
	DeleteUser(context.Context, uint64) error
}

type AuthService struct {
	repo              Repo
	jwtSecret         string
	accessTTLMinutes  int
	refreshTTLDays    int
}

func New(repo Repo, jwtSecret string, accessTTLMinutes, refreshTTLDays int) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret, accessTTLMinutes: accessTTLMinutes, refreshTTLDays: refreshTTLDays}
}

func (s *AuthService) Register(ctx context.Context, in model.RegisterInput) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	return s.repo.CreateUser(ctx, model.User{ID: in.ID, Login: in.Login, Email: in.Email, PasswordHash: string(hash), CreatedAt: time.Now(), UpdatedAt: time.Now()})
}

func (s *AuthService) Login(ctx context.Context, in model.LoginInput) (model.TokenPair, error) {
	user, err := s.repo.GetUserByLoginOrEmail(ctx, in.LoginOrEmail)
	if errors.Is(err, repository.ErrNotFound) {
		return model.TokenPair{}, ErrInvalidCredentials
	}
	if err != nil {
		return model.TokenPair{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return model.TokenPair{}, ErrInvalidCredentials
	}
	return s.issueTokens(ctx, user.ID)
}

func (s *AuthService) Validate(ctx context.Context, token string) (uint64, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return 0, ErrInvalidToken
	}
	uid, ok := claims["sub"].(float64)
	if !ok {
		return 0, ErrInvalidToken
	}
	return uint64(uid), nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (model.TokenPair, error) {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if errors.Is(err, repository.ErrNotFound) {
		return model.TokenPair{}, ErrInvalidToken
	}
	if err != nil {
		return model.TokenPair{}, err
	}
	if rt.RevokedAt != nil || time.Now().After(rt.ExpiresAt) {
		return model.TokenPair{}, fmt.Errorf("refresh token expired or revoked")
	}
	return s.issueTokens(ctx, rt.UserID)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.repo.RevokeRefreshToken(ctx, refreshToken)
}

func (s *AuthService) DeleteUser(ctx context.Context, id uint64) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *AuthService) issueTokens(ctx context.Context, userID uint64) (model.TokenPair, error) {
	now := time.Now()
	accessExp := now.Add(time.Duration(s.accessTTLMinutes) * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": userID, "exp": accessExp.Unix()})
	accessStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return model.TokenPair{}, err
	}

	raw := fmt.Sprintf("%d:%d:%s", userID, now.UnixNano(), s.jwtSecret)
	h := sha256.Sum256([]byte(raw))
	refreshStr := hex.EncodeToString(h[:])

	if err := s.repo.SaveRefreshToken(ctx, model.RefreshToken{
		UserID:    userID,
		Token:     refreshStr,
		ExpiresAt: now.Add(time.Duration(s.refreshTTLDays) * 24 * time.Hour),
		CreatedAt: now,
	}); err != nil {
		return model.TokenPair{}, err
	}
	return model.TokenPair{AccessToken: accessStr, RefreshToken: refreshStr}, nil
}
