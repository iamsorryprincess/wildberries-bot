package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type ProductClientConfig struct {
	RequestURL string `config:"request_url"`
	ProductURL string `config:"product_url"`

	RetryCount uint          `config:"retry_count"`
	RetryDelay time.Duration `config:"retry_delay"`
}

type ProductClient struct {
	logger log.Logger
	config ProductClientConfig
	client *http.Client
}

func NewProductClient(logger log.Logger, config ProductClientConfig) *ProductClient {
	return &ProductClient{
		logger: logger,
		config: config,
		client: &http.Client{},
	}
}

type response struct {
	State          int `json:"state"`
	Version        int `json:"version"`
	PayloadVersion int `json:"payloadVersion"`
	Data           struct {
		Products []struct {
			ID     uint64  `json:"id"`
			Name   string  `json:"name"`
			Rating float32 `json:"reviewRating"`

			Brand   string `json:"brand"`
			BrandID uint64 `json:"brandId"`

			Colors []struct {
				Name string `json:"name"`
			} `json:"colors"`

			Sizes []struct {
				Name  string `json:"name"`
				Price struct {
					Basic     float32 `json:"basic"`
					Product   float32 `json:"product"`
					Total     float32 `json:"total"`
					Logistics float32 `json:"logistics"`
					Return    int     `json:"return"`
				} `json:"price"`
			} `json:"sizes"`
		} `json:"products"`
	} `json:"data"`
}

func (c *ProductClient) GetProducts(_ context.Context, request model.ProductsRequest) ([]model.Product, error) {
	url := fmt.Sprintf(c.config.RequestURL, request.Category, request.Page)
	httpRequest, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ProductClient.GetData making http request error: %w", err)
	}

	httpRequest.Header.Add("Accept", "*/*")
	httpRequest.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

	var httpResponse *http.Response
	var retryCount uint

	for {
		httpResponse, err = c.client.Do(httpRequest)
		if err != nil {
			return nil, fmt.Errorf("ProductClient.GetData making http request error: %w", err)
		}

		if httpResponse.StatusCode == http.StatusTooManyRequests ||
			httpResponse.StatusCode == http.StatusGatewayTimeout ||
			httpResponse.StatusCode == http.StatusServiceUnavailable {
			if retryCount == c.config.RetryCount {
				break
			}

			c.logger.Warn().
				Int("status", httpResponse.StatusCode).
				Uint("try", retryCount+1).
				Msg("ProductClient.GetData response bad status trying to retry")

			if err = httpResponse.Body.Close(); err != nil {
				c.logger.Warn().Err(err).Msg("ProductClient.GetData failed to close response body")
			}

			time.Sleep(c.config.RetryDelay * (1 + time.Duration(retryCount)))
			retryCount++
			continue
		}

		break
	}

	defer func() {
		if cErr := httpResponse.Body.Close(); cErr != nil {
			c.logger.Warn().Err(cErr).Msg("ProductClient.GetData failed to close response body")
		}
	}()

	if httpResponse.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if httpResponse.StatusCode != http.StatusOK {
		if httpResponse.StatusCode == http.StatusTooManyRequests {
			return nil, model.ErrRequestLimit
		}
		return nil, fmt.Errorf("ProductClient.GetData http status is not ok; status: %d", httpResponse.StatusCode)
	}

	var respData response
	if err = json.NewDecoder(httpResponse.Body).Decode(&respData); err != nil {
		return nil, fmt.Errorf("ProductClient.GetData decode http request body error: %w", err)
	}

	var result []model.Product
	for _, item := range respData.Data.Products {
		colors := make([]string, 0, len(item.Colors))
		for _, color := range item.Colors {
			colors = append(colors, color.Name)
		}

		product := model.Product{
			ID:     item.ID,
			Name:   item.Name,
			Rating: item.Rating,
			URL:    fmt.Sprintf(c.config.ProductURL, item.ID),

			Brand:   item.Brand,
			BrandID: item.BrandID,

			Colors: colors,
		}

		for _, size := range item.Sizes {
			priceValue := fmt.Sprintf("%.f", size.Price.Total)
			strPriceValue := priceValue[:len(priceValue)-2]
			floatPriceValue, pErr := strconv.ParseFloat(strPriceValue, 32)
			if pErr != nil {
				return nil, fmt.Errorf("ProductClient.GetData parse product:%d price error: %w", item.ID, err)
			}

			size := model.ProductSize{
				Name:         size.Name,
				CurrentPrice: float32(floatPriceValue),
			}

			product.Sizes = append(product.Sizes, size)
		}

		result = append(result, product)
	}

	return result, nil
}
