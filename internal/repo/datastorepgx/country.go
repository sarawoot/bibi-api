package datastorepgx

import (
	"api/internal/model"
	"context"
	"errors"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

func (r *DataStoreRepo) ListCountry(ctx context.Context) ([]model.Country, error) {
	rows, err := r.pool.Query(ctx, "SELECT * FROM "+countryTable+" WHERE deleted_time IS NULL ORDER BY created_time asc")
	if err != nil {
		return nil, err
	}

	var countryRows []Country
	if err := pgxscan.ScanAll(&countryRows, rows); err != nil {
		return nil, err
	}

	countries := make([]model.Country, 0, len(countryRows))
	for _, row := range countryRows {
		countries = append(countries, *row.toModel())
	}

	return countries, nil
}

func (r *DataStoreRepo) GetCountryByID(ctx context.Context, id uuid.UUID) (*model.Country, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is blank")
	}

	uid := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	row, err := r.pool.Query(ctx, "SELECT * FROM "+countryTable+" WHERE id = $1 limit 1", uid)
	if err != nil {
		return nil, err
	}

	var countryRow Country
	if err := pgxscan.ScanOne(&countryRow, row); err != nil {
		return nil, err
	}

	return countryRow.toModel(), nil
}

func (r *DataStoreRepo) CreateCountry(ctx context.Context, country *model.Country) error {
	id := country.ID
	if id == uuid.Nil {
		id = uuid.New()
	}
	country.ID = id

	countryID := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	_, err := r.pool.Exec(
		ctx,
		"INSERT INTO "+countryTable+"(id,name) VALUES ($1,$2)",
		countryID, country.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DataStoreRepo) UpdateCountry(ctx context.Context, id uuid.UUID, country *model.Country) error {
	if id == uuid.Nil {
		return errors.New("id is blank")
	}

	countryID := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	cmdTag, err := r.pool.Exec(
		ctx,
		"UPDATE "+countryTable+" SET name = $1 WHERE id = $2",
		country.Name, countryID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *DataStoreRepo) DeleteCountryByID(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is blank")
	}

	countryID := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	cmdTag, err := r.pool.Exec(
		ctx,
		"UPDATE "+countryTable+" SET deleted_time = CURRENT_TIMESTAMP WHERE id = $1",
		countryID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
