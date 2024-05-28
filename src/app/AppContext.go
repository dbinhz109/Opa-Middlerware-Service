package app

import (
	"go-app/src/configuration"
	"go-app/src/controller"
	"go-app/src/logger"
	"go-app/src/service/impl"
	"time"

	"sync"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	ctx     *configuration.AppContext
	ctxOnce sync.Once
)

func GetAppContext() *configuration.AppContext {
	ctxOnce.Do(func() {
		logger.Info("Initializing ...")
		ctx = &configuration.AppContext{}
		InitializeAppContext(ctx)
		RunScheduledTasks()
		RegisterWithServiceRegistry(ctx)
	})
	return ctx
}

func RunScheduledTasks() {
	s := gocron.NewScheduler(time.UTC)
	s.StartAsync()
}

// InitializeAppContext initialize app components like services and controllers, resolving their dependencies properly
func InitializeAppContext(ctx *configuration.AppContext) {
	viper.SetDefault("security.auth.jwt.ttl", 1800)
	ctx.ConsulService = impl.NewConsulService()
	ctx.RetailsService = impl.NewRetailsService()
	ctx.ExternalService = impl.NewExternalService()

	ctx.RetailsControlHandler = controller.NewRetailsController(ctx.RetailsService)

}

func RegisterWithServiceRegistry(ctx *configuration.AppContext) {
	if ctx.ConsulService != nil {
		appName := viper.GetString("application.name")
		appInstance := viper.GetString("application.instanceId")
		logger.Info("appInstance", zap.Any("appInstance", appInstance))
		// regAddress := viper.GetString("server.publicAddress")
		regPort := viper.GetString("server.publicPort")
		serverAddress := viper.GetString("server.ip")
		ctx.ConsulService.RegisterService(appName, appInstance, serverAddress, regPort)
	}
}
