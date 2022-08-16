package datastorepgx

import (
	"api/internal/model"

	"github.com/jackc/pgtype"
)

type ProductRecommend struct {
	ID          pgtype.UUID `db:"id"`
	ProductID   pgtype.UUID `db:"product_id"`
	Product     ProductScan `db:"product"`
	CreatedTime pgtype.Time `db:"created_time"`
}

func (p *ProductRecommend) toModel() *model.ProductRecommend {
	rs := model.ProductRecommend{
		ID:        p.ID.Bytes,
		ProductID: p.ProductID.Bytes,
		Product:   *p.Product.toModel(),
	}

	return &rs
}
