package models

import "time"

type Order struct {
	ID            string `json:"id" bson:"_id"`
	UserID        string `json:"userId" bson:"userId"`
	Status        string `json:"status" bson:"status"`
	City          string `json:"city" bson:"city"`
	AddressDetail string `json:"addressDetail" bson:"addressDetail"`
	Product       []struct {
		Name     string  `json:"name" bson:"name"`
		Quantity int     `json:"quantity" bson:"quantity"`
		Price    float64 `json:"price" bson:"price"`
	} `json:"product" bson:"product"`
	Total     float64   `json:"total" bson:"total"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
