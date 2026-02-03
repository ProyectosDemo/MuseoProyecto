package main

import (
	"fmt"
	"main/middleware"
	"github.com/gin-gonic/gin"
	"main/tablas/trabajador"
)

func main() {
	// ----------------------  main ----------------------
	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	// Rutas para la tabla Trabajador
	router.GET("/trabajadores", tablas.LeerTrabajadores)
	router.POST("/trabajador", tablas.AgregarTrabajador)
	// Rutas para otras tablas pueden ser añadidas aqui

	fmt.Println("Servidor corriendo en http://localhost:8080")
	router.Run(":8080")
}
