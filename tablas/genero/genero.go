package tablas

import (
	"main/mysql"
	"main/middleware"
	"main/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func AgregarGenero(c *gin.Context) {
	db := mysql.ConexionBD()
	defer db.Close()
	var input models.Genero
	middleware.PanicButton(c.ShouldBindJSON(&input))

	id := mysql.Insertar(db, `
		INSERT INTO genero (nombre)
		VALUES (?)`,
		input.Nombre,
	)

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Genero agregado",
		"id":      id,
	})
}

func LeerGeneros(c *gin.Context) {
	db := mysql.ConexionBD()
	defer db.Close()

	rows := mysql.Leer(db,
		`SELECT id_genero, nombre FROM genero`,
	)
	defer rows.Close()

	var generos []models.Genero

	for rows.Next() {
		var t models.Genero

		err := rows.Scan(
			&t.Id_genero,
			&t.Nombre,
		)
		middleware.PanicButton(err)
		generos = append(generos, t)
	}

	c.JSON(http.StatusOK, gin.H{
		"datos": generos,
	})
}

