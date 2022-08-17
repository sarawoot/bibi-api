package datastorepgx

import (
	"api/internal/model"

	"github.com/jackc/pgtype"
)

type Country struct {
	ID          pgtype.UUID `db:"id"`
	Name        string      `db:"name"`
	CreatedTime pgtype.Time `db:"created_time"`
	DeletedTime pgtype.Time `db:"deleted_time"`
}

func (s *Country) toModel() model.Country {
	return model.Country{
		ID:   s.ID.Bytes,
		Name: s.Name,
	}
}

// func CountryModelToRow(p model.Country) Country {
// 	row := Country{
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
