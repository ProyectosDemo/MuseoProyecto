package cliente

import (
	"main/middleware"
	"main/models"
	"main/mysql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
)

// Podemos copiar estas funciones en las demas tablas, cambiando los campos, es casi lo mismo de antes, pero estas funciones trabajan con field ahora

func CreateClienteType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Cliente",
		Fields: graphql.Fields{
			"id_cliente":              &graphql.Field{Type: graphql.Int},
			"nombre":          &graphql.Field{Type: graphql.String},
			"email":           &graphql.Field{Type: graphql.String},
			"telefono":        &graphql.Field{Type: graphql.String},
			"login":           &graphql.Field{Type: graphql.String},
			"password":        &graphql.Field{Type: graphql.String},
			"codigo_seguridad": &graphql.Field{Type: graphql.String},
		},
	})
}


func GetClienteField(clienteType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(clienteType),
		Description: "Lista de clientes",
		Args: graphql.FieldConfigArgument{
			"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
			"offset": &graphql.ArgumentConfig{Type: graphql.Int},
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
			return GetClientes(limit, offset)
		},
	}
}

func CreateClienteField(clienteType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        clienteType,
		Description: "Crear un nuevo cliente",
		Args: graphql.FieldConfigArgument{
			"nombre":           &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"email":            &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"telefono":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"login":            &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"password":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"codigo_seguridad": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			nombre := p.Args["nombre"].(string)
			email := p.Args["email"].(string)
			telefono := p.Args["telefono"].(string)
			login := p.Args["login"].(string)
			password := p.Args["password"].(string)
			codigoSeguridad := p.Args["codigo_seguridad"].(string)

			id := mysql.Insertar(
				mysql.GetBD(),
				"INSERT INTO cliente (nombre, email, telefono, login, password, codigo_seguridad) VALUES (?, ?, ?, ?, ?, ?)",
				nombre, email, telefono, login, password, codigoSeguridad,
			)

			return models.Cliente{
				Id:              id,
				Nombre:          nombre,
				Email:           email,
				Telefono:        telefono,
				Login:           login,
				Password:        password,
				CodigoSeguridad: codigoSeguridad,
			}, nil
		},
	}
}

func GetClientes(limit int, offset int) ([]models.Cliente, error) {
	var clientes []models.Cliente
	registros, err := mysql.GetBD().Query(
		"SELECT id_cliente, nombre, email, telefono, login, password, codigo_seguridad FROM cliente LIMIT " +
			strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset),
	)
	middleware.PanicButton(err)
	defer registros.Close()

	for registros.Next() {
		var aux models.Cliente
		if err := registros.Scan(
			&aux.Id, &aux.Nombre, &aux.Email, &aux.Telefono,
			&aux.Login, &aux.Password, &aux.CodigoSeguridad,
		); err != nil {
			return nil, err
		}
		clientes = append(clientes, aux)
	}
	return clientes, nil
}
