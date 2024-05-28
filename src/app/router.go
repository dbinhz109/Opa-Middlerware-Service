package app

import (
	"go-app/src/configuration"
	"go-app/src/controller/middleware"
)

// @host      localhost:8866
// @BasePath  /api/v2

// NewRouter defines api paths
func InitRouting(ctx *configuration.AppContext) {
	apiGroup := ctx.R.Group("/lancs")
	Middleware := middleware.NewOpaMiddlewareFactory().OPAMiddleware()
	apiGroup.POST("test", Middleware, ctx.RetailsControlHandler.CustomerOrderStatus)
	apiGroup.GET("test", Middleware, ctx.RetailsControlHandler.CustomerOrderStatus)
	apiGroup.PUT("test", Middleware, ctx.RetailsControlHandler.CustomerOrderStatus)
}
