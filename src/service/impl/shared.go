package impl

import (
	"context"
	"go-app/src/database"
	"go-app/src/logger"
	"go-app/src/model"
	"go-app/src/service"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type jmap = service.JsonMap

// Get context with existed Transaction Tx or create new Tx
func GetContextTx(c context.Context, create bool) (context.Context, *sqlx.Tx, error) {
	val := c.Value("__tx")
	if val != nil {
		tx, ok := val.(*sqlx.Tx)
		if ok {
			return c, tx, nil
		}
	}
	if !create {
		return c, nil, nil
	}
	tx, retErr := database.GetDbTransaction()
	if retErr != nil || tx == nil {
		logger.Warn("TxFailure", zap.Error(retErr))
		return c, nil, retErr
	}
	return context.WithValue(c, "__tx", tx), tx, nil
}

type IEsCrudRepository interface {
	Index(id string, item any) error
	GetOne(id string) (any, error)
	GetList(pattern jmap, listParam *model.ListParam) (ret []any, total int64, retErr error)
	Delete(id string) error
}
