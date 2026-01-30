package store

import (
	"bytes"
	"context"
	"encoding/json"
	"analytics/types"
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

func (s *Store) QuerySpans(ctx context.Context, query map[string]interface{}) ([]types.Span, error) {
	data, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %v", err)
	}

	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(s.index),
		s.client.Search.WithBody(bytes.NewReader(data)),
	)
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error in search response: %s", res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source types.Span `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %v", err)
	}

	spans := make([]types.Span, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		spans[i] = hit.Source
	}

	return spans, nil
}