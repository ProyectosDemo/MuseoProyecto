package grapql

import (
	"log"

	"main/middleware"
	"main/mysql"
	"main/tablas/artista"
	"main/tablas/cliente"
	"main/tablas/genero"
	"main/tablas/obra"
	"main/tablas/trabajador"
	"main/tablas/orden"

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
	artistaType := artista.CreateArtistaType()
	generoType := genero.CreateGeneroType()
	obraType := obra.CreateObraType(artistaType, generoType)
	ordenType := orden.CreateOrdenType(clienteType, obraType, trabajadorType)

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"trabajador": trabajador.GetTrabajadorField(trabajadorType),
			"cliente":    cliente.GetClienteField(clienteType),
			"artista":    artista.GetArtistaField(artistaType),
			"genero":     genero.GetGenerosField(generoType),
			"obra":       obra.GetObrasField(obraType),
			"orden":      orden.GetOrdenesField(ordenType),
			"loginCliente": cliente.LoginClienteField(clienteType),
			"loginTrabajador": trabajador.LoginTrabajadorField(trabajadorType),
			"obraById": obra.GetObraByIDField(obraType),
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"crearCliente":      cliente.CreateClienteField(clienteType),
			"eliminarCliente":   cliente.DeleteClienteField(clienteType),
			"actualizarCliente": cliente.UpdateClienteField(clienteType),

			"crearTrabajador":      trabajador.CreateTrabajadorField(trabajadorType),
			"eliminarTrabajador":   trabajador.DeleteTrabajadorField(trabajadorType),
			"actualizarTrabajador": trabajador.UpdateTrabajadorField(trabajadorType),

			"crearArtista":      artista.CreateArtistaField(artistaType),
			"eliminarArtista":   artista.DeleteArtistaField(artistaType),
			"actualizarArtista": artista.UpdateArtistaField(artistaType),

			"crearGenero":      genero.CreateGeneroField(generoType),
			"eliminarGenero":   genero.DeleteGeneroField(generoType),
			"actualizarGenero": genero.UpdateGeneroField(generoType),

			"crearObra":      obra.CreateObraField(obraType),
			"eliminarObra":   obra.DeleteObraField(obraType),
			"actualizarObra": obra.UpdateObraField(obraType),

			"crearOrden":      orden.CreateOrdenField(ordenType),
			"eliminarOrden":   orden.DeleteOrdenField(ordenType),
			"actualizarOrden": orden.UpdateOrdenField(ordenType),
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	middleware.PanicButton(err)

	// ya esto es la parte del servidor
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	router.Any("/graphql", gin.WrapH(h))
	port := "8080"
	log.Println("Servidor GraphQL en http://localhost:" + port + "/graphql")

	err = router.Run(":" + port)
	middleware.PanicButton(err)
}
