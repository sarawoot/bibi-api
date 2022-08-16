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

func (r *DataStoreRepo) CreateBanner(ctx context.Context, banner *model.Banner) error {
	return r.pool.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(
			ctx,
			"SET CONSTRAINTS ALL DEFERRED;",
		)
		if err != nil {
			return err
		}

		id := banner.ID
		if id == uuid.Nil {
			id = uuid.New()
		}
		banner.ID = id
		bannerID := pgtype.UUID{Bytes: id, Status: pgtype.Present}

		_, err = tx.Exec(
			ctx,
			"INSERT INTO "+bannerTable+"(id,name,area_code) VALUES ($1,$2,$3);",
			bannerID, banner.Name, banner.AreaCode,
		)
		if err != nil {
			return err
		}

		for _, image := range banner.Images {
			id := image.ID
			if id == uuid.Nil {
				id = uuid.New()
			}
			image.ID = id
			imageID := pgtype.UUID{Bytes: id, Status: pgtype.Present}

			_, err := tx.Exec(
				ctx,
				"INSERT INTO "+bannerImageTable+"(id,banner_id,path) VALUES ($1,$2,$3);",
				imageID, bannerID, image.Path,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *DataStoreRepo) CreateBannerImage(ctx context.Context, bannerID uuid.UUID, path string) (uuid.UUID, error) {
	if bannerID == uuid.Nil {
		return uuid.Nil, errors.New("id is blank")
	}

	id := uuid.New()
	_, err := r.pool.Exec(
		ctx,
		"INSERT INTO "+bannerImageTable+"(id,banner_id,path) VALUES ($1,$2,$3);",
		pgtype.UUID{Bytes: id, Status: pgtype.Present},
		pgtype.UUID{Bytes: bannerID, Status: pgtype.Present},
		path,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *DataStoreRepo) UpdateBanner(ctx context.Context, id uuid.UUID, banner *model.Banner) error {
	if id == uuid.Nil {
		return errors.New("id is blank")
	}

	bannerID := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	cmdTag, err := r.pool.Exec(
		ctx,
		"UPDATE "+bannerTable+" SET name=$1,area_code=$2 WHERE id=$3",
		banner.Name, banner.AreaCode, bannerID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *DataStoreRepo) DeleteBannerByID(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is blank")
	}

	bannerID := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	cmdTag, err := r.pool.Exec(
		ctx,
		"DELETE FROM "+bannerTable+" WHERE id = $1",
		bannerID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *DataStoreRepo) DeleteBannerImageByID(ctx context.Context, bannerID, bannerImageID uuid.UUID) error {
	if bannerID == uuid.Nil || bannerImageID == uuid.Nil {
		return errors.New("banner_id or banner_image_id is blank")
	}

	bannerUUID := pgtype.UUID{Bytes: bannerID, Status: pgtype.Present}
	bannerImageUUID := pgtype.UUID{Bytes: bannerImageID, Status: pgtype.Present}
	cmdTag, err := r.pool.Exec(
		ctx,
		"DELETE FROM "+bannerImageTable+" WHERE banner_id = $1 AND id = $2",
		bannerUUID, bannerImageUUID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *DataStoreRepo) ListBanner(ctx context.Context) ([]model.Banner, error) {
	rows, err := r.pool.Query(ctx, "SELECT * FROM "+bannerTable+" ORDER BY created_time asc")
	if err != nil {
		return nil, err
	}

	var bannerRows []Banner
	if err := pgxscan.ScanAll(&bannerRows, rows); err != nil {
		return nil, err
	}

	banners := make([]model.Banner, 0, len(bannerRows))
	for _, row := range bannerRows {
		images, err := r.listBannerImageByBannerID(ctx, row.ID.Bytes)
		if err != nil {
			return nil, err
		}
		banner := *row.toModel()

		banner.Images = make([]model.BannerImage, 0, len(images))
		for _, image := range images {
			banner.Images = append(banner.Images, *image.toModel())
		}

		banners = append(banners, banner)
	}

	return banners, nil
}

func (r *DataStoreRepo) GetBannerByID(ctx context.Context, id uuid.UUID) (*model.Banner, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is blank")
	}

	uid := pgtype.UUID{Bytes: id, Status: pgtype.Present}
	row, err := r.pool.Query(ctx, "SELECT * FROM "+bannerTable+" WHERE id = $1 limit 1", uid)
	if err != nil {
		return nil, err
	}

	var bannerRow Banner
	if err := pgxscan.ScanOne(&bannerRow, row); err != nil {
		return nil, err
	}

	images, err := r.listBannerImageByBannerID(ctx, id)
	if err != nil {
		return nil, err
	}
	bannerRow.Images = images

	return bannerRow.toModel(), nil
}

func (r *DataStoreRepo) GetBannerByAreaCode(ctx context.Context, areaCode string) (*model.Banner, error) {
	row, err := r.pool.Query(ctx, "SELECT * FROM "+bannerTable+" WHERE area_code = $1 limit 1", areaCode)
	if err != nil {
		return nil, err
	}

	var bannerRow Banner
	if err := pgxscan.ScanOne(&bannerRow, row); err != nil {
		return nil, err
	}

	images, err := r.listBannerImageByBannerID(ctx, bannerRow.ID.Bytes)
	if err != nil {
		return nil, err
	}
	bannerRow.Images = images

	return bannerRow.toModel(), nil
}

func (r *DataStoreRepo) listBannerImageByBannerID(ctx context.Context, id uuid.UUID) ([]BannerImage, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is blank")
	}
	uid := pgtype.UUID{Bytes: id, Status: pgtype.Present}

	rows, err := r.pool.Query(ctx, "SELECT * FROM "+bannerImageTable+" WHERE banner_id=$1 ORDER BY created_time asc", uid)
	if err != nil {
		return nil, err
	}

	var bannerImageRows []BannerImage
	if err := pgxscan.ScanAll(&bannerImageRows, rows); err != nil {
		return nil, err
	}

	return bannerImageRows, nil
}
