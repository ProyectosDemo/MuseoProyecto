package mysql

import (
	"database/sql"
	"main/middleware"
)

var base_datos *sql.DB

func ConectarBD() {
	var err error
	base_datos, err = sql.Open(
		"mysql",
		"root:Delfin#100@tcp(localhost:3306)/museo_proyecto",
	)
	middleware.PanicButton(err)
}

func GetBD() *sql.DB {
	return base_datos
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
