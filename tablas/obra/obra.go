package obra

import (
	"main/middleware"
	"main/models"
	"main/mysql"
	"main/tablas/artista"
	"main/tablas/genero"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
)

func CreateObraType(artistaType *graphql.Object, generoType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Obra",
		Fields: graphql.Fields{
			"id_obra": &graphql.Field{Type: graphql.Int},
			"nombre":  &graphql.Field{Type: graphql.String},
			"precio_obra":  &graphql.Field{Type: graphql.Float},
			"fecha_creacion": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					if obra, ok := p.Source.(models.Obra); ok {
						return obra.Fecha_creacion.Format("02-01-2006"), nil
					}
					return nil, nil
				},
			},
			"estatus":     &graphql.Field{Type: graphql.String},
			"foto":        &graphql.Field{Type: graphql.String},
			"material":    &graphql.Field{Type: graphql.String},
			"peso":        &graphql.Field{Type: graphql.Float},
			"dimensiones": &graphql.Field{Type: graphql.String},
			// Relaciones FK usando los tipos ya existentes
			"artista": &graphql.Field{
				Type: artistaType,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					if obra, ok := p.Source.(models.Obra); ok {
						return artista.GetArtistaByID(int(obra.Id_artista))
					}
					return nil, nil
				},
			},
			"genero": &graphql.Field{
				Type: generoType,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					if obra, ok := p.Source.(models.Obra); ok {
						return genero.GetGeneroByID(int(obra.Id_genero))
					}
					return nil, nil
				},
			},
			"id_artista": &graphql.Field{Type: graphql.Int},
			"id_genero":  &graphql.Field{Type: graphql.Int},
		},
	})
}

func ObraExiste(id int) bool {
	var existe bool
	err := mysql.GetBD().QueryRow(
		"SELECT EXISTS(SELECT 1 FROM obra WHERE id_obra = ?)",
		id,
	).Scan(&existe)
	middleware.PanicButton(err)
	return existe
}

func GetObrasField(obraType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(obraType),
		Description: "Lista de obras",
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
			return GetObras(limit, offset)
		},
	}
}

