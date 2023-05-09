package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/StarkovPO/Go-shop-final/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

func getLoyaltySystem(ctx context.Context, ID int, baseurl string) (models.Orders, error) {
	var order models.Orders
	client := &http.Client{}

	url := fmt.Sprintf("%s/api/orders/%v", baseurl, ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

	if err != nil {
		logrus.Warnf("error while creating the get request to the service: %v", err)
		return models.Orders{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Warnf("error while sending the get request to the service: %v", err)
		return models.Orders{}, err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Warnf("External service responce with code: %v and body: %v", resp.StatusCode, resp.Body)
		return models.Orders{}, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		logrus.Warnf("error while unmarshaling the request from service: %v", err)
		return models.Orders{}, err
	}
	return order, nil
}
