package repo

import "github.com/jackc/pgx/v5"

type rows struct {
	rows pgx.Rows
}

// Err implements modb.Rows.
func (r *rows) Err() error {
	return r.rows.Err()
}

func (r *rows) Close() error {
	r.rows.Close()
	return nil
}

func (r *rows) Next() bool {
	return r.rows.Next()
}

func (r *rows) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}
