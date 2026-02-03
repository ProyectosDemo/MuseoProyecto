package tablas

import (
	"main/mysql"
	"main/middleware"
	"main/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func AgregarOrden(c *gin.Context) {
	db := mysql.ConexionBD()
	defer db.Close()
	var input models.Orden
	middleware.PanicButton(c.ShouldBindJSON(&input))

	id := mysql.Insertar(db, `
		INSERT INTO orden (id_cliente, id_obra, id_trabajador, precio, iva, ganancia_museo, total, fecha_orden, estatus)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		input.Id_cliente,
		input.Id_obra,
		input.Id_trabajador,
		input.Precio,
		input.Iva,
		input.Ganancia_museo,
		input.Total,
		input.Fecha_orden,
		input.Estatus,
	)

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Orden agregada",
		"id":      id,
	})
}

func LeerOrdenes(c *gin.Context) {
	db := mysql.ConexionBD()
	defer db.Close()

	rows := mysql.Leer(db,
		`SELECT id_orden, id_cliente, id_obra, id_trabajador, precio, iva, ganancia_museo, total, fecha_orden, estatus FROM orden`,
	)
	defer rows.Close()

	var ordenes []models.Orden

	for rows.Next() {
		var t models.Orden

		err := rows.Scan(
			&t.Id_orden,
			&t.Id_cliente,
			&t.Id_obra,
			&t.Id_trabajador,
			&t.Precio,
			&t.Iva,
			&t.Ganancia_museo,
			&t.Total,
			&t.Fecha_orden,
			&t.Estatus,
		)
		middleware.PanicButton(err)
		ordenes = append(ordenes, t)
	}

	c.JSON(http.StatusOK, gin.H{
		"datos": ordenes,
	})
}

