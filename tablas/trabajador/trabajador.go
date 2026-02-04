package tablas

import (
	"database/sql"
	"strconv"
	"github.com/graphql-go/graphql"
	_ "github.com/go-sql-driver/mysql"
	"main/models"
	"main/middleware"
)


func CreateTrabajadorType() *graphql.Object {
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

func QueryTrabajadorType(trabajadorType *graphql.Object) *graphql.Object {
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
						return GetTrabajadores(limit, offset)
					},
				},
			},
		},
	)
}

func GetTrabajadores(limit int, offset int) ([]models.Trabajador, error) {
	var trabajadores []models.Trabajador
	registros, err := base_datos.Query("SELECT id, nombre, login, password, admin FROM trabajador limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset))
	middleware.PanicButton(err)
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
