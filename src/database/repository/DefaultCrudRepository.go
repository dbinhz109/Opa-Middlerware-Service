package repository

import (
	"errors"
	"fmt"
	"go-app/src/logger"
	"go-app/src/model"
	"go-app/src/service/spec"
	"go-app/src/service/spec/apperror"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type DefaultCrudRepository[E model.EntityModel[any]] struct {
	tx                *sqlx.Tx
	RepositoryFactory spec.IRepositoryFactory[E]
}

// Create new repository instance on per transaction basis
func NewDefaultCrudRepository[E model.EntityModel[any]](tx *sqlx.Tx) *DefaultCrudRepository[E] {
	return &DefaultCrudRepository[E]{
		tx: tx,
	}
}

func (repo *DefaultCrudRepository[E]) Create(entity E) (E, error) {
	var defaultInstance E

	qb := QueryBuilder{}
	query, qparam := qb.BuildCreate(entity)
	query = repo.tx.Rebind(query)

	logger.Debug("Create", zap.Any("query", query), zap.Any("entity", entity))
	// _, err := repo.tx.NamedExec(query, entity)
	_, err := repo.tx.Exec(query, qparam...)
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Class().Name() == "integrity_constraint_violation" {
		err = apperror.ErrConstraintViolation
	}
	if err == nil {
		return entity, nil
	} else {
		logger.Error("Create", zap.Any("entity", entity))
		return defaultInstance, err
	}
}

// GetOne record with criteria provided by a sample object
func (repo *DefaultCrudRepository[E]) GetOne(entity E) (E, error) {
	var defaultInstance E

	qb := &QueryBuilder{}
	rule := NewRule(entity, entity, REL_EQUAL)
	qb.Select(entity).Where(rule).Limit(1)
	query, param := qb.BuildSelect()
	query = repo.tx.Rebind(query)
	rows, err := repo.tx.Queryx(query, param...)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			row := entity.CreateInstance()
			err = rows.StructScan(row)
			if err != nil {
				logger.Error("StructScan", zap.Any("param", entity), zap.Error(err))
			} else {
				rows.Close()
				return row.(E), nil
			}
		}
		rows.Close()
		return defaultInstance, nil
	} else {
		logger.Error("Create", zap.String("query", query), zap.Any("entity", entity), zap.Error(err))
		return defaultInstance, err
	}
}

// GetList find all records filtered by fields of *entity* with support counting, sorting, pagination
// @return data and total count
func (repo *DefaultCrudRepository[E]) GetList(entity E, param_ *model.ListParam) ([]E, int64, error) {
	var total int64 = 0
	param := &model.ListParam{Offset: 0, Limit: 11000}
	if param_ != nil {
		param = param_
	}
	if param.Limit <= 0 {
		param.Limit = 100
	}
	if param.Limit > 11000 {
		param.Limit = 11000
	}

	// whereClause := QB.WhereEqual(entity)
	// if len(whereClause) > 0 {
	// 	whereClause = "WHERE " + whereClause
	// }

	// sortClause := ""
	// if param.Sort != nil {
	// 	if len(*param.Sort) > 0 {
	// 		sortClause = strings.ReplaceAll(*param.Sort, " ", "")
	// 		if strings.HasSuffix(sortClause, "-") {
	// 			sortClause = strings.ReplaceAll(sortClause, "-", " DESC")
	// 		} else {
	// 			sortClause += " ASC"
	// 		}
	// 		sortClause = "ORDER BY " + sortClause
	// 	}
	// }

	qb := &QueryBuilder{}
	rule := NewRule(entity, entity, REL_EQUAL)
	qb.Select(entity).Where(rule)
	if param.Sort != nil && len(*param.Sort) > 0 {
		qb.OrderBy(strings.Split(*param.Sort, ","))
	}

	if param.Offset > 0 {
		qb.Offset(param.Offset)
	}

	if param.Limit > 0 {
		qb.Limit(param.Limit)
	}

	if param.Count {
		qb.count = true
		qb.Offset(0)
		qb.Limit(1)
		query, qparam := qb.BuildSelect()
		query = repo.tx.Rebind(query)
		logger.Debug("query", zap.Any("query", query))
		rows, err := repo.tx.Queryx(query, qparam...)
		if err != nil {
			logger.Warn("countFailed", zap.Error(err))
			return nil, 0, err
		}
		defer rows.Close()
		if rows.Next() {
			err = rows.Scan(&total)
			if err != nil {
				rows.Close()
				return nil, 0, err
			}
		}
		rows.Close()
	}

	qb.count = false
	qb.Offset(param.Offset)
	qb.Limit(param.Limit)
	query, qparam := qb.BuildSelect()
	query = repo.tx.Rebind(query)
	logger.Debug("query", zap.Any("query", query))
	rows, err := repo.tx.Queryx(query, qparam...)

	if err == nil {
		defer rows.Close()
		ret := make([]E, 0, param.Limit)
		for rows.Next() {
			row := entity.CreateInstance()
			err = rows.StructScan(row)
			if err != nil {
				logger.Error("", zap.Error(err))
			} else {
				ret = append(ret, row.(E))
			}
		}
		rows.Close()
		return ret, total, nil
	}
	logger.Warn("FindWithCount", zap.String("query", query), zap.Any("entity", entity), zap.Error(err))
	return nil, 0, err
}

