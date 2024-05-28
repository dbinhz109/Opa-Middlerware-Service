// Code generated from dbml table
package model

import (
	"github.com/google/uuid"
)

// ExternalService is generated type for table 'external_service'
type ExternalService struct {
	Id        *uuid.UUID `json:"id,omitempty" db:"id" mapstructure:"id"`
	Name      *string    `json:"name,omitempty" db:"name" mapstructure:"name"`
	AccessKey *string    `json:"access_key,omitempty" db:"access_key" mapstructure:"access_key"`
	SecretKey *string    `json:"secret_key,omitempty" db:"secret_key" mapstructure:"secret_key"`
}

// table 'ExternalService' columns list struct
type __tblExternalServiceColumns struct {
	Id        string
	Name      string
	AccessKey string
	SecretKey string
}

// table 'external_service' metadata struct
type __tblExternalService struct {
	Name_       string
	Columns_    __tblExternalServiceColumns
	AllColumns_ []string
}

// table 'external_service' metadata info
var _tblExternalService = __tblExternalService{
	Name_: "public.external_service",
	Columns_: __tblExternalServiceColumns{
		Id:        "id",
		Name:      "name",
		AccessKey: "access_key",
		SecretKey: "secret_key",
	},
	AllColumns_: []string{"id", "name", "access_key", "secret_key"},
}

func (e *ExternalService) TName() string {
	return e.T().Name()
}
func (e *ExternalService) CreateInstance() any {
	return &ExternalService{}
}
func (e *ExternalService) PK() any {
	return e.Id
}
func (e *ExternalService) SetPK(k any) {
	e.Id = k.(*uuid.UUID)
}

// ForEach iterate all fields to call given callback to that field value of target object
// target is typecasted to keep common interface entity independent
func (e *ExternalService) ForEach(target any, shared any, f func(fieldName string, fieldValue any, sharedValue any)) {
	tgt := target.(*ExternalService)
	c := e.T().Columns()
	// Convert typed nil pointers to untyped
	if e.Id != nil {
		f(c.Id, tgt.Id, shared)
	}
	if e.Name != nil {
		f(c.Name, tgt.Name, shared)
	}
	if e.AccessKey != nil {
		f(c.AccessKey, tgt.AccessKey, shared)
	}
	if e.SecretKey != nil {
		f(c.SecretKey, tgt.SecretKey, shared)
	}
}

// T return metadata info for table 'external_service'
func (*ExternalService) T() *__tblExternalService {
	return &_tblExternalService
}
func (m *__tblExternalService) Name() string {
	return m.Name_
}
func (m *__tblExternalService) Columns() *__tblExternalServiceColumns {
	return &m.Columns_
}
func (m *__tblExternalService) AllColumns() []string {
	return m.AllColumns_
}
