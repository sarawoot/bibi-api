package datastorepgx

import (
	"api/internal/model"

	"github.com/jackc/pgtype"
)

type SkinType struct {
	ID   pgtype.UUID `db:"id"`
	Name string      `db:"name"`
}

func (s *SkinType) toModel() model.SkinType {
	return model.SkinType{
		ID:   s.ID.Bytes,
		Name: s.Name,
	}
}
