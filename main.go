package main

import (
	"database/sql"
	"fmt"
	"main/crud"
	"main/middleware"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() *sql.DB {
	db, err := sql.Open(
		"mysql",
		"root:16944577aA@tcp(localhost:3306)/museo_proyecto",
	)
	middleware.PanicButton(err)
	return db
}

func main() {
	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.POST("/trabajador", agregarTrabajador)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	router.Run(":8080")
}

func agregarTrabajador(c *gin.Context) {
	db := conexionBD()
	defer db.Close()
	var input models.TrabajadorInput
	middleware.PanicButton(c.ShouldBindJSON(&input))

	id := crud.Insertar(db, `
		INSERT INTO trabajador (nombre, login, password, administrador)
		VALUES (?, ?, ?, ?)`,
		input.Nombre,
		input.Login,
		input.Password,
		input.Admin,
	)

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Trabajador agregado",
		"id":      id,
	})
}
