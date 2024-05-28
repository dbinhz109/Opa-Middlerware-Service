package spec

import (
	"context"
	"go-app/src/model"
	"go-app/src/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type WazuhQueryParam struct {
	AccountId   int64             `json:"account_id,omitempty"`
	TimeStart   int64             `json:"time_start,omitempty"`
	TimeEnd     int64             `json:"time_end,omitempty"`
	Offset      int               `json:"offset,omitempty"`
	Limit       int               `json:"limit,omitempty"`
	Label       *string           `json:"label,omitempty"`
	Aggregation bool              `json:"aggregation,omitempty"`
	Must        []service.JsonMap `json:"must,omitempty"`
	Should      []service.JsonMap `json:"should,omitempty"`
	MustNot     []service.JsonMap `json:"must_not,omitempty"`
}

type QueryParam struct {
	Search      *string       `json:"search,omitempty"`
	AccountId   *int64        `json:"account_id,omitempty"`
	TimeStart   int64         `json:"time_start,omitempty"`
	TimeEnd     int64         `json:"time_end,omitempty"`
	Offset      int           `json:"offset,omitempty"`
	Limit       int           `json:"limit,omitempty"`
	Label       *string       `json:"label,omitempty"`
	Aggregation bool          `json:"aggregation,omitempty"`
	LevelLTE    int           `json:"level_lte,omitempty"`
	LevelGTE    int           `json:"level_gte,omitempty"`
	AgentId     *string       `json:"agent_id,omitempty"`
	SortAf      []interface{} `json:"sort_af,omitempty"`
	SearchAfter []any         `json:"search_after,omitempty"`
}

type WzApiQueryParam struct {
	RuleIds  *int64  `json:"rule_ids,omitempty"`
	Offset   *int64  `json:"offset,omitempty"`
	Limit    *int64  `json:"limit,omitempty"`
	Sort     *string `json:"sort,omitempty"`
	Search   *string `json:"search,omitempty"`
	Level    *string `json:"level,omitempty"`
	Filename *string `json:"filename,omitempty"`
	Status   *string `json:"status,omitempty"`
}

// IDefaultCrudController ...
type IDefaultCrudController interface {
	Create(c *gin.Context)
	GetOne(c *gin.Context)
	GetList(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// Default CRUD service interface with T being pointer to an LSEntity type
type IDefaultCrudService[T model.EntityModel[any]] interface {
	Create(context context.Context, entity T) (T, error)
	GetOne(context context.Context, sample T) (T, error)
	GetList(context context.Context, sample T, param *model.ListParam) ([]T, int64, error)
	Update(context context.Context, sample T, value T) error
	Delete(context context.Context, sample T) (retErr error)
}

// Default CRUD repository interface with T being pointer to an LSEntity type
type IDefaultCrudRepository[T model.EntityModel[any]] interface {
	Create(entity T) (T, error)
	GetOne(entity T) (T, error)
	GetList(entity T, param_ *model.ListParam) ([]T, int64, error)
	UpdateAll(sample T, value T) (int64, error)
	UpdateNull(sample T, value T) (int64, error)
	DeleteSome(entity T) (int64, error)
}

type IRepositoryFactory[T model.EntityModel[any]] func(tx *sqlx.Tx) IDefaultCrudRepository[T]
