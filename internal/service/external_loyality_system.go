package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StarkovPO/Go-shop-final/internal/models"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func getLoyaltySystem(ctx context.Context, ID int, baseurl string) (models.OrderFromService, error) {
	var order models.OrderFromService
	client := &http.Client{}

	url := fmt.Sprintf("%s/api/orders/%v", baseurl, ID)
	logrus.Printf("request to the service: %v", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

	logrus.Printf("Ask getLoyaltySystem")
	if err != nil {
		logrus.Errorf("error while creating the get request to the service: %v", err)
		return models.OrderFromService{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Printf("ops something went wrong: %v", err)
		logrus.Errorf("error while sending the get request to the service: %v", err)
		return models.OrderFromService{}, err
	}

	// {"level":"error","msg":"External service responce with code: 200 and body: {\"order\":\"40534687\",\"status\":\"PROCESSED\",\"accrual\":729.98}\n","time":"2023-05-14T07:20:11Z"}
	if resp.StatusCode == http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		logrus.Printf("External service responce with code: 204 and body: %v", string(b))
		return models.OrderFromService{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		logrus.Errorf("External service responce with code: %v and body: %v", resp.StatusCode, string(b))
		return models.OrderFromService{}, errors.New("external service response with bad status code")
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		b, _ := io.ReadAll(resp.Body)
		logrus.Printf("ops something went wrong: %v", err)
		logrus.Errorf("error while unmarshaling the request from service: %v and body: %v", err, string(b))
		return models.OrderFromService{}, err
	}

	logrus.Printf("response succesfull unmarshaled.ID: %v, status: %v, accural: %v", order.ID, order.Status, order.Accrual)
	return order, nil
}
