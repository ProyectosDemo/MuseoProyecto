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


func DeleteClienteField(clienteType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        clienteType,
		Description: "Eliminar un cliente por ID",
		Args: graphql.FieldConfigArgument{
			"id_cliente": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id := p.Args["id_cliente"].(int)

			var cliente models.Cliente
			err := mysql.GetBD().QueryRow(
				"SELECT id_cliente, nombre, email, telefono, login, password, codigo_seguridad FROM cliente WHERE id_cliente = ?",
				id,
			).Scan(&cliente.Id, &cliente.Nombre, &cliente.Email, &cliente.Telefono,
				&cliente.Login, &cliente.Password, &cliente.CodigoSeguridad)

			if err != nil {
				return nil, err
			}

			_, err = mysql.GetBD().Exec(
				"DELETE FROM cliente WHERE id_cliente = ?",
				id,
			)

			if err != nil {
				return nil, err
			}

			return cliente, nil
		},
	}
}

func UpdateClienteField(clienteType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        clienteType,
		Description: "Actualizar un cliente por ID",
		Args: graphql.FieldConfigArgument{
			"id_cliente": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"nombre": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"email": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"telefono": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"login": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"password": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"codigo_seguridad": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id := p.Args["id_cliente"].(int)

			// Valores opcionales
			nombre, _ := p.Args["nombre"].(string)
			email, _ := p.Args["email"].(string)
			telefono, _ := p.Args["telefono"].(string)
			login, _ := p.Args["login"].(string)
			password, _ := p.Args["password"].(string)
			codigoSeguridad, _ := p.Args["codigo_seguridad"].(string)

			// mantiene original si es vacio o invalido
			_, err := mysql.GetBD().Exec(`
				UPDATE cliente
				SET nombre = COALESCE(NULLIF(?, ''), nombre),
				    email = COALESCE(NULLIF(?, ''), email),
				    telefono = COALESCE(NULLIF(?, ''), telefono),
					login = COALESCE(NULLIF(?, ''), login),
				    password = COALESCE(NULLIF(?, ''), password),
				    codigo_seguridad = COALESCE(NULLIF(?, ''), codigo_seguridad)
				WHERE id_cliente = ?
			`,
				nombre,
				email,
				telefono,
				login,
				password,
				codigoSeguridad,
				id,
			)
			if err != nil {
				return nil, err
			}

			// Traer el cliente actualizado
			var cliente models.Cliente
			err = mysql.GetBD().QueryRow(
				"SELECT id_cliente, nombre, email, telefono, login, password, codigo_seguridad FROM cliente WHERE id_cliente = ?",
				id,
			).Scan(&cliente.Id, &cliente.Nombre, &cliente.Email, &cliente.Telefono,
				&cliente.Login, &cliente.Password, &cliente.CodigoSeguridad)

			if err != nil {
				return nil, err
			}

			return cliente, nil
		},
	}
}