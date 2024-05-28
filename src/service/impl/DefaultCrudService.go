package impl

import (
	"context"
	"go-app/src/database"
	"go-app/src/logger"
	"go-app/src/model"
	"go-app/src/service"
	"go-app/src/service/spec"
	"go-app/src/service/spec/apperror"
	"reflect"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DefaultCrudService[T model.EntityModel[any]] struct {
	DefaultInstance   T
	RepositoryFactory spec.IRepositoryFactory[T]
}

func NewDefaultCrudService[T model.EntityModel[any]](repoFactory spec.IRepositoryFactory[T]) *DefaultCrudService[T] {
	return &DefaultCrudService[T]{
		RepositoryFactory: repoFactory,
	}
}
func (svc *DefaultCrudService[T]) validatePrimaryKey(entity T) {
	pk := entity.PK()
	// do not check nil pointer here as generic type always have (pk != nil) true
	supportedPKType := false
	_, ok := pk.(*uuid.UUID)
	if ok {
		supportedPKType = true
		// generate uuid primary key
		newId, _ := uuid.NewRandom()
		entity.SetPK(&newId)
	}
	int64Val, ok := pk.(*int64)
	if ok {
		supportedPKType = true
		if int64Val != nil {
			logger.Warn("CreateRequiresEmptyPK", zap.Any("typeofpk", reflect.TypeOf(pk)), zap.Any("primaryKey", pk))
			panic(apperror.ErrNotImplemented)
		}
	}
	if !supportedPKType {
		logger.Warn("UnsupportedPKType", zap.Any("typeofpk", reflect.TypeOf(pk)))
		panic(apperror.ErrUnsupportedDataType)
	}
}
func (svc *DefaultCrudService[T]) Create(context context.Context, entity T) (ret T, retErr error) {
	ret = svc.DefaultInstance
	tx, retErr := database.GetDbTransaction()
	if retErr != nil || tx == nil {
		logger.Warn("TxFailure", zap.Error(retErr))
		return ret, retErr
	}
	defer service.PanicRecoverTx(tx)
	svc.validatePrimaryKey(entity)
	repo := svc.RepositoryFactory(tx)
	ret, retErr = repo.Create(entity)
	if retErr != nil {
		logger.Warn("CreateFailed", zap.Any("entity", entity))
		tx.Rollback()
		return ret, retErr
	}
	tx.Commit()
	return ret, retErr
}

func (svc *DefaultCrudService[T]) GetOne(context context.Context, sample T) (ret T, retErr error) {
	ret = svc.DefaultInstance
	tx, retErr := database.GetDbTransaction()
	if retErr != nil || tx == nil {
		logger.Warn("TxFailure", zap.Error(retErr))
		return ret, retErr
	}
	defer service.PanicRecoverTx(tx)
	defer tx.Rollback()

	repo := svc.RepositoryFactory(tx)
	return repo.GetOne(sample)
}

func (svc *DefaultCrudService[T]) GetList(context context.Context, sample T, param *model.ListParam) (ret []T, total int64, retErr error) {
	ret = nil
	tx, retErr := database.GetDbTransaction()
	if retErr != nil || tx == nil {
		logger.Warn("TxFailure", zap.Error(retErr))
		return ret, 0, retErr
	}
	defer service.PanicRecoverTx(tx)
	defer tx.Rollback()

	repo := svc.RepositoryFactory(tx)
	return repo.GetList(sample, param)
}

func (svc *DefaultCrudService[T]) Update(context context.Context, sample T, value T) (retErr error) {
	tx, retErr := database.GetDbTransaction()
	if retErr != nil || tx == nil {
		logger.Warn("TxFailure", zap.Error(retErr))
		return retErr
	}
	defer service.PanicRecoverTx(tx)

	repo := svc.RepositoryFactory(tx)
	_, retErr = repo.UpdateAll(sample, value)
	if retErr != nil {
		logger.Warn("UpdateFailure", zap.Error(retErr), zap.Any("sample", sample), zap.Any("value", value))
		tx.Rollback()
		return retErr
	}
	tx.Commit()
	return retErr
}

func (svc *DefaultCrudService[T]) UpdateNull(context context.Context, sample T, value T) (retErr error) {
	tx, retErr := database.GetDbTransaction()
	if retErr != nil || tx == nil {
		logger.Warn("TxFailure", zap.Error(retErr))
		return retErr
	}
	defer service.PanicRecoverTx(tx)

	repo := svc.RepositoryFactory(tx)
	_, retErr = repo.UpdateNull(sample, value)
	if retErr != nil {
		logger.Warn("UpdateFailure", zap.Error(retErr), zap.Any("sample", sample), zap.Any("value", value))
		tx.Rollback()
		return retErr
	}
	tx.Commit()
	return retErr
}

func (svc *DefaultCrudService[T]) Delete(context context.Context, sample T) (retErr error) {
	tx, retErr := database.GetDbTransaction()
	if retErr != nil || tx == nil {
		logger.Warn("TxFailure", zap.Error(retErr))
		return retErr
	}
	defer service.PanicRecoverTx(tx)

	repo := svc.RepositoryFactory(tx)
	_, retErr = repo.DeleteSome(sample)
	if retErr != nil {
		logger.Warn("DeleteFailure", zap.Error(retErr), zap.Any("sample", sample))
		tx.Rollback()
		return retErr
	}
	tx.Commit()
	return nil
}
