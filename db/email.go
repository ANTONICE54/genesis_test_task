package db

import (
	"context"
	"time"
)

type Email struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (store *PostgresDB) AddEmail(ctx context.Context, arg string) (Email, error) {
	query := "INSERT INTO emails (email) VALUES ($1) RETURNING id, email, created_at"
	row := store.QueryRowContext(ctx, query, arg)

	var result Email
	err := row.Scan(
		&result.ID,
		&result.Email,
		&result.CreatedAt,
	)

	return result, err

}

func (store *PostgresDB) ListEmails(ctx context.Context) ([]Email, error) {
	query := "SELECT * FROM emails"
	rows, err := store.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := []Email{}

	for rows.Next() {
		var i Email
		err = rows.Scan(
			&i.ID,
			&i.Email,
			&i.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil

}
