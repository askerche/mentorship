package repo

import "github.com/jackc/pgx/v5"

type Repo struct {
	DB *pgx.Conn
}
