package database

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go-app/src/logger"
	"go-app/src/model"
	"go-app/src/service"
	"go-app/src/service/impl/query"
	"go-app/src/service/spec/apperror"
	"regexp"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"go.uber.org/zap"
)

type es8Repository struct {
	EsClient  *elasticsearch.Client
	IndexName string
	Pipeline  string
}

func NewEs8Repository(client *elasticsearch.Client, indexName string, pipeline string) *es8Repository {
	repo := &es8Repository{EsClient: client, IndexName: indexName}
	if pipeline != "" {
		repo.Pipeline = pipeline
	}
	return repo
}

func (r *es8Repository) Index(id string, item any) error {
	data, err := json.Marshal(item)
	if err != nil {
		logger.Warn("EsCrud", zap.Error(err), zap.Any("item", item))
		return err
	}
	var req esapi.IndexRequest
	if id != "" {
		req = esapi.IndexRequest{
			Index:      r.IndexName,
			DocumentID: id,
			Body:       bytes.NewReader(data),
		}
	} else {
		req = esapi.IndexRequest{
			Index: r.IndexName,
			Body:  bytes.NewReader(data),
		}
	}
	// Set up the request object.

	if r.Pipeline != "" {
		req.Pipeline = r.Pipeline
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

func (r *es8Repository) GetOne(id string) (any, error) {
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

func removeNumbers(input string) string {
	re := regexp.MustCompile("[0-9]+")
	return re.ReplaceAllString(input, "")
}

type jmap = map[string]any

func (r *es8Repository) GetList(must jmap, should jmap, mustNot jmap, listParam *model.ListParam) (ret []any, total int64, retErr error) {
	q, _ := service.DeepCopyMap(query.QUERY_CRUD_GETLIST)
	mustcondition := []service.JsonMap{}
	for k, v := range must {
		term := jmap{}
		result := removeNumbers(k)
		term[result] = jmap{"value": v}
		mustcondition = append(mustcondition, jmap{"term": term})
	}

	if should != nil {
		q["query"].(jmap)["bool"].(jmap)["minimum_should_match"] = 1
	}

	mustNotcondition := []service.JsonMap{}
	for k, v := range mustNot {
		term := jmap{}
		result := removeNumbers(k)
		term[result] = jmap{"value": v}
		mustNotcondition = append(mustNotcondition, jmap{"term": term})
	}

	shouldcondition := []service.JsonMap{}
	for k, v := range should {
		term := jmap{}
		result := removeNumbers(k)
		term[result] = jmap{"value": v}
		shouldcondition = append(shouldcondition, jmap{"term": term})
	}

	q["size"] = listParam.Size

	if listParam.TimeEnd != nil && listParam.TimeStart != nil {
		mustcondition = append(mustcondition, jmap{"range": jmap{"timestamp": jmap{"gte": listParam.TimeStart, "lt": listParam.TimeEnd}}})
	}
	q["query"].(jmap)["bool"].(jmap)["must"] = mustcondition
	q["query"].(jmap)["bool"].(jmap)["should"] = shouldcondition
	q["query"].(jmap)["bool"].(jmap)["must_not"] = mustNotcondition

	if listParam.Aggs != nil {
		q["aggs"] = listParam.Aggs
	}

	if listParam.SearchAfter != nil && listParam.SortAf != nil {
		q["search_after"] = listParam.SearchAfter
		q["sort"] = listParam.SortAf
	}

	data, _ := json.Marshal(q)
	req := esapi.SearchRequest{
		Index: []string{r.IndexName},
		Body:  bytes.NewReader(data),
	}
	res, err := req.Do(context.Background(), r.EsClient)
	if err != nil {
		logger.Warn("EsCrud", zap.Error(err))
		return nil, 0, err
	}

	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	logger.Info("q", zap.Any("q", q))
	if res.IsError() {
		logger.Warn("EsCrud", zap.String("status", res.Status()))
		return nil, 0, errors.New(res.Status())
	}
	if err != nil {
		logger.Warn("EsCrud", zap.Error(err))
		return nil, 0, err
	}
	var result jmap
	ret = []any{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		service.PanicOnError(apperror.ErrInternal, zap.Any("error", err), zap.String("cause", "ParseResponse"))
	}
	if *listParam.AggsType == "default" {
		hits := result["hits"].(jmap)
		fTotal, _ := hits["total"].(jmap)["value"].(float64)
		total = int64(fTotal)
		hitsAny := hits["hits"].([]any)
		for _, itemAny := range hitsAny {
			itemM := itemAny.(jmap)
			ret = append(ret, itemM["_source"])
		}
	} else if *listParam.AggsType == "aggs" {
		aggs := result["aggregations"].(jmap)["key"].(jmap)["buckets"].([]any)
		ret = aggs
	}

	return ret, total, err
}

func (r *es8Repository) Delete(id string) error {
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
