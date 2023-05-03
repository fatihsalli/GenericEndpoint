package order_api

type OrderCreateRequest struct {
	UserID        string `json:"userId" bson:"userId"`
	Status        string `json:"status" bson:"status"`
	City          string `json:"city" bson:"city"`
	AddressDetail string `json:"addressDetail" bson:"addressDetail"`
	Product       []struct {
		Name     string  `json:"name" bson:"name"`
		Quantity int     `json:"quantity" bson:"quantity"`
		Price    float64 `json:"price" bson:"price"`
	} `json:"product" bson:"product"`
}

type OrderGetRequest struct {
	ExactFilters map[string][]interface{} `json:"exact_filters"`
	Fields       []string                 `json:"fields"`
	Match        map[string]interface{}   `json:"match"`
	Sort         map[string]int           `json:"sort"`
}
