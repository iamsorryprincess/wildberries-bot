package model

import (
	"context"
	"time"
)

const (
	ProductCategoryDresses = "dresses"
)

type ProductsRequest struct {
	Page     int    `json:"page"`
	Category string `json:"category"`
}

type ProductClient interface {
	GetProducts(ctx context.Context, request ProductsRequest) ([]Product, error)
}

type Product struct {
	ID     uint64  `json:"id"`
	Name   string  `json:"name"`
	Rating float32 `json:"rating"`
	Size   string  `json:"size"`

	Brand   string `json:"brand"`
	BrandID uint64 `json:"brandId"`

	Colors []string `json:"colors"`

	FirstValue    float32 `json:"firstValue"`
	PreviousValue float32 `json:"previousValue"`
	CurrentValue  float32 `json:"currentValue"`

	FirstDate time.Time `json:"firstDate"`
	LastDate  time.Time `json:"lastDate"`
}
