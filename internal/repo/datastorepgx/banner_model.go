package datastorepgx

import (
	"api/internal/model"

	"github.com/jackc/pgtype"
)

type Banner struct {
	ID          pgtype.UUID `db:"id"`
	Name        string      `db:"name"`
	AreaCode    string      `db:"area_code"`
	CreatedTime pgtype.Time `db:"created_time"`

	Images []BannerImage `db:"-"`
}

type BannerImage struct {
	ID          pgtype.UUID `db:"id"`
	BannerID    pgtype.UUID `db:"banner_id"`
	Path        string      `db:"path"`
	CreatedTime pgtype.Time `db:"created_time"`
}

func (b Banner) toModel() model.Banner {
	rs := model.Banner{
		ID:       b.ID.Bytes,
		Name:     b.Name,
		AreaCode: b.AreaCode,
	}

	rs.Images = make([]model.BannerImage, 0, len(b.Images))
	for _, image := range b.Images {
		rs.Images = append(rs.Images, image.toModel())
	}

	return rs
}

func (b BannerImage) toModel() model.BannerImage {
	return model.BannerImage{
		ID:       b.ID.Bytes,
		BannerID: b.BannerID.Bytes,
		Path:     b.Path,
	}
}
