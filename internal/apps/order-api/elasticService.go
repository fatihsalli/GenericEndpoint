package order_api

import (
	"GenericEndpoint/internal/configs"
	"GenericEndpoint/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/labstack/gommon/log"
)

type ElasticService struct {
	Config        *configs.Config
	ElasticClient *elasticsearch.Client
}

func NewElasticService(config *configs.Config) *ElasticService {
	// client with default config
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.Elasticsearch.Addresses["Address 1"],
		},
	}

	elasticClient, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Errorf("Error creating the client: ", err)
	}

	elasticService := &ElasticService{Config: config, ElasticClient: elasticClient}
	return elasticService
}

func (e *ElasticService) SaveOrderToElasticsearch(order models.Order) error {
	// Build the request body.
	data, err := json.Marshal(order)
	if err != nil {
		log.Errorf("Error marshaling document: %s", err)
		return err
	}

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      e.Config.Elasticsearch.IndexName["Order"],
		DocumentID: order.ID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), e.ElasticClient)
	if err != nil {
		log.Errorf("Error getting response: %s", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Errorf("Error parsing the response body: %s", err)
			return err
		} else {
			// Print the error information.
			log.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	return nil
}

func (e *ElasticService) DeleteOrderFromElasticsearch(orderID string) error {
	// Create request object
	req := esapi.DeleteRequest{
		Index:      e.Config.Elasticsearch.IndexName["Order"],
		DocumentID: orderID,
		Refresh:    "true",
	}

	// Execute the request
	res, err := req.Do(context.Background(), e.ElasticClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Errorf("Error parsing the response body: %s", err)
			return err
		} else {
			// Print the error information.
			log.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	return nil
}

func (e *ElasticService) FromModelConvertToSearchRequest(req OrderGetRequest) (bytes.Buffer, error) {
	// Create a buffer to hold the query
	var buf bytes.Buffer

	// Create a map to hold the query filters
	queryFilters := make(map[string]interface{})

	// Add exact filter criteria to the query filters if provided
	if len(req.ExactFilters) > 0 {
		for key, values := range req.ExactFilters {
			if len(values) > 1 {
				// If there are multiple values for the same key, we use the terms query
				queryFilters["terms"] = map[string]interface{}{key: values}
			} else {
				// Otherwise, we use the term query
				queryFilters["term"] = map[string]interface{}{key: values[0]}
			}
		}
	}

	// Add match criteria to the query filters if provided
	if len(req.Match) > 0 {
		for key, value := range req.Match {
			queryFilters["match"] = map[string]interface{}{key: value}
		}
	}

	// Add sort criteria to the query if provided
	if len(req.Sort) > 0 {
		sort := make([]map[string]interface{}, 0)
		for key, value := range req.Sort {
			sort = append(sort, map[string]interface{}{key: map[string]interface{}{"order": value}})
		}
		queryFilters["sort"] = sort
	}

	// Add projection criteria to the query if provided
	if len(req.Fields) > 0 {
		projection := make(map[string]interface{})
		for _, field := range req.Fields {
			projection[field] = true
		}
		queryFilters["_source"] = projection
	}

	// Construct the Elasticsearch query from the query filters
	query := map[string]interface{}{"query": map[string]interface{}{"bool": queryFilters}}

	// Serialize the query to JSON and write it to the buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return buf, err
	}

	return buf, nil
}

func (e *ElasticService) GetFromElasticSearch(query bytes.Buffer) (*esapi.Response, error) {
	// Create the search request
	req := esapi.SearchRequest{
		Index:          []string{e.Config.Elasticsearch.IndexName["Order"]},
		Body:           bytes.NewReader(query.Bytes()),
		TrackTotalHits: true,
	}

	// Execute the search request and return the response
	res, err := req.Do(context.Background(), e.ElasticClient)
	if err != nil {
		return nil, err
	}

	return res, nil
}
