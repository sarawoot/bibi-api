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

func (r *DataStoreRepo) ListProductCategory(ctx context.Context) ([]model.ProductCategory, error) {
	rows, err := r.pool.Query(ctx, "SELECT id,name FROM "+productCategoryTable+" WHERE deleted_time IS NULL ORDER BY created_time asc")
	if err != nil {
		return nil, err
	}

	var productCategoryRows []ProductCategory
	if err := pgxscan.ScanAll(&productCategoryRows, rows); err != nil {
		return nil, err
	}

	productCategories := make([]model.ProductCategory, 0, len(productCategoryRows))
	for _, row := range productCategoryRows {
		productCategories = append(productCategories, *row.toModel())
	}

	return productCategories, nil
}

func (r *DataStoreRepo) GetProductCategoryByID(ctx context.Context, id uuid.UUID) (*model.ProductCategory, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is blank")
	}

	uid := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	row, err := r.pool.Query(ctx, "SELECT id,name FROM "+productCategoryTable+" WHERE id = $1 limit 1", uid)
	if err != nil {
		return nil, err
	}

	var productCategoryRow ProductCategory
	if err := pgxscan.ScanOne(&productCategoryRow, row); err != nil {
		return nil, err
	}

	return productCategoryRow.toModel(), nil
}

func (r *DataStoreRepo) CreateProductCategory(ctx context.Context, productCategory *model.ProductCategory) error {
	id := productCategory.ID
	if id == uuid.Nil {
		id = uuid.New()
	}
	productCategory.ID = id

	productCategoryID := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	_, err := r.pool.Exec(
		ctx,
		"INSERT INTO "+productCategoryTable+"(id,name) VALUES ($1,$2)",
		productCategoryID, productCategory.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DataStoreRepo) UpdateProductCategory(ctx context.Context, id uuid.UUID, productCategory *model.ProductCategory) error {
	if id == uuid.Nil {
		return errors.New("id is blank")
	}

	productCategoryID := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	cmdTag, err := r.pool.Exec(
		ctx,
		"UPDATE "+productCategoryTable+" SET name = $1 WHERE id = $2",
		productCategory.Name, productCategoryID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *DataStoreRepo) DeleteProductCategoryByID(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is blank")
	}

	productCategoryID := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	cmdTag, err := r.pool.Exec(
		ctx,
		"UPDATE "+productCategoryTable+" SET deleted_time = CURRENT_TIMESTAMP WHERE id = $1",
		productCategoryID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
