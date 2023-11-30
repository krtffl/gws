package repository

import (
	"database/sql"

	"github.com/krtffl/get-well-soon/internal/domain"
)

type gwsRepo struct {
	db *sql.DB
}

func NewGWSRepo(db *sql.DB) domain.GWSRepo {
	return &gwsRepo{
		db: db,
	}
}

func (rep *gwsRepo) List() (
	[]*domain.GWS, error,
) {
	rows, err := rep.db.Query(
		`SELECT * FROM "Messages" ORDER BY "CreatedAt" DESC`,
	)
	if err != nil {
		return nil, handleErrors(err)
	}

	defer rows.Close()
	var gwss []*domain.GWS
	for rows.Next() {
		gws := &domain.GWS{}
		if err := rows.Scan(
			&gws.Id,
			&gws.From,
			&gws.Message,
			&gws.Memory,
			&gws.CreatedAt,
		); err != nil {
			return nil, handleErrors(err)
		}
		gwss = append(gwss, gws)
	}

	return gwss, nil
}

func (rep *gwsRepo) Create(gws *domain.GWS) (
	*domain.GWS, error,
) {
	err := rep.db.QueryRow(
		`
        INSERT INTO "Messages" ("Id", "From", "Message", "Memory")
        VALUES ($1, $2, $3, $4)
        RETURNING *`,
		gws.Id,
		gws.From,
		gws.Message,
		gws.Memory,
	).Scan(&gws.Id,
		&gws.From,
		&gws.Message,
		&gws.Memory,
		&gws.CreatedAt)
	if err != nil {
		return nil, handleErrors(err)
	}
	return gws, nil
}
