package tablas

import (
	"main/middleware"
	"main/models"
	"main/mysql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
)

func CreateTrabajadorType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Trabajador",
			Fields: graphql.Fields{
				"id_trabajador": &graphql.Field{
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

func MutationTrabajadorType(trabajadorType *graphql.Object) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"crearTrabajador": &graphql.Field{
					Type:        trabajadorType,
					Description: "Crear un nuevo trabajador",
					Args: graphql.FieldConfigArgument{
						"nombre": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"login": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"password": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"admin": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Boolean),
						},
					},
					Resolve: func(p graphql.ResolveParams) (any, error) {
						nombre := p.Args["nombre"].(string)
						login := p.Args["login"].(string)
						password := p.Args["password"].(string)
						admin := p.Args["admin"].(bool)

						id := mysql.Insertar(
							mysql.GetBD(),
							"INSERT INTO trabajador (nombre, login, password, admin) VALUES (?, ?, ?, ?)",
							nombre,
							login,
							password,
							admin,
						)

						return models.Trabajador{
							Id:       id,
							Nombre:   nombre,
							Login:    login,
							Password: password,
							Admin:    admin,
						}, nil
					},
				},
			},
		},
	)
}

func GetTrabajadores(limit int, offset int) ([]models.Trabajador, error) {
	var trabajadores []models.Trabajador
	registros, err := mysql.GetBD().Query("SELECT id_trabajador, nombre, login, password, admin FROM trabajador limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset))
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
