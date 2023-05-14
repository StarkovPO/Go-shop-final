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

	logrus.Printf("Ask getLoyaltySystem")
	if err != nil {
		logrus.Errorf("error while creating the get request to the service: %v", err)
		return models.Orders{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Printf("ops something went wrong: %v", err)
		logrus.Errorf("error while sending the get request to the service: %v", err)
		return models.Orders{}, err
	}

	// {"level":"error","msg":"External service responce with code: 200 and body: \u0026{0xc0000a85c0 {0 0} false \u003cnil\u003e 0x748800 0x748900}","time":"2023-05-14T07:11:51Z"}
	//if resp.StatusCode != http.StatusOK|http.StatusNoContent {
	//	b, _ := io.ReadAll(resp.Body)
	//	logrus.Printf("ops something went wrong. Http code: %v", resp.StatusCode)
	//	logrus.Errorf("External service responce with code: %v and body: %v", resp.StatusCode, string(b))
	//	return models.Orders{}, err
	//}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		logrus.Printf("ops something went wrong: %v", err)
		logrus.Errorf("error while unmarshaling the request from service: %v", err)
		return models.Orders{}, err
	}
	return order, nil
}
