package model

type TrackingSettings struct {
	ChatID     int64  `json:"chatId"`
	SizeID     uint64 `json:"sizeId"`
	CategoryID uint64 `json:"categoryId"`
	DiffValue  int    `json:"diffValue"`
}

type TrackingResult struct {
	ChatID          int64
	ProductID       uint64
	ProductName     string
	ProductURL      string
	SizeID          uint64
	Size            string
	PreviousPrice   float32
	CurrentPrice    float32
	CurrentPriceInt uint64
	DiffPercent     int
}

type TrackingLog struct {
	ChatID    int64
	SizeID    uint64
	ProductID uint64
	Price     uint64
}

type TrackingSettingsInfo struct {
	ChatID        int64  `json:"chatId"`
	CategoryID    uint64 `json:"categoryId"`
	CategoryTitle string `json:"categoryTitle"`
	CategoryEmoji string `json:"categoryEmoji"`
	SizeID        uint64 `json:"sizeId"`
	Size          string `json:"size"`
	DiffPercent   int    `json:"diffPercent"`
}
