package grapql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"main/models"
	"net/http"
	"strconv"
)

func createTrabajadorType() *graphql.Object {
	return graphql.NewObject(
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
		})
}

func queryTrabajadorType(trabajadorType *graphql.Object) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"trabajador": &graphql.Field{
					Type:        graphql.NewList(trabajadorType),
					Description: "Retorna la lista de trabajadores",
					Args: graphql.FieldConfigArgument{
						"limit": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"offset": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (any, error) {
						limit, _ := p.Args["limit"].(int)
						if limit <= 0 || limit > 20 {
							limit = 10
						}
						offset, _ := p.Args["offset"].(int)
						if offset < 0 {
							offset = 0
						}
						return getTrabajadores(limit, offset)
					},
				},
			},
		},
	)
}

func getTrabajadores(limit int, offset int) ([]models.Trabajador, error) {
	var trabajadores []models.Trabajador
	registros, err := base_datos.Query("SELECT id, nombre, login, password, admin FROM trabajador limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer registros.Close()

	for registros.Next() {
		var aux models.Trabajador
		if err := registros.Scan(&aux.Id, &aux.Nombre, &aux.Login, &aux.Password, &aux.Admin); err != nil {
			return nil, err
		}
		trabajadores = append(trabajadores, aux)
	}
	return trabajadores, nil
}

var base_datos *sql.DB

func conectarBD() {
	var err error
	base_datos, err = sql.Open(
		"mysql",
		"root:16944577aA@tcp(localhost:3306)/museo_proyecto",
	)
	if err != nil {
		log.Fatal(err)
	}
}

func Coco() {
	// conectar a la base de datos
	conectarBD()
	// crear el esquema GraphQL
	trabajadorType := createTrabajadorType()

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryTrabajadorType(trabajadorType),
		})
	if err != nil {
		log.Fatalf("error al crear el esquema: %v", err)
	}
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