func CreateObraField(obraType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        obraType,
		Description: "Crear una nueva obra",
		Args: graphql.FieldConfigArgument{
			"nombre":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"id_artista":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			"id_genero":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			"precio_obra":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
			"fecha_creacion": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"estatus":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"foto":           &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"material":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"peso":           &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
			"dimensiones":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			nombre := p.Args["nombre"].(string)
			id_artista := p.Args["id_artista"].(int)
			id_genero := p.Args["id_genero"].(int)
			precio_obra := p.Args["precio_obra"].(float64)
			fecha_creacion_str := p.Args["fecha_creacion"].(string)
			estatus := p.Args["estatus"].(string)
			foto := p.Args["foto"].(string)
			material := p.Args["material"].(string)
			peso := p.Args["peso"].(float64)
			dimensiones := p.Args["dimensiones"].(string)

			// validacion de entrada
			errores := []string{}
			if !artista.ArtistaExiste(id_artista) {
				errores = append(errores, "Error: El artista con ID "+strconv.Itoa(id_artista)+" no existe.")
			}

			if !genero.GeneroExiste(id_genero) {
				errores = append(errores, "Error: El género con ID "+strconv.Itoa(id_genero)+" no existe.")
			}

			var i int
			if len(errores) != 0 {
				for i = 0; i < len(errores); i++ {
					println(errores[i])
				}
				return nil, nil
			}

			fecha_creacion, err := time.Parse("2006-01-02", fecha_creacion_str)
			middleware.PanicButton(err)

			id := mysql.Insertar(
				mysql.GetBD(),
				`INSERT INTO obra 
					(nombre, id_artista, id_genero, precio_obra, fecha_creacion, estatus, foto, material, peso, dimensiones)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				nombre, id_artista, id_genero, precio_obra, fecha_creacion, estatus, foto, material, peso, dimensiones,
			)

			return models.Obra{
				Id_obra:        id,
				Nombre:         nombre,
				Id_artista:     int64(id_artista),
				Id_genero:      int64(id_genero),
				Precio_obra:         precio_obra,
				Fecha_creacion: fecha_creacion,
				Estatus:        estatus,
				Foto:           foto,
				Material:       material,
				Peso:           peso,
				Dimensiones:    dimensiones,
			}, nil
		},
	}
}


func GetObras(limit int, offset int) ([]models.Obra, error) {
	var obras []models.Obra
	rows, err := mysql.GetBD().Query(
		`SELECT id_obra, nombre, id_artista, id_genero, precio_obra, fecha_creacion, estatus, foto, material, peso, dimensiones 
		 FROM obra LIMIT ` + strconv.Itoa(limit) + ` OFFSET ` + strconv.Itoa(offset),
	)
	middleware.PanicButton(err)
	defer rows.Close()

	for rows.Next() {
		var o models.Obra
		if err := rows.Scan(
			&o.Id_obra, &o.Nombre, &o.Id_artista, &o.Id_genero, &o.Precio_obra, &o.Fecha_creacion,
			&o.Estatus, &o.Foto, &o.Material, &o.Peso, &o.Dimensiones,
		); err != nil {
			return nil, err
		}
		obras = append(obras, o)
	}
	return obras, nil
}

func UpdateObraField(obraType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        obraType,
		Description: "Actualizar una obra por ID",
		Args: graphql.FieldConfigArgument{
			"id_obra":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			"nombre":         &graphql.ArgumentConfig{Type: graphql.String},
			"id_artista":     &graphql.ArgumentConfig{Type: graphql.Int},
			"id_genero":      &graphql.ArgumentConfig{Type: graphql.Int},
			"precio_obra":         &graphql.ArgumentConfig{Type: graphql.Float},
			"fecha_creacion": &graphql.ArgumentConfig{Type: graphql.String},
			"estatus":        &graphql.ArgumentConfig{Type: graphql.String},
			"foto":           &graphql.ArgumentConfig{Type: graphql.String},
			"material":       &graphql.ArgumentConfig{Type: graphql.String},
			"peso":           &graphql.ArgumentConfig{Type: graphql.Float},
			"dimensiones":    &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id_obra := p.Args["id_obra"].(int)

			// Obtener obra actual
			var o models.Obra
			err := mysql.GetBD().QueryRow(
				`SELECT id_obra, nombre, id_artista, id_genero, precio_obra, fecha_creacion, estatus, foto, material, peso, dimensiones
				 FROM obra WHERE id_obra = ?`,
				id_obra,
			).Scan(&o.Id_obra, &o.Nombre, &o.Id_artista, &o.Id_genero, &o.Precio_obra, &o.Fecha_creacion,
				&o.Estatus, &o.Foto, &o.Material, &o.Peso, &o.Dimensiones)
			middleware.PanicButton(err)

			// Actualizar solo si vienen valores
			if nombre, ok := p.Args["nombre"].(string); ok && nombre != "" {
				o.Nombre = nombre
			}
			if id_artista, ok := p.Args["id_artista"].(int); ok {
				o.Id_artista = int64(id_artista)
			}
			if id_genero, ok := p.Args["id_genero"].(int); ok {
				o.Id_genero = int64(id_genero)
			}
			if precio_obra, ok := p.Args["precio_obra"].(float64); ok {
				o.Precio_obra = precio_obra
			}
			if fechaStr, ok := p.Args["fecha_creacion"].(string); ok && fechaStr != "" {
				fechaParseada, err := time.Parse("2006-01-02", fechaStr)
				middleware.PanicButton(err)
				o.Fecha_creacion = fechaParseada
			}
			if estatus, ok := p.Args["estatus"].(string); ok && estatus != "" {
				o.Estatus = estatus
			}
			if foto, ok := p.Args["foto"].(string); ok && foto != "" {
				o.Foto = foto
			}
			if material, ok := p.Args["material"].(string); ok && material != "" {
				o.Material = material
			}
			if peso, ok := p.Args["peso"].(float64); ok {
				o.Peso = peso
			}
			if dimensiones, ok := p.Args["dimensiones"].(string); ok && dimensiones != "" {
				o.Dimensiones = dimensiones
			}

			// Guardar cambios
			_, err = mysql.GetBD().Exec(
				`UPDATE obra 
				SET 
					nombre=?, 
					id_artista=?, 
					id_genero=?, 
					precio_obra=?, 
					fecha_creacion=?, 
					estatus=?, 
					foto=?, 
					material=?, 
					peso=?, 
					dimensiones=? 
				WHERE id_obra=?`,

				o.Nombre,
				o.Id_artista,
				o.Id_genero,
				o.Precio_obra,
				o.Fecha_creacion,
				o.Estatus,
				o.Foto,
				o.Material,
				o.Peso,
				o.Dimensiones,
				id_obra,
			)
			middleware.PanicButton(err)

			return o, nil
		},
	}
}

func DeleteObraField(obraType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        obraType,
		Description: "Eliminar una obra por ID",
		Args: graphql.FieldConfigArgument{
			"id_obra": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id_obra := p.Args["id_obra"].(int)

			var o models.Obra
			err := mysql.GetBD().QueryRow(
				`SELECT id_obra, nombre, id_artista, id_genero, precio_obra, fecha_creacion, estatus, foto, material, peso, dimensiones
				 FROM obra WHERE id_obra=?`,
				id_obra,
			).Scan(&o.Id_obra, &o.Nombre, &o.Id_artista, &o.Id_genero, &o.Precio_obra, &o.Fecha_creacion,
				&o.Estatus, &o.Foto, &o.Material, &o.Peso, &o.Dimensiones)
			middleware.PanicButton(err)

			_, err = mysql.GetBD().Exec(
				"DELETE FROM obra WHERE id_obra=?",
				id_obra,
			)
			middleware.PanicButton(err)

			return o, nil
		},
	}
}

func GetObraByID(id int) (models.Obra, error) {
	var obra models.Obra
	err := mysql.GetBD().QueryRow(
		"SELECT id_obra, nombre, id_artista, id_genero, precio_obra, fecha_creacion, estatus, foto, material, peso, dimensiones FROM obra WHERE id_obra = ?",
		id,
	).Scan(&obra.Id_obra, &obra.Nombre, &obra.Id_artista, &obra.Id_genero,
		&obra.Precio_obra, &obra.Fecha_creacion, &obra.Estatus,
		&obra.Foto, &obra.Material, &obra.Peso, &obra.Dimensiones)

	if err != nil {
		return models.Obra{}, err
	}
	return obra, nil
}

func GetObraByIDField(obraType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        obraType,
		Description: "Obtener una obra por ID",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id := p.Args["id"].(int)
			return GetObraByID(id) // llama a tu función existente
		},
	}
}