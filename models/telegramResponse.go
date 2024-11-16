package models

type TelegramResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"result"`
	Description string `json:"description"`
}
