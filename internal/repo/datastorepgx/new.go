package datastorepgx

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type DataStoreRepo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *DataStoreRepo {
	return &DataStoreRepo{
		pool: pool,
	}
}
