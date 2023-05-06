package models

import "time"

type Order struct {
	ID            string `json:"id,omitempty" bson:"_id"`
	UserID        string `json:"userId,omitempty" bson:"userId"`
	Status        string `json:"status,omitempty" bson:"status"`
	City          string `json:"city,omitempty" bson:"city"`
	AddressDetail string `json:"addressDetail,omitempty" bson:"addressDetail"`
	Product       []struct {
		Name     string  `json:"name" bson:"name"`
		Quantity int     `json:"quantity" bson:"quantity"`
		Price    float64 `json:"price" bson:"price"`
	} `json:"product,omitempty" bson:"product"`
	Total     float64   `json:"total,omitempty" bson:"total"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt"`
}
