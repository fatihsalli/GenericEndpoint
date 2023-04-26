package models

import "time"

type Order struct {
	ID       string `json:"id" bson:"_id"`
	UserID   string `json:"userId" bson:"userId"`
	Status   string `json:"status" bson:"status"`
	Country  string `json:"country" bson:"country"`
	City     string `json:"city" bson:"city"`
	District string `json:"district" bson:"district"`
	Product  []struct {
		Name     string  `json:"name" bson:"name"`
		Quantity int     `json:"quantity" bson:"quantity"`
		Price    float64 `json:"price" bson:"price"`
	} `json:"product" bson:"product"`
	AddressDetail string    `json:"addressDetail" bson:"addressDetail"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
}

type OrderRequest struct {
	exactFilters map[string]string
}
