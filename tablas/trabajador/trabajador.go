package tablas

import (
	"main/mysql"
	"main/middleware"
	"main/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func AgregarTrabajador(c *gin.Context) {
	db := mysql.ConexionBD()
	defer db.Close()
	var input models.Trabajador
	middleware.PanicButton(c.ShouldBindJSON(&input))

	id := mysql.Insertar(db, `
		INSERT INTO trabajador (nombre, login, password, admin)
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

func LeerTrabajadores(c *gin.Context) {
	db := mysql.ConexionBD()
	defer db.Close()

	rows := mysql.Leer(db,
		`SELECT id_trabajador, nombre, login, admin FROM trabajador`,
	)
	defer rows.Close()

	var trabajadores []models.Trabajador

	for rows.Next() {
		var t models.Trabajador

		err := rows.Scan(
			&t.Id,
			&t.Nombre,
			&t.Login,
			&t.Admin,
		)
		middleware.PanicButton(err)
		trabajadores = append(trabajadores, t)
	}

	c.JSON(http.StatusOK, gin.H{
		"datos": trabajadores,
	})
}

