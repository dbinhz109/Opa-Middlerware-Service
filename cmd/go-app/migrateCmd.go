package main

import (
	"go-app/src/database"
	"go-app/src/service"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "DB migration",
	Long:  "DB migration",
	Run:   migrateCmdRun,
}

func migrateCmdRun(cmd *cobra.Command, args []string) {
	db := database.GetDbInstance()
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	service.PanicOnError(err)
	m, err := migrate.NewWithDatabaseInstance("file:///app/public/migrations", "postgres", driver)
	service.PanicOnError(err)
	if len(args) > 0 {
		if args[0] == "up" {
			m.Up()
		}
		if args[0] == "down" {
			m.Down()
		}
	}
}
