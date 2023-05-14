package service

import (
	"context"
	"errors"
	"github.com/StarkovPO/Go-shop-final/internal/appErrors"
	"github.com/StarkovPO/Go-shop-final/internal/config"
	"github.com/StarkovPO/Go-shop-final/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"strconv"
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
	CreateUserOrderDB(ctx context.Context, order models.OrderFromService) error
	GetUserOrders(ctx context.Context, UID string) ([]models.Orders, error)
	GetUserID(ctx context.Context, login string) (string, error)
	IncreaseUserBalance(ctx context.Context, accrual float64, UID string) error
	GetUserBalanceDB(ctx context.Context, UID string) (models.Balance, error)
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
		UID, err := s.store.GetUserID(ctx, req.Login)
		if err != nil {
			return "", errors.New("error while getting UID: %v")
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(tokenTTL).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			UID,
		})

		return token.SignedString([]byte(signingKey))
	}
	return "", appErrors.ErrInvalidLoginOrPass
}

func (s *Service) CreateUserOrder(ctx context.Context, req models.Orders) error {

	intId, _ := strconv.Atoi(req.ID)
	if !IsOrderNumberValid(intId) {
		return appErrors.ErrInvalidOrderNumber
	}
	/* it works only with external service */
	res, err := getLoyaltySystem(ctx, intId, s.config.AccrualSystemAddressValue)

	if err != nil {
		logrus.Printf("ops something went wrong: %v", err)
		return err
	}
	res.UserID = req.UserID
	logrus.Printf("calling createUserOrderDB")
	err = s.store.CreateUserOrderDB(ctx, res)

	if res.Accrual != 0 {
		logrus.Printf("accural: %v", res.Accrual)
		err = s.store.IncreaseUserBalance(ctx, res.Accrual, res.UserID)
	}

	//err := s.store.CreateUserOrderDB(ctx, req)

	return err
}

func (s *Service) GetUserOrders(ctx context.Context, UID string) ([]models.Orders, error) {

	req, err := s.store.GetUserOrders(ctx, UID)

	if err != nil {
		return nil, err
	}

	if req != nil {
		return req, nil
	}

	return req, appErrors.ErrOrderNotFound
}

func (s *Service) GetUserBalance(ctx context.Context, UID string) (models.Balance, error) {
	b, err := s.store.GetUserBalanceDB(ctx, UID)

	if err != nil {
		return models.Balance{}, err
	}

	return b, nil
}
