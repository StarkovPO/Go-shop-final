package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/StarkovPO/Go-shop-final/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

func getLoyaltySystem(ctx context.Context, ID int, baseurl string) (models.OrderFromService, error) {
	var order models.OrderFromService
	client := &http.Client{}

	url := fmt.Sprintf("%s/api/orders/%v", baseurl, ID)
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
	//if resp.StatusCode != http.StatusOK && http.StatusNoContent {
	//	b, _ := io.ReadAll(resp.Body)
	//	logrus.Printf("ops something went wrong. Http code: %v", resp.StatusCode)
	//	logrus.Errorf("External service responce with code: %v and body: %v", resp.StatusCode, string(b))
	//	return models.Orders{}, err
	//}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		logrus.Printf("ops something went wrong: %v", err)
		logrus.Errorf("error while unmarshaling the request from service: %v and body: ", err)
		return models.OrderFromService{}, err
	}
	return order, nil
}
