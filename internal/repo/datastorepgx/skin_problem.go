package datastorepgx

import (
	"api/internal/model"
	"context"

	"github.com/georgysavva/scany/pgxscan"
)

func (r *DataStoreRepo) ListSkinProblem(ctx context.Context) ([]model.SkinProblem, error) {
	rows, err := r.pool.Query(ctx, "SELECT id,name FROM "+skinProblemTable+" WHERE deleted_time IS NULL ORDER BY created_time asc")
	if err != nil {
		return nil, err
	}

	var SkinProblemRows []SkinProblem
	if err := pgxscan.ScanAll(&SkinProblemRows, rows); err != nil {
		return nil, err
	}

	SkinProblems := make([]model.SkinProblem, 0, len(SkinProblemRows))
	for _, row := range SkinProblemRows {
		SkinProblems = append(SkinProblems, *row.toModel())
	}

	return SkinProblems, nil
}
