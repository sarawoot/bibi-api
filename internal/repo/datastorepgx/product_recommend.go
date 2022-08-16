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

func (r *DataStoreRepo) AdminListProductRecommend(ctx context.Context, limit, offset int) ([]model.ProductRecommend, error) {
	rows, err := r.pool.Query(
		ctx,
		`SELECT
			product_recommends.*,
			products.id "product.id",
			products.brand "product.brand",
			products.name "product.name",
			products.short_description "product.short_description",
			products.description "product.description",
			products.size "product.size",
			products.price "product.price",
			products.product_type_id "product.product_type_id",
			products.product_category_id "product.product_category_id",
			products.skin_type_id "product.skin_type_id",
			products.country_id "product.country_id",
			products.tags "product.tags",
			products.created_time "product.created_time",
			products.deleted_time "product.deleted_time"
		FROM `+productRecommendTable+","+productTable+
			" WHERE products.id = product_recommends.product_id AND products.deleted_time IS NULL"+
			" ORDER BY product_recommends.created_time ASC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}

	var productRecommendRows []ProductRecommend
	if err := pgxscan.ScanAll(&productRecommendRows, rows); err != nil {
		return nil, err
	}

	productRecommends := make([]model.ProductRecommend, 0, len(productRecommendRows))
	for _, row := range productRecommendRows {
		productRecommends = append(productRecommends, *row.toModel())
	}

	return productRecommends, nil
}

func (r *DataStoreRepo) AdminCountProductRecommend(ctx context.Context) (int, error) {
	rows, err := r.pool.Query(ctx, "SELECT COUNT(id) AS cnt FROM "+productRecommendTable)
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

func (r *DataStoreRepo) CreateProductRecommend(ctx context.Context, productID uuid.UUID) (uuid.UUID, error) {
	if productID == uuid.Nil {
		return uuid.Nil, errors.New("id is blank")
	}

	id := uuid.New()
	_, err := r.pool.Exec(
		ctx,
		"INSERT INTO "+productRecommendTable+"(id,product_id) VALUES ($1,$2);",
		pgtype.UUID{Bytes: id, Status: pgtype.Present},
		pgtype.UUID{Bytes: productID, Status: pgtype.Present},
	)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *DataStoreRepo) DeleteProductRecommendByID(ctx context.Context, productRecommendID uuid.UUID) error {
	if productRecommendID == uuid.Nil {
		return errors.New("product_recommend_id is blank")
	}

	cmdTag, err := r.pool.Exec(
		ctx,
		"DELETE FROM "+productRecommendTable+" WHERE id = $1",
		pgtype.UUID{Bytes: productRecommendID, Status: pgtype.Present},
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *DataStoreRepo) MobileListProductRecommend(ctx context.Context, limit int) ([]model.ProductRecommend, error) {
	rows, err := r.pool.Query(
		ctx,
		`SELECT
			product_recommends.*,
			products.id "product.id",
			products.brand "product.brand",
			products.name "product.name",
			products.short_description "product.short_description",
			products.description "product.description",
			products.size "product.size",
			products.price "product.price",
			products.product_type_id "product.product_type_id",
			products.product_category_id "product.product_category_id",
			products.skin_type_id "product.skin_type_id",
			products.country_id "product.country_id",
			products.tags "product.tags",
			products.created_time "product.created_time",
			products.deleted_time "product.deleted_time"
		FROM `+productRecommendTable+","+productTable+
			" WHERE products.id = product_recommends.product_id AND products.deleted_time IS NULL"+
			" ORDER BY product_recommends.created_time ASC LIMIT $1",
		limit,
	)
	if err != nil {
		return nil, err
	}

	var productRecommendRows []ProductRecommend
	if err := pgxscan.ScanAll(&productRecommendRows, rows); err != nil {
		return nil, err
	}

	productRecommends := make([]model.ProductRecommend, 0, len(productRecommendRows))
	for _, row := range productRecommendRows {
		images, err := r.listProductImageByProductID(ctx, row.ProductID.Bytes)
		if err != nil {
			return nil, err
		}
		row.Product.Images = images

		productRecommends = append(productRecommends, *row.toModel())
	}

	return productRecommends, nil
}
