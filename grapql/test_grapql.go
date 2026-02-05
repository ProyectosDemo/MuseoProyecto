package grapql

import (
	"log"
	"main/middleware"
	"main/mysql"
	"main/tablas/trabajador"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func Coco() {
	// conectar a la base de datos
	mysql.ConectarBD()
	// crear el esquema GraphQL
	trabajadorType := tablas.CreateTrabajadorType()

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    tablas.QueryTrabajadorType(trabajadorType),
			Mutation: tablas.MutationTrabajadorType(trabajadorType),
		})
	middleware.PanicButton(err)
	// manejador GraphQL
	handler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true, // habilita GraphQL en el navegador
	})
	// iniciar el servidor HTTP
	port := "8080"
	http.Handle("/graphql", handler)
	log.Println("Servidor GraphQL corriendo en http://localhost:" + port + "/graphql")
	log.Println(http.ListenAndServe(":"+port, nil))

	// en efecto, no se lo que hago
}
