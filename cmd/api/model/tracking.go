package model

type TrackingSettings struct {
	ChatID     int64  `json:"chatId"`
	SizeID     uint64 `json:"sizeId"`
	CategoryID uint64 `json:"categoryId"`
	DiffValue  int    `json:"diffValue"`
}

type TrackingResult struct {
	ChatID        int64
	ProductName   string
	ProductURL    string
	Size          string
	PreviousPrice string
	CurrentPrice  string
	DiffPercent   int
}
