package model

type TrackingSettings struct {
	ChatID    int64  `json:"chatId"`
	Size      string `json:"size"`
	Category  string `json:"category"`
	DiffValue int    `json:"diffValue"`
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
