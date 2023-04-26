package main

import "time"

type Order struct {
	ID       string
	UserID   string
	Status   string
	Country  string
	City     string
	District string
	Product  []struct {
		Name     string
		Quantity int
		Price    float64
	}
	AddressDetail string
	CreatedDate   time.Time
	UpdatedDate   time.Time
}

type OrderRequest struct {
	exactFilters map[string]string
}
