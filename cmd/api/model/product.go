package model

import "errors"

var ErrRequestLimit = errors.New("request limit exceeded")

type ProductsRequest struct {
	Page       int    `json:"page"`
	Category   string `json:"category"`
	CategoryID uint64 `json:"category_id"`
	RequestURL string `json:"requestUrl"`
	ProductURL string `json:"productUrl"`
}

type ProductSize struct {
	Name          string  `json:"name"`
	FirstPrice    float32 `json:"firstPrice"`
	PreviousPrice float32 `json:"previousPrice"`
	CurrentPrice  float32 `json:"currentPrice"`
}

type Product struct {
	ID         uint64  `json:"id"`
	CategoryID uint64  `json:"category_id"`
	Name       string  `json:"name"`
	Rating     float32 `json:"rating"`
	URL        string  `json:"url"`

	Brand   string `json:"brand"`
	BrandID uint64 `json:"brandId"`

	Colors []string      `json:"colors"`
	Sizes  []ProductSize `json:"sizes"`
}
