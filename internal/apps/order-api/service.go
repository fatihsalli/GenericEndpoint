package order_api

import (
	"GenericEndpoint/internal/models"
	"GenericEndpoint/internal/repository"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	Repository *repository.Repository
}

func NewService(Repository *repository.Repository) *Service {
	service := &Service{Repository: Repository}
	return service
}

func (s *Service) GetAll() ([]models.Order, error) {
	result, err := s.Repository.GetAll()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) Insert(order models.Order) (models.Order, error) {
	// Create id and created date value
	order.ID = uuid.New().String()
	order.CreatedAt = time.Now()
	// We don't want to set null, so we put CreatedAt value.
	order.UpdatedAt = order.CreatedAt

	var total float64
	for _, product := range order.Product {
		total = product.Price * float64(product.Quantity)
		order.Total += total
	}

	_, err := s.Repository.CreateOrder(order)

	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}
