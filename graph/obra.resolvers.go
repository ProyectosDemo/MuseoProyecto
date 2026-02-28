package graph

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/graph/model"
)

// Obras lista de los clientes
func (r *queryResolver) Obras(ctx context.Context, limit *int32, offset *int32) ([]*model.Obra, error) {
    rows, err := r.DB.QueryContext(ctx,
        `SELECT id_obra, nombre, id_artista, id_genero, precio_obra, fecha_creacion, estatus, foto 
         FROM obra LIMIT ? OFFSET ?`,
        limit, offset)
    if err != nil {
        log.Printf("Obras DB error: %v", err)
        return nil, err
    }
    defer rows.Close()

    var obras []*model.Obra
    for rows.Next() {
        var obra model.Obra
        // IDs como strings
        var idObra, idArtista, idGenero string
        err := rows.Scan(&idObra, &obra.Nombre, &idArtista, &idGenero, &obra.Precio, &obra.FechaCreacion, &obra.Status, &obra.Foto)
        if err != nil {
            log.Printf("Obra scan error: %v", err)
            return nil, err
        }

        obra.ID = idObra
        obra.IDArtista = idArtista
        obra.IDGenero = idGenero

        obras = append(obras, &obra)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Obra rows error: %v", err)
        return nil, err
    }

    return obras, nil
}

// FindObra is the resolver for the findobra field.
func (r *queryResolver) FindObra(ctx context.Context, id string) ([]*model.Obra, error) {
	log.Printf("FindCliente llamado con id: %s", id)

	var obra model.Obra
	query := `SELECT id_obra, nombre, id_artista, id_genero, precio_obra, fecha_creacion, estatus, foto FROM obra WHERE id_obra = ?`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&obra.ID, &obra.Nombre, &obra.IDArtista, &obra.IDGenero, &obra.Precio, &obra.FechaCreacion, &obra.Status, &obra.Foto,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("obra no encontrado")
		}
		log.Printf("obraCliente DB error: %v", err)
		return nil, err
	}
	return []*model.Obra{&obra}, nil
}

// UpdateObra is the resolver for the updateObra field.
func (r *mutationResolver) UpdateObra(ctx context.Context, input model.UpdateObra) (*model.Obra, error) {
	log.Printf("UpdateObra llamado con input: %+v", input)

	if input.ID == "" {
		return nil, fmt.Errorf("el ID de la obra es obligatorio para la actualización")
	}

	query := `UPDATE obra SET nombre = COALESCE(?, nombre), id_artista = COALESCE(?, id_artista), id_genero = COALESCE(?, id_genero), precio_obra = COALESCE(?, precio_obra), fecha_creacion = COALESCE(?, fecha_creacion), estatus = COALESCE(?, estatus), foto = COALESCE(?, foto), material = COALESCE(?, material), peso = COALESCE(?, peso), dimensiones = COALESCE(?, dimensiones) WHERE id_obra = ?`

	_, err := r.DB.ExecContext(ctx, query, input.Nombre, input.IDArtista, input.IDGenero, input.Precio, input.FechaCreacion, input.Status, input.Foto, input.Material, input.Peso, input.Dimensiones, input.ID)
	if err != nil {
		log.Printf("UpdateObra DB error: %v", err)
		return nil, err
	}

	var obra model.Obra
	selectQuery := `SELECT id_obra, nombre, id_artista, id_genero, precio_obra, fecha_creacion, estatus, foto, material, peso, dimensiones
	FROM obra WHERE id_obra = ?`

	err = r.DB.QueryRowContext(ctx, selectQuery, input.ID).Scan(
		&obra.ID, &obra.Nombre, &obra.IDArtista, &obra.IDGenero, &obra.Precio, &obra.FechaCreacion, &obra.Status, &obra.Foto, &obra.Material, &obra.Peso, &obra.Dimensiones,
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

// CreateObra is the resolver for the createobra field.
func (r *mutationResolver) CreateObra(ctx context.Context, input model.NewObra) (*model.Obra, error) {
	log.Printf("CreateObra called with input: %+v", input)
	if input.Nombre == "" || input.Foto == "" || input.IDArtista == "" || input.IDGenero == "" {
		return nil, fmt.Errorf("hay campos obligatorios que no pueden estar vacios")
	}
	_, err := r.DB.ExecContext(ctx, "INSERT INTO obra (nombre) VALUES (?)", input.Nombre)
	if err != nil {
		log.Printf("Createobra DB error: %v", err)
		return nil, err
	}
	obra := &model.Obra{
		Nombre:        input.Nombre,
		IDArtista:     input.IDArtista,
		IDGenero:      input.IDGenero,
		Precio:        input.Precio,
		FechaCreacion: input.FechaCreacion,
		Status:        input.Status,
		Foto:          input.Foto,
		Material:      input.Material,
		Peso:          input.Peso,
		Dimensiones:   input.Dimensiones,
	}
	return obra, nil
}

// KillObra is the resolver for the killobra field.
func (r *mutationResolver) KillObra(ctx context.Context, id string) (bool, error) {
	// log de la llamada para depuración
	log.Printf("KillCliente called with id: %s", id)
	if id == "" {
		return false, fmt.Errorf("el ID de la obra es obligatorio para eliminar")
	}
	// ejecucion de la consulta
	_, err := r.DB.ExecContext(ctx, "DELETE FROM obra WHERE id_obra = ?", id)
	if err != nil {
		log.Printf("KillObra DB error: %v", err)
		return false, err
	}
	return true, nil
}