// UpdateAll updates all records filtered by fields of *sample* with value fields from *value*
func (repo *DefaultCrudRepository[E]) UpdateAll(condition E, value E) (int64, error) {
	// Parameterized using ? syntax and rebind later, can't use named parameters here because of 2 objects with same structure
	args := make([]interface{}, 0, 1024)

	valueClause := QB.NamedSetValueList(value)
	valueArgs := QB.ValueList(value)
	args = append(args, valueArgs...)

	whereClause := QB.WhereEqualOrder(condition)
	if len(whereClause) > 0 {
		whereClause = "WHERE " + whereClause
	}
	whereArgs := QB.ValueList(condition)
	args = append(args, whereArgs...)

	query := fmt.Sprintf(`UPDATE %s SET %s %s`, condition.TName(), valueClause, whereClause)
	// Rebind into native sql driver parameter syntan
	query = repo.tx.Rebind(query)
	result, err := repo.tx.Exec(query, args...)
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Class().Name() == "integrity_constraint_violation" {
		err = apperror.ErrConstraintViolation
	}
	if err == nil {
		affected, err_ := result.RowsAffected()
		if err_ == nil {
			return affected, nil
		} else {
			logger.Error("UpdateAll", zap.String("query", query), zap.Any("sample", condition), zap.Any("value", value), zap.Error(err_))
			return 0, nil
		}
	} else {
		logger.Error("UpdateAll", zap.String("query", query), zap.Any("entity", condition), zap.Any("value", value), zap.Error(err))
		return 0, err
	}
}

func (repo *DefaultCrudRepository[E]) UpdateNull(condition E, value E) (int64, error) {
	valueClause := QB.NamedSetNullList(value)
	if len(valueClause) == 0 {
		return 0, errors.New("NoValues")
	}

	whereClause := QB.WhereEqual(condition)
	if len(whereClause) > 0 {
		whereClause = "WHERE " + whereClause
	}

	query := fmt.Sprintf(`UPDATE %s SET %s %s`, condition.TName(), valueClause, whereClause)
	result, err := repo.tx.NamedExec(query, condition)
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Class().Name() == "integrity_constraint_violation" {
		err = apperror.ErrConstraintViolation
	}
	if err == nil {
		affected, err_ := result.RowsAffected()
		if err_ == nil {
			return affected, nil
		} else {
			logger.Error("UpdateAll", zap.String("query", query), zap.Any("sample", condition), zap.Any("value", value), zap.Error(err_))
			return 0, nil
		}
	} else {
		logger.Error("UpdateAll", zap.String("query", query), zap.Any("entity", condition), zap.Any("value", value), zap.Error(err))
		return 0, err
	}
}

func (repo *DefaultCrudRepository[E]) DeleteSome(condition E) (int64, error) {
	qb := QueryBuilder{}
	rule := NewRule(condition, condition, REL_EQUAL)
	qb.Select(condition).Where(rule)
	query, qparam := qb.BuildDelete()

	query = repo.tx.Rebind(query)
	result, err := repo.tx.Exec(query, qparam...)

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Class().Name() == "integrity_constraint_violation" {
		err = apperror.ErrConstraintViolation
	}
	if err == nil {
		affected, err_ := result.RowsAffected()
		if err_ != nil {
			return affected, nil
		} else {
			return 0, nil
		}
	} else {
		logger.Warn("DeleteEntity", zap.Any("conidition", condition))
		return 0, err
	}
}
