package datastorepgx

import (
	"api/internal/model"
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

func (r *DataStoreRepo) CreateUser(ctx context.Context, user *model.User) error {
	row := userModelToRow(*user)

	userID := user.ID
	if userID == uuid.Nil {
		userID = uuid.New()
	}
	user.ID = userID

	id := pgtype.UUID{Bytes: userID, Status: pgtype.Present}
	_, err := r.pool.Exec(
		ctx,
		"INSERT INTO "+userTable+"(id,email,password_hash,gender,birthdate,skin_type_id,skin_problem_id) VALUES ($1,$2,$3,$4,$5,$6,$7)",
		id, row.Email, row.PasswordHash, row.Gender, row.Birthdate, row.SkinTypeID, row.SkinProblemID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DataStoreRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	row, err := r.pool.Query(ctx, "SELECT * FROM "+userTable+" WHERE email=$1 LIMIT 1", email)
	if err != nil {
		return nil, err
	}

	var user User
	if err := pgxscan.ScanOne(&user, row); err != nil {
		return nil, err
	}

	return user.toModel(), nil
}
