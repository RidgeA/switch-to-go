package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type (
	Material struct {
		Url string
	}

	MaterialRepository struct {
		conn *pgx.Conn
	}
)

func NewMaterialRepository(conn *pgx.Conn) *MaterialRepository {
	return &MaterialRepository{conn}
}

func (mr MaterialRepository) GetMaterials(ctx context.Context) ([]Material, error) {
	rows, err := mr.conn.Query(ctx, "SELECT url FROM materials")
	if err != nil {
		return nil, fmt.Errorf("cannot load list of materials: %w", err)
	}
	defer rows.Close()

	result := make([]Material, 0)

	for rows.Next() {
		m := Material{}
		err := rows.Scan(&m.Url)
		if err != nil {
			return nil, fmt.Errorf("cannot read db record: %w", err)
		}
		result = append(result, m)
	}

	return result, nil
}
