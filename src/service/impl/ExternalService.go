package impl

import (
	"go-app/src/database/repository"
	"go-app/src/model"
	"go-app/src/service/spec"

	"github.com/jmoiron/sqlx"
)

type externalService struct {
	externalServiceCrud *DefaultCrudService[*model.ExternalService]
}

func NewExternalService() *externalService {
	endUserRepoFactory := func(tx *sqlx.Tx) spec.IDefaultCrudRepository[*model.ExternalService] {
		return repository.NewExternalServiceRepository(tx)
	}
	svc := &externalService{
		externalServiceCrud: NewDefaultCrudService(endUserRepoFactory),
	}
	return svc
}

func (s *externalService) ExternalServiceCrud() spec.IDefaultCrudService[*model.ExternalService] {
	return s.externalServiceCrud
}
