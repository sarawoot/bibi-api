package main

import (
	"fmt"
	"os"

	"api/internal/config"
	migrator "api/pkg/migrators/pgx"
)

const migrationsSource = "file://third_party/db/migration"

func main() {
	if err := migrator.RunMigrator(migrationsSource, config.App.PostgresConnection, os.Args[1:]); err != nil {
		fmt.Println("Run migrator failed!", err)
		os.Exit(1)
	}
}
