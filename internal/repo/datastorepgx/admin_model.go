package datastorepgx

import (
	"api/internal/model"

	"github.com/jackc/pgtype"
)

type Admin struct {
	ID           pgtype.UUID `db:"id"`
	Username     string      `db:"username"`
	PasswordHash string      `db:"password_hash"`
	CreatedTime  pgtype.Time `db:"created_time"`
}

func (a *Admin) toModel() *model.Admin {
	return &model.Admin{
		ID:           a.ID.Bytes,
		Username:     a.Username,
		PasswordHash: a.PasswordHash,
	}
}
