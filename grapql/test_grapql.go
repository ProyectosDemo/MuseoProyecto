package grapql

import (
	"log"

	"main/middleware"
	"main/mysql"
	"main/tablas/trabajador"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func Coco() {

	mysql.ConectarBD()

	trabajadorType := trabajador.CreateTrabajadorType()

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    trabajador.QueryTrabajadorType(trabajadorType),
			Mutation: trabajador.MutationTrabajadorType(trabajadorType),
		})

	middleware.PanicButton(err)

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Any("/graphql", gin.WrapH(h))

	port := "8080"

	log.Println("Servidor GraphQL en http://localhost:" + port + "/graphql")

	// Run server
	err = router.Run(":" + port)
	middleware.PanicButton(err)
}
