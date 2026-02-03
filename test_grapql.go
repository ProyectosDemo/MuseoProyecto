package main

import (
	_"github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"log"
	"net/http"
	"main/mysql"
	"main/models"
)

var trabajadorType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Trabajador",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"nombre": &graphql.Field{
				Type: graphql.String,
			},
			"login": &graphql.Field{
				Type: graphql.String,
			},
			"password": &graphql.Field{
				Type: graphql.String,
			},
			"admin": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		/* Create new product item
		http://localhost:8080/product?query=mutation+_{create(name:"Inca Kola",info:"Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",price:1.99){id,name,info,price}}
		*/
		"create": &graphql.Field{
			Type:        trabajadorType,
			Description: "Crea un nuevo trabajador",
			Args: graphql.FieldConfigArgument{
				"nombre": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"login": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
				"admin": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (any, error) {
				// trabajador de la archivo con los modelos
				trabajador := models.Trabajador{
					Nombre:  params.Args["nombre"].(string),
					Login:  params.Args["login"].(string),
					Password:  params.Args["password"].(string),
					Admin:  params.Args["admin"].(bool),
				}
				return trabajador, nil
			},
		},

	},
})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		//Query:    queryType,
		Mutation: mutationType,
	},
)

func coco() {

	
	base_datos := mysql.ConexionBD()
	defer base_datos.Close()

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
		GraphiQL: true, // habilita GraphQL en el navegador
	})

	

	http.Handle("/graphql", h)
	port := "8080"
	log.Println("Servidor GraphQL corriendo en http://localhost:" + port + "/graphql")
	log.Println(http.ListenAndServe(":"+port, nil))

	// en efecto, no se lo que hago
}

/*

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

*/

