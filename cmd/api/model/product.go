package model

const (
	ProductCategoryDresses = "dresses"
)

type ProductsRequest struct {
	Page     int    `json:"page"`
	Category string `json:"category"`
}

type ProductSize struct {
	Name          string  `json:"name"`
	FirstPrice    float32 `json:"firstPrice"`
	PreviousPrice float32 `json:"previousPrice"`
	CurrentPrice  float32 `json:"currentPrice"`
}

type Product struct {
	ID     uint64  `json:"id"`
	Name   string  `json:"name"`
	Rating float32 `json:"rating"`
	URL    string  `json:"url"`

	Brand   string `json:"brand"`
	BrandID uint64 `json:"brandId"`

	Colors []string      `json:"colors"`
	Sizes  []ProductSize `json:"sizes"`
}
