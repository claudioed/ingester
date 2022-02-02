package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/google/uuid"
	"ingester/internal/infra/elastic"
	ingester_v1 "ingester/pkg/pb/analytics"
	"time"
)

type ElasticApiCallRepository struct {
	elastic *elastic.ElasticSearch
	timeout time.Duration
}

func (er *ElasticApiCallRepository) Add(ctx context.Context, call *ingester_v1.ApiCall) (*string, error) {
	id := uuid.New().String()
	m := jsonpb.Marshaler{}
	js, err := m.MarshalToString(call)
	req := esapi.CreateRequest{
		Index:      er.elastic.Alias,
		DocumentID: id,
		Body:       bytes.NewReader([]byte(js)),
	}
	ctx, cancel := context.WithTimeout(ctx, er.timeout)
	defer cancel()
	res, err := req.Do(ctx, er.elastic.Client)
	if err != nil {
		return nil, fmt.Errorf("insert: request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode == 409 {
		return nil, fmt.Errorf("conflict: %w", err)
	}
	if res.IsError() {
		return nil, fmt.Errorf("insert: response: %s", res.String())
	}
	return &id, nil
}

func NewElasticApiCallRepository(elastic *elastic.ElasticSearch) ApiCallRepo {
	return &ElasticApiCallRepository{
		elastic: elastic,
		timeout: time.Second * 2,
	}
}
