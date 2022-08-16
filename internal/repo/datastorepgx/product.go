package datastorepgx

import (
	"api/internal/model"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

func (r *DataStoreRepo) CreateProduct(ctx context.Context, product *model.Product) error {
	return r.pool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(
			ctx,
			"SET CONSTRAINTS ALL DEFERRED;",
		)
		if err != nil {
			return err
		}

		row := productModelToRow(*product)
		if row.ID.Status == pgtype.Null {
			row.ID = pgtype.UUID{Bytes: uuid.New(), Status: pgtype.Present}
		}

		columns := []string{"id", "brand", "name", "short_description", "description", "size", "price",
			"product_type_id", "product_category_id", "skin_type_id", "country_id", "tags"}
		values := []interface{}{row.ID, row.Brand, row.Name, row.ShortDescription, row.Description, row.Size, row.Price,
			row.ProductTypeID, row.ProductCategoryID, row.SkinTypeID, row.CountryID, row.Tags}

		columnInsert := ""
		valueInsert := ""
		paramInsert := make([]interface{}, 0, len(values))
		i := 0
		for idx, v := range values {
			if !isNil(v) {
				i += 1
				columnInsert += "," + columns[idx]
				valueInsert += fmt.Sprintf(",$%d", i)
				paramInsert = append(paramInsert, v)
			}
		}

		_, err = tx.Exec(
			ctx,
			"INSERT INTO "+productTable+" ("+strings.TrimLeft(columnInsert, ",")+") "+
				"VALUES ("+strings.TrimLeft(valueInsert, ",")+");",
			paramInsert...,
		)
		if err != nil {
			return err
		}

		for _, image := range product.Images {
			rowImage := productImageModelToRow(image)
			if rowImage.ID.Status == pgtype.Null {
				rowImage.ID = pgtype.UUID{Bytes: uuid.New(), Status: pgtype.Present}
			}
			rowImage.ProductID = row.ID

			_, err := tx.Exec(
				ctx,
				"INSERT INTO "+productImageTable+" (id,product_id,path) VALUES ($1,$2,$3);",
				rowImage.ID, rowImage.ProductID, rowImage.Path,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *DataStoreRepo) CreateProductImage(ctx context.Context, productID uuid.UUID, path string) (uuid.UUID, error) {
	if productID == uuid.Nil {
		return uuid.Nil, errors.New("id is blank")
	}

	id := uuid.New()
	_, err := r.pool.Exec(
		ctx,
		"INSERT INTO "+productImageTable+"(id,product_id,path) VALUES ($1,$2,$3);",
		pgtype.UUID{Bytes: id, Status: pgtype.Present},
		pgtype.UUID{Bytes: productID, Status: pgtype.Present},
		path,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *DataStoreRepo) UpdateProduct(ctx context.Context, id uuid.UUID, product *model.Product) error {
	if id == uuid.Nil {
		return errors.New("id is blank")
	}
	row := productModelToRow(*product)
	columns := []string{"brand", "name", "short_description", "description", "size", "price",
		"product_type_id", "product_category_id", "skin_type_id", "country_id", "tags"}
	values := []interface{}{row.Brand, row.Name, row.ShortDescription, row.Description, row.Size, row.Price,
		row.ProductTypeID, row.ProductCategoryID, row.SkinTypeID, row.CountryID, row.Tags}

	columnUpdate := ""
	paramInsert := make([]interface{}, 0, len(values)+1)
	i := 0
	for idx, v := range values {
		if !isNil(v) {
			i += 1
			columnUpdate += fmt.Sprintf(",%s=$%d", columns[idx], i)
			paramInsert = append(paramInsert, v)
		}
	}
	paramInsert = append(paramInsert, pgtype.UUID{Bytes: id, Status: pgtype.Present})

	cmdTag, err := r.pool.Exec(
		ctx,
		"UPDATE "+productTable+" SET "+strings.TrimLeft(columnUpdate, ",")+fmt.Sprintf(" WHERE id=$%d", len(paramInsert)),
		paramInsert...,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *DataStoreRepo) DeleteProductByID(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is blank")
	}

	cmdTag, err := r.pool.Exec(
		ctx,
		"UPDATE "+productTable+" SET deleted_time = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_time IS NULL",
		pgtype.UUID{Bytes: id, Status: pgtype.Present},
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *DataStoreRepo) DeleteProductImageByID(ctx context.Context, productID, productImageID uuid.UUID) error {
	if productID == uuid.Nil || productImageID == uuid.Nil {
		return errors.New("product_id or product_image_id is blank")
	}

	cmdTag, err := r.pool.Exec(
		ctx,
		"DELETE FROM "+productImageTable+" WHERE product_id = $1 AND id = $2",
		pgtype.UUID{Bytes: productID, Status: pgtype.Present},
		pgtype.UUID{Bytes: productImageID, Status: pgtype.Present},
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *DataStoreRepo) AdminListProduct(ctx context.Context, limit, offset int) ([]model.Product, error) {
	rows, err := r.pool.Query(
		ctx,
		"SELECT * FROM "+productTable+
			" WHERE deleted_time IS NULL ORDER BY created_time asc"+
			" LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}

	var productRows []ProductScan
	if err := pgxscan.ScanAll(&productRows, rows); err != nil {
		return nil, err
	}

	products := make([]model.Product, 0, len(productRows))
	for _, row := range productRows {
		products = append(products, *row.toModel())
	}

	return products, nil
}

func (r *DataStoreRepo) AdminCountProduct(ctx context.Context) (int, error) {
	rows, err := r.pool.Query(ctx, "SELECT COUNT(id) AS cnt FROM "+productTable+" WHERE deleted_time IS NULL")
	if err != nil {
		return 0, err
	}

	var cnt int
	for rows.Next() {
		if err := rows.Scan(&cnt); err != nil {
			return 0, err
		}
	}

	return cnt, rows.Err()
}

func (r *DataStoreRepo) GetProductByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is blank")
	}

	row, err := r.pool.Query(
		ctx,
		"SELECT * FROM "+productTable+" WHERE id = $1 limit 1",
		pgtype.UUID{Bytes: id, Status: pgtype.Present},
	)
	if err != nil {
		return nil, err
	}

	var productRow ProductScan
	if err := pgxscan.ScanOne(&productRow, row); err != nil {
		return nil, err
	}

	images, err := r.listProductImageByProductID(ctx, id)
	if err != nil {
		return nil, err
	}
	productRow.Images = images

	return productRow.toModel(), nil
}

func (r *DataStoreRepo) ListProductNewArrival(ctx context.Context, limit int) ([]model.Product, error) {
	rows, err := r.pool.Query(
		ctx,
		"SELECT * FROM "+productTable+
			" WHERE deleted_time IS NULL ORDER BY created_time desc"+
			" LIMIT $1",
		limit,
	)
	if err != nil {
		return nil, err
	}

	var productRows []ProductScan
	if err := pgxscan.ScanAll(&productRows, rows); err != nil {
		return nil, err
	}

	products := make([]model.Product, 0, len(productRows))
	for _, row := range productRows {
		images, err := r.listProductImageByProductID(ctx, row.ID.Bytes)
		if err != nil {
			return nil, err
		}
		row.Images = images

		products = append(products, *row.toModel())
	}

	return products, nil
}

func (r *DataStoreRepo) listProductImageByProductID(ctx context.Context, id uuid.UUID) ([]ProductImage, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is blank")
	}

	rows, err := r.pool.Query(
		ctx,
		"SELECT * FROM "+productImageTable+
			" WHERE product_id=$1 ORDER BY created_time asc",
		pgtype.UUID{Bytes: id, Status: pgtype.Present},
	)
	if err != nil {
		return nil, err
	}

	var productImageRows []ProductImage
	if err := pgxscan.ScanAll(&productImageRows, rows); err != nil {
		return nil, err
	}

	return productImageRows, nil
}
