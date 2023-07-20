package model

import "time"

type Order struct {
	Id       string           `json:"id"`
	Total    float32          `json:"total"`
	Products []ProductInOrder `json:"products"`
	Date     time.Time        `json:"date"`
}
