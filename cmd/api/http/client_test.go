package http

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

func TestClient(t *testing.T) {
	logger := log.New("debug", "test")

	config := ProductClientConfig{
		RequestURL: "https://catalog.wb.ru/catalog/%s/v2/catalog?ab_testing=false&appType=1&cat=8137&curr=rub&dest=-1257786&hide_dtype=13&lang=ru&page=%d&sort=popular&spp=30",
		ProductURL: "https://www.wildberries.ru/catalog/%d/detail.aspx",
		RetryCount: 3,
		RetryDelay: time.Second,
	}

	wbClient := NewProductClient(logger, config)
	ctx := context.Background()
	request := model.ProductsRequest{
		Page:     1,
		Category: model.ProductCategoryDresses,
	}

	for {
		fmt.Println(request.Page)

		products, err := wbClient.GetProducts(ctx, request)
		if err != nil {
			t.Error(err)
			break
		}

		if len(products) == 0 {
			fmt.Println("done")
			return
		}

		time.Sleep(time.Second)
		request.Page++
	}
}
