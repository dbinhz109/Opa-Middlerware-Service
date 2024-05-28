package configuration

import (
	"go-app/src/service/spec"

	"github.com/gin-gonic/gin"
)

// AppContext will hold reference to service interfaces that will eventually be injected into this layer on initialization
type AppContext struct {
	R               *gin.Engine
	ConsulService   spec.IConsulService
	RetailsService  spec.IRetailsService
	ExternalService spec.IExternalService

	RetailsControlHandler spec.IRetailsController
}
