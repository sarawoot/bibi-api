package datastorepgx

import (
	"api/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type User struct {
	ID            pgtype.UUID `db:"id"`
	Email         string      `db:"email"`
	PasswordHash  string      `db:"password_hash"`
	Gender        string      `db:"gender"`
	Birthdate     pgtype.Date `db:"birthdate"`
	SkinTypeID    pgtype.UUID `db:"skin_type_id"`
	SkinProblemID pgtype.UUID `db:"skin_problem_id"`
	CreatedTime   pgtype.Time `db:"created_time"`
}

func (u *User) toModel() model.User {
	return model.User{
		ID:            u.ID.Bytes,
		Email:         u.Email,
		PasswordHash:  u.PasswordHash,
		Gender:        u.Gender,
		Birthdate:     u.Birthdate.Time,
		SkinTypeID:    u.SkinTypeID.Bytes,
		SkinProblemID: u.SkinProblemID.Bytes,
	}
}

func userModelToRow(user model.User) User {
	row := User{
		ID: pgtype.UUID{
			Bytes:  user.ID,
			Status: pgtype.Present,
		},
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Gender:       user.Gender,
		Birthdate: pgtype.Date{
			Time:   user.Birthdate,
			Status: pgtype.Present,
		},
		SkinTypeID: pgtype.UUID{
			Bytes:  user.SkinTypeID,
			Status: pgtype.Present,
		},
		SkinProblemID: pgtype.UUID{
			Bytes:  user.SkinProblemID,
			Status: pgtype.Present,
		},
	}

	if user.ID == uuid.Nil {
		row.ID.Status = pgtype.Null
	}

	if user.Birthdate.IsZero() {
		row.Birthdate.Status = pgtype.Null
	}

	if user.SkinTypeID == uuid.Nil {
		row.SkinTypeID.Status = pgtype.Null
	}

	if user.SkinProblemID == uuid.Nil {
		row.SkinProblemID.Status = pgtype.Null
	}

	return row
}
