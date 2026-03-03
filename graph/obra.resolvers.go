package graph

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/graph/model"
)

// Obras lista todas las obras
func (r *queryResolver) Obras(ctx context.Context, limit *int32, offset *int32) ([]*model.Obra, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT o.id_obra, o.nombre, o.precio_obra, o.fecha_creacion, o.status, o.foto, o.material, o.peso, o.dimensiones,
		        a.id_artista, a.nombre,
		        g.id_genero, g.nombre
		 FROM obra o
		 JOIN artista a ON o.id_artista = a.id_artista
		 JOIN genero g ON o.id_genero = g.id_genero
		 LIMIT ? OFFSET ?`,
		limit, offset)
	if err != nil {
		log.Printf("Obras DB error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var obras []*model.Obra
	for rows.Next() {
		var obra model.Obra
		var artista model.Artista
		var genero model.Genero

		err := rows.Scan(&obra.ID, &obra.Nombre, &obra.Precio, &obra.FechaCreacion, &obra.Status, &obra.Foto, &obra.Material, &obra.Peso, &obra.Dimensiones,
			&artista.ID, &artista.Nombre,
			&genero.ID, &genero.Nombre)
		if err != nil {
			log.Printf("Obra scan error: %v", err)
			return nil, err
		}

		obra.Artista = &artista
		obra.Genero = &genero

		obras = append(obras, &obra)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Obra rows error: %v", err)
		return nil, err
	}

	return obras, nil
}

// FindObra busca una obra por ID
func (r *queryResolver) FindObra(ctx context.Context, id string) ([]*model.Obra, error) {
	var obra model.Obra
	var artista model.Artista
	var genero model.Genero

	query := `
	SELECT o.id_obra, o.nombre, o.precio_obra, o.fecha_creacion, o.status, o.foto, o.material, o.peso, o.dimensiones,
	       a.id_artista, a.nombre, a.fecha_nacimiento, a.nacionalidad, a.biografia, a.foto,
	       g.id_genero, g.nombre
	FROM obra o
	JOIN artista a ON o.id_artista = a.id_artista
	JOIN genero g ON o.id_genero = g.id_genero
	WHERE o.id_obra = ?
	`

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&obra.ID, &obra.Nombre, &obra.Precio, &obra.FechaCreacion, &obra.Status, &obra.Foto, &obra.Material, &obra.Peso, &obra.Dimensiones,
		&artista.ID, &artista.Nombre, &artista.FechaNacimiento, &artista.Nacionalidad, &artista.Biografia, &artista.Foto,
		&genero.ID, &genero.Nombre,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("obra no encontrada")
		}
		return nil, err
	}

	obra.Artista = &artista
	obra.Genero = &genero

	return []*model.Obra{&obra}, nil
}

var validStatuses = map[string]bool{
	"DISPONIBLE": true,
	"RESERVADA":  true,
	"VENDIDA":    true,
}

func (r *mutationResolver) UpdateObra(ctx context.Context, input model.UpdateObra) (*model.Obra, error) {
	log.Printf("UpdateObra llamado con input: %+v", input)

	if input.ID == "" {
		return nil, fmt.Errorf("el ID de la obra es obligatorio para la actualización")
	}

	// Validar status si viene
	var statusValue interface{}
	if input.Status != nil {
		statusStr := string(*input.Status)
		if !validStatuses[statusStr] {
			return nil, fmt.Errorf("status inválido, debe ser DISPONIBLE, RESERVADA o VENDIDA")
		}
		statusValue = statusStr
	} else {
		statusValue = nil
	}

	query := `UPDATE obra SET 
		nombre = COALESCE(?, nombre),
		id_artista = COALESCE(?, id_artista),
		id_genero = COALESCE(?, id_genero),
		precio_obra = COALESCE(?, precio_obra),
		fecha_creacion = COALESCE(?, fecha_creacion),
		status = COALESCE(?, status),
		foto = COALESCE(?, foto),
		material = COALESCE(?, material),
		peso = COALESCE(?, peso),
		dimensiones = COALESCE(?, dimensiones)
		WHERE id_obra = ?`

	_, err := r.DB.ExecContext(ctx, query,
		input.Nombre, input.IDArtista, input.IDGenero, input.Precio, input.FechaCreacion,
		statusValue, input.Foto, input.Material, input.Peso, input.Dimensiones, input.ID,
	)
	if err != nil {
		log.Printf("UpdateObra DB error: %v", err)
		return nil, err
	}

	// Leer la obra actualizada
	var obra model.Obra
	selectQuery := `SELECT id_obra, nombre, id_artista, id_genero, precio_obra, fecha_creacion, status, foto, material, peso, dimensiones FROM obra WHERE id_obra = ?`
	err = r.DB.QueryRowContext(ctx, selectQuery, input.ID).Scan(
		&obra.ID, &obra.Nombre, &obra.IDArtista, &obra.IDGenero, &obra.Precio, &obra.FechaCreacion,
		&obra.Status, &obra.Foto, &obra.Material, &obra.Peso, &obra.Dimensiones,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("obra no encontrada")
		}
		log.Printf("Error al leer obra actualizada: %v", err)
		return nil, err
	}

	return &obra, nil
}

// -------------------------------
// CreateObra inserta una nueva obra
// -------------------------------
func (r *mutationResolver) CreateObra(ctx context.Context, input model.NewObra) (*model.Obra, error) {
	log.Printf("CreateObra called with input: %+v", input)

	if input.Nombre == "" || input.Foto == "" || input.IDArtista == "" || input.IDGenero == "" {
		return nil, fmt.Errorf("hay campos obligatorios que no pueden estar vacíos")
	}

	// Validar status
	statusStr := string(input.Status)
	if !validStatuses[statusStr] {
		return nil, fmt.Errorf("status inválido, debe ser DISPONIBLE, RESERVADA o VENDIDA")
	}

	insertQuery := `INSERT INTO obra 
		(nombre, id_artista, id_genero, precio_obra, fecha_creacion, status, foto, material, peso, dimensiones) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.DB.ExecContext(ctx, insertQuery,
		input.Nombre, input.IDArtista, input.IDGenero, input.Precio, input.FechaCreacion,
		statusStr, input.Foto, input.Material, input.Peso, input.Dimensiones,
	)
	if err != nil {
		log.Printf("CreateObra DB error: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	obra := &model.Obra{
		ID:            fmt.Sprintf("%d", id),
		Nombre:        input.Nombre,
		IDArtista:     input.IDArtista,
		IDGenero:      input.IDGenero,
		Precio:        input.Precio,
		FechaCreacion: input.FechaCreacion,
		Status:        input.Status, // sigue usando StatusObra de gqlgen
		Foto:          input.Foto,
		Material:      input.Material,
		Peso:          input.Peso,
		Dimensiones:   input.Dimensiones,
	}

	return obra, nil
}

// KillObra elimina una obra por ID
func (r *mutationResolver) KillObra(ctx context.Context, id string) (bool, error) {
	log.Printf("KillObra called with id: %s", id)
	if id == "" {
		return false, fmt.Errorf("el ID de la obra es obligatorio para eliminar")
	}
	_, err := r.DB.ExecContext(ctx, "DELETE FROM obra WHERE id_obra = ?", id)
	if err != nil {
		log.Printf("KillObra DB error: %v", err)
		return false, err
	}
	return true, nil
}