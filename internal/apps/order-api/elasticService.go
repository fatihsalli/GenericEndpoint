package order_api

import (
	"GenericEndpoint/internal/configs"
	"GenericEndpoint/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func (e *ElasticService) GetFromElasticsearch(req OrderGetRequest) ([]interface{}, error) {

	searchBody := make(map[string]interface{})
	query := make(map[string]interface{})

	// Creating query for exact filters
	if len(req.ExactFilters) > 0 {
		boolQuery := make(map[string]interface{})
		mustClauses := make([]map[string]interface{}, 0)

		for field, values := range req.ExactFilters {
			if len(values) > 0 {
				mustClause := make(map[string]interface{})
				mustClause["terms"] = map[string]interface{}{
					field: values,
				}
				mustClauses = append(mustClauses, mustClause)
			}
		}

		boolQuery["must"] = mustClauses
		query["bool"] = boolQuery
	}

	// TODO: Match çalışmıyor kontrol edilecek
	// Creating query for match
	if len(req.Match) > 0 {
		for field, value := range req.Match {
			if query[field] != nil {
				// Adding a match query to exist bool query
				boolQuery := query[field].(map[string]interface{})
				boolQuery["match"] = map[string]interface{}{
					field: value,
				}
				query[field] = boolQuery
			} else {
				// Creating new bool query
				boolQuery := make(map[string]interface{})
				boolQuery["match"] = map[string]interface{}{
					field: value,
				}
				query[field] = boolQuery
			}
		}
	}

	searchBody["query"] = query

	if len(req.Sort) > 0 {
		for field, value := range req.Sort {
			if value == -1 {
				searchBody["sort"] = map[string]interface{}{
					field: "desc",
				}
			} else if value == 1 {
				searchBody["sort"] = map[string]interface{}{
					field: "asc",
				}
			}
		}
	}

	if len(req.Fields) > 0 {
		searchBody["_source"] = req.Fields
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(searchBody); err != nil {
		fmt.Println("Error encoding the query: ", err)
		return nil, err
	}

	res, err := e.ElasticClient.Search(
		e.ElasticClient.Search.WithIndex(e.Config.Elasticsearch.IndexName["Order"]),
		e.ElasticClient.Search.WithBody(buf),
	)
	if err != nil {
		fmt.Println("Error executing the search: ", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		fmt.Println("Error executing the decode: ", err)
		return nil, err
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Println("Error executing the decode: ", err)
		return nil, err
	}

	var orders []interface{}

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {

		// Casting with type assertion
		source, ok := hit.(map[string]interface{})["_source"]
		if !ok {
			fmt.Println("Source not found in the hit", err)
			return nil, err
		}
		orders = append(orders, source)
	}

	return orders, nil
}
