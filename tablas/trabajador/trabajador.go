package trabajador

import (
	"main/middleware"
	"main/models"
	"main/mysql"
	"strconv"

	"github.com/graphql-go/graphql"
	_ "github.com/go-sql-driver/mysql"
)

func CreateTrabajadorType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Trabajador",
		Fields: graphql.Fields{
			"id_trabajador": &graphql.Field{Type: graphql.Int},
			"nombre":        &graphql.Field{Type: graphql.String},
			"login":         &graphql.Field{Type: graphql.String},
			"password":      &graphql.Field{Type: graphql.String},
			"admin":         &graphql.Field{Type: graphql.Boolean},
		},
	})
}

func GetTrabajadorField(trabajadorType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(trabajadorType),
		Description: "Lista de trabajadores",
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
			return GetTrabajadores(limit, offset)
		},
	}
}

func CreateTrabajadorField(trabajadorType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        trabajadorType,
		Description: "Crear un nuevo trabajador",
		Args: graphql.FieldConfigArgument{
			"nombre":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"login":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"admin":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Boolean)},
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
	}
}

func GetTrabajadores(limit int, offset int) ([]models.Trabajador, error) {
	var trabajadores []models.Trabajador
	registros, err := mysql.GetBD().Query("SELECT id_trabajador, nombre, login, password, admin FROM trabajador LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset))
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

func DeleteTrabajadorField(trabajadorType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        trabajadorType,
		Description: "Eliminar un trabajador por ID",
		Args: graphql.FieldConfigArgument{
			"id_trabajador": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id := p.Args["id_trabajador"].(int)

			var trabajador models.Trabajador
			err := mysql.GetBD().QueryRow(
				"SELECT id_trabajador, nombre, login, password, admin FROM trabajador WHERE id_trabajador = ?",
				id,
			).Scan(&trabajador.Id, &trabajador.Nombre, &trabajador.Login, &trabajador.Password, &trabajador.Admin)

			middleware.PanicButton(err)

			_, err = mysql.GetBD().Exec(
				"DELETE FROM trabajador WHERE id_trabajador = ?",
				id,
			)

			middleware.PanicButton(err)

			return trabajador, nil
		},
	}
}

func UpdateTrabajadorField(trabajadorType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        trabajadorType,
		Description: "Actualizar un trabajador por ID",
		Args: graphql.FieldConfigArgument{
			"id_trabajador": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"nombre": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"login": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"password": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"admin": &graphql.ArgumentConfig{
				Type: graphql.Boolean,
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id := p.Args["id_trabajador"].(int)

			// Valores opcionales
			nombre, _ := p.Args["nombre"].(string)
			login, _ := p.Args["login"].(string)
			password, _ := p.Args["password"].(string)
			admin, adminOk := p.Args["admin"].(bool)

			// mantenemos original si es vacio o invalido
			_, err := mysql.GetBD().Exec(`
				UPDATE trabajador
				SET nombre = COALESCE(NULLIF(?, ''), nombre),
				    login = COALESCE(NULLIF(?, ''), login),
				    password = COALESCE(NULLIF(?, ''), password),
				    admin = COALESCE(?, admin)
				WHERE id_trabajador = ?
			`,
				nombre,
				login,
				password,
				func() any { if adminOk { return admin } else { return nil } }(),
				id,
			)
			middleware.PanicButton(err)

			// Traer el trabajador actualizado
			var trabajador models.Trabajador
			err = mysql.GetBD().QueryRow(
				"SELECT id_trabajador, nombre, login, password, admin FROM trabajador WHERE id_trabajador = ?",
				id,
			).Scan(&trabajador.Id, &trabajador.Nombre, &trabajador.Login, &trabajador.Password, &trabajador.Admin)

			middleware.PanicButton(err)

			return trabajador, nil
		},
	}
}
