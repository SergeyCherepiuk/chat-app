package utils

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func MustPrepareNamed(db *sqlx.DB, field **sqlx.NamedStmt, query string) {
	var err error
	*field, err = db.PrepareNamed(query)
	if err != nil {
		panic(err)
	}
}

func MustPrepareNamedMap(db *sqlx.DB, columns []string, field map[string]*sqlx.NamedStmt, query string) {
	var err error
	for _, column := range columns {
		field[column], err = db.PrepareNamed(fmt.Sprintf(query, column))
		if err != nil {
			panic(err)
		}
	}
}
