package artista

import (
	"main/middleware"
	"main/models"
	"main/mysql"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
)

//artistaType
// Podemos copiar estas funciones en las demas tablas, cambiando los campos, es casi lo mismo de antes, pero estas funciones trabajan con field ahora

func CreateArtistaType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Artista",
		Fields: graphql.Fields{
			"id_artista": &graphql.Field{Type: graphql.Int},
			"nombre":     &graphql.Field{Type: graphql.String},
			"fecha_nacimiento": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					if artista, ok := p.Source.(models.Artista); ok {
						return artista.Fecha_nacimiento.Format("02-01-2006"), nil
					}
					return nil, nil
				},
			},
			"nacionalidad": &graphql.Field{Type: graphql.String},
			"biografia":    &graphql.Field{Type: graphql.String},
			"foto":         &graphql.Field{Type: graphql.String},
		},
	})
}

func GetArtistaField(artistaType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(artistaType),
		Description: "Lista de artistas",
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
			return GetArtistas(limit, offset)
		},
	}
}

func CreateArtistaField(artistaType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        artistaType,
		Description: "Crear un nuevo artista",
		Args: graphql.FieldConfigArgument{
			"nombre":           &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"fecha_nacimiento": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"nacionalidad":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"biografia":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"foto":             &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			nombre := p.Args["nombre"].(string)
			fechaNacimientoStr := p.Args["fecha_nacimiento"].(string)
			nacionalidad := p.Args["nacionalidad"].(string)
			biografia := p.Args["biografia"].(string)
			foto := p.Args["foto"].(string)

			// Convertir string a time.Time
			fechaNacimiento, err := time.Parse("2006-01-02", fechaNacimientoStr)
			middleware.PanicButton(err)

			id := mysql.Insertar(
				mysql.GetBD(),
				"INSERT INTO artista (nombre, fecha_nacimiento, nacionalidad, biografia, foto) VALUES (?, ?, ?, ?, ?)",
				nombre, fechaNacimiento, nacionalidad, biografia, foto,
			)

			return models.Artista{
				Id_artista:       id,
				Nombre:           nombre,
				Fecha_nacimiento: fechaNacimiento,
				Nacionalidad:     nacionalidad,
				Biografia:        biografia,
				Foto:             foto,
			}, nil
		},
	}
}

func GetArtistas(limit int, offset int) ([]models.Artista, error) {
	var artistas []models.Artista
	registros, err := mysql.GetBD().Query(
		"SELECT id_artista, nombre, fecha_nacimiento, nacionalidad, biografia, foto FROM artista LIMIT " +
			strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset),
	)
	middleware.PanicButton(err)
	defer registros.Close()

	for registros.Next() {
		var aux models.Artista
		if err := registros.Scan(
			&aux.Id_artista, &aux.Nombre, &aux.Fecha_nacimiento, &aux.Nacionalidad,
			&aux.Biografia, &aux.Foto,
		); err != nil {
			return nil, err
		}
		artistas = append(artistas, aux)
	}
	return artistas, nil
}

func DeleteArtistaField(artistaType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        artistaType,
		Description: "Eliminar un artista por ID",
		Args: graphql.FieldConfigArgument{
			"id_artista": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id := p.Args["id_artista"].(int)

			var artista models.Artista
			err := mysql.GetBD().QueryRow(
				"SELECT id_artista, nombre, fecha_nacimiento, nacionalidad, biografia, foto FROM artista WHERE id_artista = ?",
				id,
			).Scan(&artista.Id_artista, &artista.Nombre, &artista.Fecha_nacimiento, &artista.Nacionalidad,
				&artista.Biografia, &artista.Foto)

			middleware.PanicButton(err)

			_, err = mysql.GetBD().Exec(
				"DELETE FROM artista WHERE id_artista = ?",
				id,
			)

			middleware.PanicButton(err)

			return artista, nil
		},
	}
}
func UpdateArtistaField(artistaType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        artistaType,
		Description: "Actualizar un artista por ID",
		Args: graphql.FieldConfigArgument{
			"id_artista": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"nombre": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"fecha_nacimiento": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"nacionalidad": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"biografia": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"foto": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {

			id_artista := p.Args["id_artista"].(int)

			// Obtener artista actual
			var artista models.Artista
			err := mysql.GetBD().QueryRow(
				"SELECT id_artista, nombre, fecha_nacimiento, nacionalidad, biografia, foto FROM artista WHERE id_artista = ?",
				id_artista,
			).Scan(&artista.Id_artista, &artista.Nombre, &artista.Fecha_nacimiento,
				&artista.Nacionalidad, &artista.Biografia, &artista.Foto)

			middleware.PanicButton(err)

			// Actualizar solo si vienen valores
			if nombre, ok := p.Args["nombre"].(string); ok && nombre != "" {
				artista.Nombre = nombre
			}

			if fechaStr, ok := p.Args["fecha_nacimiento"].(string); ok && fechaStr != "" {
				fechaParseada, err := time.Parse("2006-01-02", fechaStr)
				middleware.PanicButton(err)
				artista.Fecha_nacimiento = fechaParseada
			}

			if nacionalidad, ok := p.Args["nacionalidad"].(string); ok && nacionalidad != "" {
				artista.Nacionalidad = nacionalidad
			}

			if biografia, ok := p.Args["biografia"].(string); ok && biografia != "" {
				artista.Biografia = biografia
			}

			if foto, ok := p.Args["foto"].(string); ok && foto != "" {
				artista.Foto = foto
			}

			// Guardar cambios
			_, err = mysql.GetBD().Exec(`
				UPDATE artista
				SET nombre = ?, 
				    fecha_nacimiento = ?, 
				    nacionalidad = ?, 
				    biografia = ?, 
				    foto = ?
				WHERE id_artista = ?
			`,
				artista.Nombre,
				artista.Fecha_nacimiento,
				artista.Nacionalidad,
				artista.Biografia,
				artista.Foto,
				id_artista,
			)

			middleware.PanicButton(err)

			return artista, nil
		},
	}
}

func ArtistaExiste(id int) bool {
	var existe bool
	err := mysql.GetBD().QueryRow(
		"SELECT EXISTS(SELECT 1 FROM artista WHERE id_artista = ?)",
		id,
	).Scan(&existe)
	middleware.PanicButton(err)
	return existe
}

func GetArtistaByID(id int) (models.Artista, error) {
	var artista models.Artista
	err := mysql.GetBD().QueryRow(
		"SELECT id_artista, nombre, fecha_nacimiento, nacionalidad, biografia, foto FROM artista WHERE id_artista = ?",
		id,
	).Scan(&artista.Id_artista, &artista.Nombre, &artista.Fecha_nacimiento,
		&artista.Nacionalidad, &artista.Biografia, &artista.Foto)

	if err != nil {
		return models.Artista{}, err
	}
	return artista, nil
}
