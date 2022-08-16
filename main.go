package main

import (
	"api/internal/repo/datastorepgx"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func main() {
	con, err := initPgx("postgres://postgres:P%40ssw0rd@localhost:5432/bibi")
	if err != nil {
		panic(err)
	}

	r := datastorepgx.New(con)

	fmt.Println(r.AdminCountProduct(context.Background()))

}
func initPgx(postgresConnection string) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(postgresConnection)
	if err != nil {
		return nil, err
	}

	looger := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.JSONFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}

	conf.ConnConfig.Logger = logrusadapter.NewLogger(looger)

	return pgxpool.ConnectConfig(context.Background(), conf)
}
