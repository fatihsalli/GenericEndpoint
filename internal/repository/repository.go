package repository

import (
	"GenericEndpoint/internal/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository struct {
	Collection *mongo.Collection
}

func NewRepository(mongoCollection *mongo.Collection) *Repository {
	repository := &Repository{Collection: mongoCollection}
	return repository
}

// GetAll Method => to list every order
func (r *Repository) GetAll() ([]models.Order, error) {
	var order models.Order
	var orders []models.Order

	// open connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//We can think of "Cursor" like a request. We pull the data from the database with the "Next" command. (C# => IQueryable)
	result, err := r.Collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	for result.Next(ctx) {
		if err := result.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *Repository) CreateOrder(order models.Order) (bool, error) {
	// open connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := r.Collection.InsertOne(ctx, order)

	if result.InsertedID == nil || err != nil {
		return false, errors.New("failed to add")
	}

	return true, nil
}
