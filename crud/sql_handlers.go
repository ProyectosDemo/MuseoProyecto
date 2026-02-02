package crud

import (
	"database/sql"
	"main/middleware"
)

func Insertar(db *sql.DB, query string, args ...interface{}) int64 {
	result, err := db.Exec(query, args...)
	middleware.PanicButton(err)

	id, err := result.LastInsertId()
	middleware.PanicButton(err)

	return id
}