package mysql

import (
	"database/sql"
	"main/middleware"
)

func ConexionBD() *sql.DB {
	db, err := sql.Open(
		"mysql",
		"root:16944577aA@tcp(localhost:3306)/museo_proyecto",
	)
	middleware.PanicButton(err)
	return db
}


func Insertar(db *sql.DB, query string, args ...any) int64 {
	result, err := db.Exec(query, args...)
	middleware.PanicButton(err)

	id, err := result.LastInsertId()
	middleware.PanicButton(err)

	return id
}


func Leer(db *sql.DB, query string, args ...any) *sql.Rows {
	rows, err := db.Query(query, args...)
	middleware.PanicButton(err)

	return rows
}

//TO DO: Update and Delete functions