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
		//solo descomenta la tuya y comenta la mia
		//"root:Delfin#100@tcp(localhost:3306)/museo_proyecto?parseTime=true",
		"root:16944577aA@tcp(localhost:3306)/museo_proyecto?parseTime=true&loc=Local",
	)
	middleware.PanicButton(err)
}

func GetBD() *sql.DB {
	return base_datos
}
