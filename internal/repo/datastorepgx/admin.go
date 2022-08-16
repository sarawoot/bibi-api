package datastorepgx

import (
	"api/internal/model"
	"context"

	"github.com/georgysavva/scany/pgxscan"
)

func (r *DataStoreRepo) GetAdminByUsername(ctx context.Context, username string) (*model.Admin, error) {
	row, err := r.pool.Query(ctx, "SELECT * FROM "+adminTable+" WHERE username=$1 LIMIT 1", username)
	if err != nil {
		return nil, err
	}

	var admin Admin
	if err := pgxscan.ScanOne(&admin, row); err != nil {
		return nil, err
	}

	return admin.toModel(), nil
}
