package model

type Category struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Emoji      string `json:"emoji"`
	RequestURL string `json:"requestUrl"`
	ProductURL string `json:"productUrl"`
}
