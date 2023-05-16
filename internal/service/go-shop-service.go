package service

import (
	"context"
	"errors"
	"github.com/StarkovPO/Go-shop-final/internal/appErrors"
	"github.com/StarkovPO/Go-shop-final/internal/config"
	"github.com/StarkovPO/Go-shop-final/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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
	CreateWithdraw(ctx context.Context, req models.Withdrawn) error
	GetUserWithdrawnDB(ctx context.Context, UID string) ([]models.Withdrawn, error)
	UpdateOrderStatus(ctx context.Context, status string, UID string) error
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
		return err
	}

	if res.Status == "" || res.Status == "REGISTERED" || res.Status == "PROCESSING" {
		statusChan := make(chan int)
		orderID, err := strconv.Atoi(res.ID)
		if err != nil {
			return err
		}
		statusChan <- orderID

		go s.updater(ctx, statusChan, s.config.AccrualSystemAddressValue)
	}

	res.UserID = req.UserID
	res.ID = req.ID
	err = s.store.CreateUserOrderDB(ctx, res)

	if err != nil {
		return err
	}

	err = s.store.IncreaseUserBalance(ctx, res.Accrual, res.UserID)

	if err != nil {

		return err
	}

	return nil
}

func (s *Service) updater(ctx context.Context, orderIDChan chan int, baseurl string) {

	done := make(chan bool)
	g, _ := errgroup.WithContext(ctx)
	logrus.Info("updater start")
	for id := range orderIDChan {
		id := id
		g.Go(func() error {
			time.Sleep(5 * time.Second)
			logrus.Info("gorutine sending request")
			res, err := getLoyaltySystem(ctx, id, baseurl)
			if err != nil {
				logrus.Errorf("fatal fetching new order status: %v", err)
				return err
			}
			if res.Status == "INVALID" || res.Status == "PROCESSED" {
				logrus.Info("gorutine starting update status")
				err = s.store.UpdateOrderStatus(ctx, res.Status, res.ID)
				if err != nil {
					logrus.Errorf("fatal updating Order status: %v", err)
					return err
				}
				logrus.Info("gorutine updated order status")
			} else {
				logrus.Info("gorutine got not final status")
				orderIDChan <- id
			}
			done <- true
			return nil
		})
	}

	go func() {
		if err := g.Wait(); err != nil {
			logrus.Errorf("unexpected error: %v", err)
		}
		done <- true
	}()
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
		logrus.Errorf("ops unhandled error on service level: %v", err)
		return models.Balance{}, err
	}

	return b, nil
}

func (s *Service) CreateUserWithdraw(ctx context.Context, req models.Withdrawn) error {

	intId, _ := strconv.Atoi(req.OrderID)
	if !IsOrderNumberValid(intId) {
		return appErrors.ErrInvalidOrderNumber
	}

	b, err := s.store.GetUserBalanceDB(ctx, req.UserID)

	if err != nil {
		logrus.Errorf("ops unhandled error on service level: %v", err)
		return err
	}

	if b.Current < req.Withdrawn {
		return appErrors.ErrNotEnoughPoints
	}

	err = s.store.CreateWithdraw(ctx, req)

	return nil
}

func (s *Service) GetUserWithdrawn(ctx context.Context, UID string) ([]models.Withdrawn, error) {
	res, err := s.store.GetUserWithdrawnDB(ctx, UID)

	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, err
	}

	return res, appErrors.ErrWithdrawnNotFound
}
