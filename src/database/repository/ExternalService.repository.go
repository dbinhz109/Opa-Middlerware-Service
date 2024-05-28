// Generated automatically from dbml Table entry
package repository

import (
	"go-app/src/model"

	"github.com/jmoiron/sqlx"
)

type externalServiceRepository struct {
	DefaultCrudRepository[*model.ExternalService]
}

// Create new repository instance on per transaction basis
func NewExternalServiceRepository(tx *sqlx.Tx) *externalServiceRepository {
	return &externalServiceRepository{
		DefaultCrudRepository: *NewDefaultCrudRepository[*model.ExternalService](tx),
	}
}
