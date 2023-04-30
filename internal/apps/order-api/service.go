package order_api

import (
	"GenericEndpoint/internal/models"
	"GenericEndpoint/internal/repository"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (s *Service) GetOrdersWithFilter(filter bson.M, findOptions *options.FindOptions) ([]interface{}, error) {
	result, err := s.Repository.GetOrdersWithFilter(filter, findOptions)

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

	_, err := s.Repository.Insert(order)

	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (s *Service) Delete(id string) (bool, error) {
	result, err := s.Repository.Delete(id)

	if err != nil {
		return false, err
	}

	return result, nil
}

func (s *Service) FromModelConvertToFilter(req OrderGetRequest) (bson.M, *options.FindOptions) {

	// Create a filter based on the exact filters and matches provided in the request
	filter := bson.M{}

	// Add exact filters to filter if provided
	if len(req.ExactFilters) > 0 {
		for key, value := range req.ExactFilters {
			filter[key] = value
		}
	}

	// Add match criteria to filter if provided
	if len(req.Match) > 0 {
		match := bson.M{}
		for key, value := range req.Match {
			match[key] = value
		}
		filter = bson.M{
			"$and": []bson.M{
				filter,
				match,
			},
		}
	}

	// Create options for the find operation, including the requested fields and sort order
	findOptions := options.Find()

	if len(req.Fields) > 0 {
		projection := bson.M{}
		findOptions.SetProjection(projection)
		for _, field := range req.Fields {
			projection[field] = 1
		}
	}

	if len(req.Sort) > 0 {
		sort := bson.M{}
		for key, value := range req.Sort {
			sort[key] = value
		}
		findOptions.SetSort(sort)
	}

	return filter, findOptions
}
