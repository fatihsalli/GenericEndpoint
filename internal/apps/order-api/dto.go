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

type OrderResponse struct {
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
	Total     float64 `json:"total,omitempty" bson:"total"`
	CreatedAt string  `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt string  `json:"updatedAt,omitempty" bson:"updatedAt"`
}
