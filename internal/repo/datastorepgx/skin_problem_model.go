package datastorepgx

import (
	"api/internal/model"

	"github.com/jackc/pgtype"
)

type SkinProblem struct {
	ID   pgtype.UUID `db:"id"`
	Name string      `db:"name"`
}

func (s SkinProblem) toModel() model.SkinProblem {
	return model.SkinProblem{
		ID:   s.ID.Bytes,
		Name: s.Name,
	}
}
