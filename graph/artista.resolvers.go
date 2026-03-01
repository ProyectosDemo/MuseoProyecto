package graph

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/graph/model"
)

func (r *queryResolver) Artistas(ctx context.Context, limit *int32, offset *int32) ([]*model.Artista, error) {

	rows, err := r.DB.QueryContext(ctx, "SELECT id_artista, nombre, fecha_nacimiento, nacionalidad, biografia, foto FROM artista LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		log.Printf("Artistas DB error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var artistas []*model.Artista
	for rows.Next() {
		var artista model.Artista
		err := rows.Scan(&artista.ID, &artista.Nombre, &artista.FechaNacimiento, &artista.Nacionalidad, &artista.Biografia, &artista.Foto)
		if err != nil {
			log.Printf("Artistas scan error: %v", err)
			return nil, err
		}
		artistas = append(artistas, &artista)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Artistas rows error: %v", err)
		return nil, err
	}
	return artistas, nil
}

// FindArtista is the resolver for the findArtista field.
func (r *queryResolver) FindArtista(ctx context.Context, id string) ([]*model.Artista, error) {
	log.Printf("FindArtista called with id: %s", id)

	var artista model.Artista
	query := `SELECT id_artista, nombre, fecha_nacimiento, nacionalidad, biografia, foto FROM artista WHERE id_artista = ?`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&artista.ID, &artista.Nombre, &artista.FechaNacimiento, &artista.Nacionalidad,
		&artista.Biografia, &artista.Foto,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("artista no encontrado")
		}
		log.Printf("FindArtista DB error: %v", err)
		return nil, err
	}
	return []*model.Artista{&artista}, nil
}
//codigo
// UpdateArtista is the resolver for the updateArtista field.
func (r *mutationResolver) UpdateArtista(ctx context.Context, input model.UpdateArtista) (*model.Artista, error) {
	log.Printf("UpdateArtista called with input: %+v", input)

	if input.ID == "" {
		return nil, fmt.Errorf("el ID del artista es obligatorio para la actualización")
	}

	query := `UPDATE artista SET 
		nombre = COALESCE(?, nombre), 
		fecha_nacimiento = COALESCE(?, fecha_nacimiento), 
		nacionalidad = COALESCE(?, nacionalidad), 
		biografia = COALESCE(?, biografia), 
		foto = COALESCE(?, foto) 
	WHERE id_artista = ?`

	_, err := r.DB.ExecContext(ctx, query, input.Nombre, input.FechaNacimiento, input.Nacionalidad, input.Biografia, input.Foto, input.ID)
	if err != nil {
		log.Printf("UpdateArtista DB error: %v", err)
		return nil, err
	}

	var artista model.Artista
	selectQuery := `SELECT id_artista, nombre, fecha_nacimiento, nacionalidad, biografia, foto 
	FROM artista WHERE id_artista = ?`

	err = r.DB.QueryRowContext(ctx, selectQuery, input.ID).Scan(
		&artista.ID, &artista.Nombre, &artista.FechaNacimiento, &artista.Nacionalidad,
		&artista.Biografia, &artista.Foto,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("artista no encontrado")
		}
		log.Printf("Error al leer artista actualizado: %v", err)
		return nil, err
	}

	return &artista, nil
}

// CreateArtista is the resolver for the createArtista field.
func (r *mutationResolver) CreateArtista(ctx context.Context, input model.NewArtista) (*model.Artista, error) {
	// Log de entrada para depuración
	log.Printf("CreateArtista called with input: %+v", input)

	//biografica
	// Validaciones básicas de los campos requeridos
	if input.Nombre == "" || input.FechaNacimiento == "" || input.Nacionalidad == "" || input.Biografia == "" || input.Foto == "" {
		return nil, fmt.Errorf("todos los campos de input son obligatorios")
	}

	// Lógica para crear un artista
	artista := &model.Artista{
		Nombre:          input.Nombre,
		FechaNacimiento: input.FechaNacimiento,
		Nacionalidad:    input.Nacionalidad,
		Biografia:       input.Biografia,
		Foto:            input.Foto,
	}
	_, err := r.DB.ExecContext(ctx, "INSERT INTO artista (nombre, fecha_nacimiento, nacionalidad, biografia, foto) VALUES (?, ?, ?, ?, ?)", artista.Nombre, artista.FechaNacimiento, artista.Nacionalidad, artista.Biografia, artista.Foto)
	if err != nil {
		log.Printf("CreateArtista DB error: %v", err)
		return nil, err
	}
	return artista, nil
}

func (r *mutationResolver) KillArtista(ctx context.Context, id string) (bool, error) {
	// log de la llamada para depuración
	log.Printf("KillArtista called with id: %s", id)
	if id == "" {
		return false, fmt.Errorf("el ID del artista es obligatorio para eliminar")
	}
	// ejecucion de la consulta
	_, err := r.DB.ExecContext(ctx, "DELETE FROM artista WHERE id_artista = ?", id)
	if err != nil {
		log.Printf("KillArtista DB error: %v", err)
		return false, err
	}
	return true, nil
}
