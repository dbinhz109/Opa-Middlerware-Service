package database

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go-app/src/logger"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"go.uber.org/zap"
)

type es7Repository struct {
	EsClient  *elasticsearch.Client
	IndexName string
}

func NewEs7Repository(client *elasticsearch.Client, indexName string) *es7Repository {
	return &es7Repository{EsClient: client, IndexName: indexName}
}

func (r *es7Repository) Index(id string, item any) error {
	data, err := json.Marshal(item)
	if err != nil {
		logger.Warn("EsCrud", zap.Error(err), zap.Any("item", item))
		return err
	}
	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      r.IndexName,
		DocumentID: id,
		Body:       bytes.NewReader(data),
	}
	// Perform the request with the client.
	res, err := req.Do(context.Background(), r.EsClient)
	if err != nil {
		logger.Warn("EsCrud", zap.Error(err))
		return err
	}

	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()

	if res.IsError() {
		logger.Warn("EsCrud", zap.String("status", res.Status()), zap.Any("item", item))
		return errors.New(res.Status())
	}
	if err != nil {
		logger.Warn("EsCrud", zap.Error(err), zap.Any("item", item))
	}
	return err
}

func (r *es7Repository) GetOne(id string) (any, error) {
	req := esapi.GetRequest{
		Index:      r.IndexName,
		DocumentID: id,
	}
	res, err := req.Do(context.Background(), r.EsClient)
	if err != nil {
		logger.Warn("EsCrud", zap.Error(err))
		return nil, err
	}

	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()

	if res.IsError() {
		logger.Warn("EsCrud", zap.String("status", res.Status()), zap.String("id", id))
		return nil, errors.New(res.Status())
	}
	if err != nil {
		logger.Warn("EsCrud", zap.Error(err), zap.String("id", id))
		return nil, err
	}
	return nil, err
}
func (r *es7Repository) Delete(id string) error {
	req := esapi.DeleteRequest{
		Index:      r.IndexName,
		DocumentID: id,
	}
	res, err := req.Do(context.Background(), r.EsClient)
	if err != nil {
		logger.Warn("EsCrud", zap.Error(err))
		return err
	}

	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()

	if res.IsError() {
		logger.Warn("EsCrud", zap.String("status", res.Status()), zap.String("id", id))
		return errors.New(res.Status())
	}
	if err != nil {
		logger.Warn("EsCrud", zap.Error(err), zap.String("id", id))
		return err
	}
	return err
}
