package repo

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Repo struct {
	db *sql.DB
}

func (r *Repo) transact(f func(tx sql.Tx)) {
	defer db.Commit()
	f(db.Begin())
}
