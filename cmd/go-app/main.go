package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.elastic.co/apm/module/apmgin/v2"
	"go.uber.org/zap"

	"go-app/src/app"
	"go-app/src/configuration"
	"go-app/src/logger"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "go-app",
	Short: "short description",
	Long:  `long description`,
	Run:   rootCmdRun,
}

func rootCmdRun(cmd *cobra.Command, args []string) {
	// configuration.InitializeAppConfig()
	appContext := app.GetAppContext()

	engine := gin.Default()
	engine.Use(apmgin.Middleware(engine))
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	engine.Use(cors.New(corsConfig)) //enable cors default to all hosts

	// engine.Use(gzip.Gzip(gzip.DefaultCompression))

	appContext.R = engine
	app.InitRouting(appContext)
	engine.GET("/healthz", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain", []byte(viper.GetString(configuration.CFG_APPLICATION_NAME)))
	})

	srv := &http.Server{
		Addr:           viper.GetString(configuration.CFG_SERVER_HOST) + ":" + viper.GetString(configuration.CFG_SERVER_PORT),
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to initialize server: %v\n", zap.Error(err))
		}
	}()

	logger.Info("Listening on port %v\n", zap.String("address", srv.Addr))

	// Wait for kill signal of channel
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}

func init() {
	cobra.OnInitialize(configuration.InitializeAppConfig)
	rootCmd.AddCommand(migrateCmd)
}
