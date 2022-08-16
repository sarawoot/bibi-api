package pgx

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	pgxstdlib "github.com/jackc/pgx/v4/stdlib"
)

func RunMigrator(migrationsSource string, postgresConnection string, commands []string) error {
	pgxConfig, err := pgx.ParseConfig(postgresConnection)
	if err != nil {
		return err
	}

	migrateDB, err := migratepgx.WithInstance(pgxstdlib.OpenDB(*pgxConfig), &migratepgx.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsSource, "", migrateDB)
	if err != nil {
		return err
	}

	if len(commands) < 1 {
		return errors.New("invalid action")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for sig := range sigChan {
			if sig == syscall.SIGTERM || sig == syscall.SIGINT {
				fmt.Println("Gracefully stopping migration")
				m.GracefulStop <- true
				break
			}
		}
	}()

	switch commands[0] {
	case "up":
		if err := m.Up(); err != nil {
			return err
		}

	case "down":
		if err := m.Down(); err != nil {
			return err
		}

	case "steps":
		if len(commands) < 2 {
			return errors.New("no steps provided")
		}
		steps, err := strconv.ParseInt(commands[1], 10, 64)
		if err != nil {
			return err
		}
		if err := m.Steps(int(steps)); err != nil {
			return err
		}

	case "goto":
		if len(commands) < 2 {
			return errors.New("no version provided")
		}
		version, err := strconv.ParseUint(commands[1], 10, 64)
		if err != nil {
			return err
		}
		if err := m.Migrate(uint(version)); err != nil {
			return err
		}

	default:
		return errors.New("invalid action")
	}

	return nil
}
