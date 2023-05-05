package order_api

import (
	"GenericEndpoint/internal/models"
	"GenericEndpoint/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoService struct {
	Repository *repository.Repository
}

func NewService(Repository *repository.Repository) *MongoService {
	service := &MongoService{Repository: Repository}
	return service
}

func (s *MongoService) GetAll() ([]models.Order, error) {
	result, err := s.Repository.GetAll()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MongoService) GetOrdersWithFilter(filter bson.M, findOptions *options.FindOptions) ([]models.Order, error) {
	result, err := s.Repository.GetOrdersWithFilter(filter, findOptions)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MongoService) Insert(order models.Order) (models.Order, error) {
	_, err := s.Repository.Insert(order)

	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (s *MongoService) Delete(id string) (bool, error) {
	result, err := s.Repository.Delete(id)

	if err != nil {
		return false, err
	}

	return result, nil
}

func (s *MongoService) FromModelConvertToFilter(req OrderGetRequest) (bson.M, *options.FindOptions) {

	// Create a filter based on the exact filters and matches provided in the request
	filter := bson.M{}

	// Add exact filter criteria to filter if provided
	if len(req.ExactFilters) > 0 {
		for key, values := range req.ExactFilters {
			filter[key] = bson.M{"$in": values}
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

	// Add projection criteria to find options if provided
	if len(req.Fields) > 0 {
		projection := bson.M{}
		findOptions.SetProjection(projection)
		for _, field := range req.Fields {
			projection[field] = 1
		}
	}

	// Add sort criteria to find options if provided
	if len(req.Sort) > 0 {
		sort := bson.M{}
		for key, value := range req.Sort {
			sort[key] = value
		}
		findOptions.SetSort(sort)
	}

	return filter, findOptions
}
