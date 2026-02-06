package grapql

import (
	"log"

	"main/middleware"
	"main/mysql"
	"main/tablas/trabajador"
	"main/tablas/cliente"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func Coco() {
	mysql.ConectarBD()

	// Crear tipos para cada tabla
	trabajadorType := trabajador.CreateTrabajadorType()
	clienteType := cliente.CreateClienteType()

	// Añadir cuantos queries quieras
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"trabajador": trabajador.GetTrabajadorField(trabajadorType),
			"cliente":    cliente.GetClienteField(clienteType),
		},
	})

	// Añadir cuantos mutagenos quieras
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"crearTrabajador": trabajador.CreateTrabajadorField(trabajadorType),
			"crearCliente":    cliente.CreateClienteField(clienteType),
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
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

	err = router.Run(":" + port)
	middleware.PanicButton(err)
}
