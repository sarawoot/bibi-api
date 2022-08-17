package datastorepgx

import (
	"api/internal/model"

	"github.com/jackc/pgtype"
)

type ProductCategory struct {
	ID   pgtype.UUID `db:"id"`
	Name string      `db:"name"`
}

func (s *ProductCategory) toModel() model.ProductCategory {
	return model.ProductCategory{
		ID:   s.ID.Bytes,
		Name: s.Name,
	}
}

// func productCategoryModelToRow(p model.ProductCategory) ProductCategory {
// 	row := ProductCategory{
// 		ID: pgtype.UUID{
// 			Bytes:  p.ID,
// 			Status: pgtype.Present,
// 		},
// 		Name: p.Name,
// 	}

// 	if p.ID == uuid.Nil {
// 		row.ID.Status = pgtype.Null
// 	}

// 	return row
// }
