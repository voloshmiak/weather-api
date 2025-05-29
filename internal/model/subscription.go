package model

type Subscription struct {
	Email     string `json:"email"`
	City      string `json:"city"`
	Frequency string `json:"frequency"`
	Confirmed bool   `json:"confirmed"`
}
