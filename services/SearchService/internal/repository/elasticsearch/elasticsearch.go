package elasticsearch

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ReilEgor/Vaca/services/SearchService/internal/config"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

func NewElasticClient(address config.ElasticSearchURL) (*elasticsearch.TypedClient, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			string(address),
		},
	}

	typedClient, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elastic client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := typedClient.Ping().Do(ctx); err != nil {
		return nil, fmt.Errorf("could not ping elasticsearch: %w", err)
	}

	log.Println("Successfully connected to Elasticsearch")
	return typedClient, nil
}
