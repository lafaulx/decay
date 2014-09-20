package model

type Dcy struct {
	Id        int    `json:"id"`
	CallCount int64  `json:"callCount"`
	Content   string `json:"content"`
}
