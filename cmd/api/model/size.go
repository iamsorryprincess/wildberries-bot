package model

type SizeInfo struct {
	ID            uint64 `json:"id"`
	Name          string `json:"name"`
	ProductsCount uint   `json:"productsCount"`
}

type SizeCategoryInfo struct {
	Name          string `json:"name"`
	CategoryTitle string `json:"categoryTitle"`
	CategoryEmoji string `json:"categoryEmoji"`
}
