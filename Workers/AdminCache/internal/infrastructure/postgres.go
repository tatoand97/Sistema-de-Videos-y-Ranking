package infrastructure

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func MustPostgres(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil { panic(err) }
	if err = db.Ping(); err != nil { panic(err) }
	return db
}
