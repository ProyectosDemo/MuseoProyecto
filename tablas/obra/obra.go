package tablas

import (
	"main/mysql"
	"main/middleware"
	"main/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func AgregarObra(c *gin.Context) {
	db := mysql.ConexionBD()
	defer db.Close()
	var input models.Obra
	middleware.PanicButton(c.ShouldBindJSON(&input))

	id := mysql.Insertar(db, `
		INSERT INTO obra (nombre, id_artista, id_genero, precio, fecha_creacion, estatus, foto, material, peso, dimensiones)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		input.Nombre,
		input.Id_artista,
		input.Id_genero,
		input.Precio,
		input.Fecha_creacion,
		input.Estatus,
		input.Foto,
		input.Material,
		input.Peso,
		input.Dimensiones,
	)

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Obra agregada",
		"id":      id,
	})
}

func LeerObras(c *gin.Context) {
	db := mysql.ConexionBD()
	defer db.Close()

	rows := mysql.Leer(db,
		`SELECT id_obra, nombre, id_artista, id_genero, precio, fecha_creacion, estatus, foto, material, peso, dimensiones FROM obra`,
	)
	defer rows.Close()

	var obras []models.Obra

	for rows.Next() {
		var t models.Obra

		err := rows.Scan(
			&t.Id_obra,
			&t.Nombre,
			&t.Id_artista,
			&t.Id_genero,
			&t.Precio,
			&t.Fecha_creacion,
			&t.Estatus,
			&t.Foto,
			&t.Material,
			&t.Peso,
			&t.Dimensiones,
		)
		middleware.PanicButton(err)
		obras = append(obras, t)
	}

	c.JSON(http.StatusOK, gin.H{
		"datos": obras,
	})
}

