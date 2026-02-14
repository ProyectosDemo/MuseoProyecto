package genero

import (
	"main/middleware"
	"main/models"
	"main/mysql"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
)

// Podemos copiar estas funciones en las demas tablas, cambiando los campos, es casi lo mismo de antes, pero estas funciones trabajan con field ahora

func CreateGeneroType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Genero",
		Fields: graphql.Fields{
			"id_genero":  &graphql.Field{Type: graphql.Int},
			"nombre":     &graphql.Field{Type: graphql.String},
		},
	})
}


func GetGenerosField(generoType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(generoType),
		Description: "Lista de generos",
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
			return GetGeneros(limit, offset)
		},
	}
}

func CreateGeneroField(generoType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        generoType,
		Description: "Crear un nuevo genero",
		Args: graphql.FieldConfigArgument{
			"nombre": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			nombre := p.Args["nombre"].(string)
			id_genero := mysql.Insertar(
				mysql.GetBD(),
				"INSERT INTO genero (nombre) VALUES (?)",
				nombre,
			)

			return models.Genero{
				Id_genero: id_genero,
				Nombre:    nombre,
			}, nil
		},
	}
}

func GetGeneros(limit int, offset int) ([]models.Genero, error) {
	var generos []models.Genero
	registros, err := mysql.GetBD().Query(
		"SELECT id_genero, nombre FROM genero LIMIT " +
			strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset),
	)
	middleware.PanicButton(err)
	defer registros.Close()

	for registros.Next() {
		var aux models.Genero
		if err := registros.Scan(
			&aux.Id_genero, &aux.Nombre,
		); err != nil {
			return nil, err
		}
		generos = append(generos, aux)
	}
	return generos, nil
}


func DeleteGeneroField(generoType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        generoType,
		Description: "Eliminar un genero por ID",
		Args: graphql.FieldConfigArgument{
			"id_genero": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id_genero := p.Args["id_genero"].(int)

			var genero models.Genero
			err := mysql.GetBD().QueryRow(
				"SELECT id_genero, nombre FROM genero WHERE id_genero = ?",
				id_genero,
			).Scan(&genero.Id_genero, &genero.Nombre)

			if err != nil {
				return nil, err
			}

			_, err = mysql.GetBD().Exec(
				"DELETE FROM genero WHERE id_genero = ?",
				id_genero,
			)

			if err != nil {
				return nil, err
			}

			return genero, nil
		},
	}
}
func UpdateGeneroField(generoType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        generoType,
		Description: "Actualizar un genero por ID",
		Args: graphql.FieldConfigArgument{
			"id_genero": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"nombre": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {

			id_genero := p.Args["id_genero"].(int)

			// Obtener genero actual
			var genero models.Genero
			err := mysql.GetBD().QueryRow(
				"SELECT id_genero, nombre FROM genero WHERE id_genero = ?",
				id_genero,
			).Scan(&genero.Id_genero, &genero.Nombre)

			if err != nil {
				return nil, err
			}

			// Actualizar solo si vienen valores
			if nombre, ok := p.Args["nombre"].(string); ok && nombre != "" {
				genero.Nombre = nombre
			}

			// Guardar cambios
			_, err = mysql.GetBD().Exec(`
				UPDATE genero
				SET nombre = ? 
				WHERE id_genero = ?
			`,
				genero.Nombre,
				id_genero,
			)

			if err != nil {
				return nil, err
			}

			return genero, nil
		},
	}
}
