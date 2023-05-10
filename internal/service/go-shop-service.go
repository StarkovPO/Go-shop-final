package service

import (
	"context"
	"github.com/StarkovPO/Go-shop-final/internal/appErrors"
	"github.com/StarkovPO/Go-shop-final/internal/config"
	"github.com/StarkovPO/Go-shop-final/internal/models"
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
	CreateUserOrderDB(ctx context.Context, order models.Orders) error
	GetUserOrders(ctx context.Context, UID string) ([]models.Orders, error)
}

type Service struct {
	store  StoreInterface
	config config.Config
}

type TokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

func NewService(s StoreInterface, c config.Config) Service {
	return Service{
		store:  s,
		config: c,
	}
}

func (s *Service) CreateUser(ctx context.Context, req models.Users) (string, error) {

	if exist := s.store.CheckLogin(ctx, req.Login); exist { // remove checker and use DB index
		return "", appErrors.ErrLoginAlreadyExist
	}

	req.Password = generatePasswordHash(req.Password)
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

	isPassValid := comparePasswordHash(passwordHash, req.Password)
	if isPassValid {
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

func (s *Service) CreateUserOrder(ctx context.Context, req models.Orders) error {

	if !IsOrderNumberValid(req.ID) {
		return appErrors.ErrInvalidOrderNumber
	}
	/* it works only with external service */
	res, err := getLoyaltySystem(ctx, req.ID, s.config.AccrualSystemAddressValue)

	if err != nil {
		return err
	}
	res.UserID = req.UserID

	err = s.store.CreateUserOrderDB(ctx, res)

	//err := s.store.CreateUserOrderDB(ctx, req)

	return err
}

func (s *Service) GetUserOrders(ctx context.Context, UID string) ([]models.Orders, error) {

	req, err := s.store.GetUserOrders(ctx, UID)

	if err != nil {
		return nil, err
	}

	return req, nil
}
