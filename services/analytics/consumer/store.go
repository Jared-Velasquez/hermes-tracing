package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	es "github.com/elastic/go-elasticsearch/v8"
)

type Store struct {
	client *es.Client
	index  string
}

func NewStore(addresses []string, index string) (*Store, error) {
	cfg := es.Config{
		Addresses: addresses,
	}
	client, err := es.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %v", err)
	}

	return &Store{
		client: client,
		index:  index,
	}, nil
}

func (s *Store) IndexSpan (ctx context.Context, doc map[string]interface{}) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %v", err)
	}

	res, err := s.client.Index(
		s.index,
		bytes.NewReader(data),
		s.client.Index.WithContext(ctx),
		s.client.Index.WithDocumentID(doc["span_id"].(string)),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document ID=%s: %s", doc["span_id"], res.String())
	}

	return nil
}