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
