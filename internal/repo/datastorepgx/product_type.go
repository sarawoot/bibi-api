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

func (r *DataStoreRepo) ListProductType(ctx context.Context) ([]model.ProductType, error) {
	rows, err := r.pool.Query(ctx, "SELECT id,name FROM "+productTypeTable+" WHERE deleted_time IS NULL ORDER BY created_time asc")
	if err != nil {
		return nil, err
	}

	var productTypeRows []ProductType
	if err := pgxscan.ScanAll(&productTypeRows, rows); err != nil {
		return nil, err
	}

	productTypes := make([]model.ProductType, 0, len(productTypeRows))
	for _, row := range productTypeRows {
		productTypes = append(productTypes, row.toModel())
	}

	return productTypes, nil
}

func (r *DataStoreRepo) GetProductTypeByID(ctx context.Context, id uuid.UUID) (*model.ProductType, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is blank")
	}

	uid := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	row, err := r.pool.Query(ctx, "SELECT id,name FROM "+productTypeTable+" WHERE id = $1 limit 1", uid)
	if err != nil {
		return nil, err
	}

	var productTypeRow ProductType
	if err := pgxscan.ScanOne(&productTypeRow, row); err != nil {
		return nil, err
	}
	productType := productTypeRow.toModel()

	return &productType, nil
}

func (r *DataStoreRepo) CreateProductType(ctx context.Context, productType *model.ProductType) error {
	id := productType.ID
	if id == uuid.Nil {
		id = uuid.New()
	}
	productType.ID = id

	productTypeID := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	_, err := r.pool.Exec(
		ctx,
		"INSERT INTO "+productTypeTable+"(id,name) VALUES ($1,$2)",
		productTypeID, productType.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DataStoreRepo) UpdateProductType(ctx context.Context, id uuid.UUID, productType *model.ProductType) error {
	if id == uuid.Nil {
		return errors.New("id is blank")
	}

	productTypeID := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	cmdTag, err := r.pool.Exec(
		ctx,
		"UPDATE "+productTypeTable+" SET name = $1 WHERE id = $2",
		productType.Name, productTypeID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *DataStoreRepo) DeleteProductTypeByID(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is blank")
	}

	productTypeID := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	cmdTag, err := r.pool.Exec(
		ctx,
		"UPDATE "+productTypeTable+" SET deleted_time = CURRENT_TIMESTAMP WHERE id = $1",
		productTypeID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
