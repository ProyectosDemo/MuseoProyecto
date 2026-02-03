package tablas

import (
	"main/mysql"
	"main/middleware"
	"main/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func AgregarArtista(c *gin.Context) {
	db := mysql.ConexionBD()
	defer db.Close()
	var input models.Artista
	middleware.PanicButton(c.ShouldBindJSON(&input))

	id := mysql.Insertar(db, `
		INSERT INTO artista (nombre, fecha_nacimiento, nacionalidad, biografia, foto)
		VALUES (?, ?, ?, ?, ?)`,
		input.Nombre,
		input.Fecha_nacimiento,
		input.Nacionalidad,
		input.Biografia,
		input.Foto,
	)

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Artista agregado",
		"id":      id,
	})
}

func LeerArtista(c *gin.Context) {
	db := mysql.ConexionBD()
	defer db.Close()

	rows := mysql.Leer(db,
		`SELECT id_artista, nombre, fecha_nacimiento, nacionalidad, biografia, foto FROM artista`,
	)
	defer rows.Close()

	var trabajadores []models.Artista

	for rows.Next() {
		var t models.Artista

		err := rows.Scan(
			&t.Id_artista,
			&t.Nombre,
			&t.Fecha_nacimiento,
			&t.Nacionalidad,
			&t.Biografia,
			&t.Foto,
		)
		middleware.PanicButton(err)
		trabajadores = append(trabajadores, t)
	}

	c.JSON(http.StatusOK, gin.H{
		"datos": trabajadores,
	})
}

