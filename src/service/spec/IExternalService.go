package spec

import (
	"go-app/src/model"
)

// IGroupEndUserService ...
type IExternalService interface {
	ExternalServiceCrud() IDefaultCrudService[*model.ExternalService]
}
