package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/StarkovPO/Go-shop-final/internal/appErrors"
	"github.com/StarkovPO/Go-shop-final/internal/config"
	"github.com/StarkovPO/Go-shop-final/internal/models"
	"github.com/StarkovPO/Go-shop-final/internal/store"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	salt       = "nalfhdp1238valls"
	signingKey = "qiausydigswig104#hlk[pzxn"
	tokenTTL   = 12 * time.Hour
)

type StoreInterface interface {
	CreateUserDB(ctx context.Context, user models.Users) error
	CheckLogin(ctx context.Context, login string) bool
	GetUserPass(ctx context.Context, login string) (string, bool)
}

type Service struct {
	store  store.Store
	config config.Config
}

type TokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

func NewService(s store.Store, c config.Config) Service {
	return Service{
		store:  s,
		config: c,
	}
}

func (s *Service) CreateUser(ctx context.Context, req models.Users) (string, error) {

	if exist := s.store.CheckLogin(ctx, req.Login); exist {
		return "", appErrors.ErrLoginAlreadyExist
	}

	req.Password = s.generatePasswordHash(req.Password)
	req.Id = generateUID()

	if err := s.store.CreateUserDB(ctx, req); err != nil {
		return "", appErrors.ErrCreateUser
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		req.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *Service) GenerateUserToken(ctx context.Context, req models.Users) (string, error) {

	passwordHash, exist := s.store.GetUserPass(ctx, req.Login)
	if !exist {
		return "", appErrors.ErrInvalidLoginOrPass
	}

	if isPassValid := s.comparePasswordHash(req.Password, passwordHash); isPassValid {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(tokenTTL).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			req.Id,
		})

		return token.SignedString([]byte(signingKey))
	}
	return "", appErrors.ErrInvalidLoginOrPass
}

func (s *Service) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *Service) comparePasswordHash(hashedPassword, password string) bool {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt))) == hashedPassword
}
