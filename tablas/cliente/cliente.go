package tablas

import (
	"main/mysql"
	"main/middleware"
	"main/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func AgregarCliente(c *gin.Context) {
	db := mysql.ConexionBD()
	defer db.Close()
	var input models.Cliente
	middleware.PanicButton(c.ShouldBindJSON(&input))

	id := mysql.Insertar(db, `
		INSERT INTO cliente (nombre, email, telefono, login, password, codigo_seguridad)
		VALUES (?, ?, ?, ?, ?, ?)`,
		input.Nombre,
		input.Email,
		input.Telefono,
		input.Login,
		input.Password,
		input.CodigoSeguridad,
	)

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Cliente agregado",
		"id":      id,
	})
}

func LeerCliente(c *gin.Context) {
	db := mysql.ConexionBD()
	defer db.Close()

	rows := mysql.Leer(db,
		`SELECT id_cliente, nombre, email, telefono, login, password, codigo_seguridad FROM cliente`,
	)
	defer rows.Close()

	var clientes []models.Cliente
	for rows.Next() {
		var t models.Cliente

		err := rows.Scan(
			&t.Id,
			&t.Nombre,
			&t.Email,
			&t.Telefono,
			&t.Login,
			&t.Password,
			&t.CodigoSeguridad,
		)
		middleware.PanicButton(err)
		clientes = append(clientes, t)
	}

	c.JSON(http.StatusOK, gin.H{
		"datos": clientes,
	})
}

