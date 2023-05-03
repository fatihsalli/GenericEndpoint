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

func (e *ElasticService) GetFromElasticsearch(req OrderGetRequest) ([]models.Order, error) {

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"status": "Shipped",
						},
					},
					map[string]interface{}{
						"match": map[string]interface{}{
							"status": "Delivered",
						},
					},
				},
			},
		},
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(query); err != nil {
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

	var orders []models.Order

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {

		source, err := json.Marshal(hit.(map[string]interface{})["_source"])
		if err != nil {
			fmt.Println("Error marshalling the source: ", err)
			return nil, err
		}

		var order models.Order
		err = json.Unmarshal(source, &order)
		if err != nil {
			fmt.Println("Error unmarshalling the order: ", err)
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}
