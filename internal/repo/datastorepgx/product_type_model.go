package datastorepgx

import (
	"api/internal/model"

	"github.com/jackc/pgtype"
)

type ProductType struct {
	ID   pgtype.UUID `db:"id"`
	Name string      `db:"name"`
}

func (s ProductType) toModel() model.ProductType {
	return model.ProductType{
		ID:   s.ID.Bytes,
		Name: s.Name,
	}
}

// func productTypeModelToRow(p model.ProductType) ProductType {
// 	row := ProductType{
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
