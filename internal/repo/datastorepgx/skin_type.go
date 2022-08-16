package datastorepgx

import (
	"api/internal/model"
	"context"

	"github.com/georgysavva/scany/pgxscan"
)

func (r *DataStoreRepo) ListSkinType(ctx context.Context) ([]model.SkinType, error) {
	rows, err := r.pool.Query(ctx, "SELECT id,name FROM "+skinTypeTable+" WHERE deleted_time IS NULL ORDER BY created_time asc")
	if err != nil {
		return nil, err
	}

	var skinTypeRows []SkinType
	if err := pgxscan.ScanAll(&skinTypeRows, rows); err != nil {
		return nil, err
	}

	skinTypes := make([]model.SkinType, 0, len(skinTypeRows))
	for _, row := range skinTypeRows {
		skinTypes = append(skinTypes, *row.toModel())
	}

	return skinTypes, nil
}
