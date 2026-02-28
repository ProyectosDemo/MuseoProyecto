package graph

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/graph/model"
)

// FindGenero is the resolver for the findgenero field.
func (r *queryResolver) FindGenero(ctx context.Context, id string) ([]*model.Genero, error) {
	log.Printf("FindGenero llamado con id: %s", id)

	var genero model.Genero
	query := `SELECT id_genero, nombre FROM genero WHERE id_genero = ?`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&genero.ID, &genero.Nombre,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("genero no encontrado")
		}
		log.Printf("Genero DB error: %v", err)
		return nil, err
	}
	return []*model.Genero{&genero}, nil
}

// UpdateGenero is the resolver for the updategenero field.
func (r *mutationResolver) UpdateGenero(ctx context.Context, input model.UpdateGenero) (*model.Genero, error) {
	log.Printf("UpdateGenero llamado con input: %+v", input)

	if input.ID == "" {
		return nil, fmt.Errorf("el ID del genero es obligatorio para la actualización")
	}

	query := `UPDATE genero SET nombre = COALESCE(?, nombre) WHERE id_genero = ?`

	_, err := r.DB.ExecContext(ctx, query, input.Nombre, input.ID)
	if err != nil {
		log.Printf("UpdateGenero DB error: %v", err)
		return nil, err
	}

	var genero model.Genero
	selectQuery := `SELECT id_genero, nombre
	FROM genero WHERE id_genero = ?`

	err = r.DB.QueryRowContext(ctx, selectQuery, input.ID).Scan(
		&genero.ID, &genero.Nombre,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("genero no encontrado")
		}
		log.Printf("Error al leer genero actualizado: %v", err)
		return nil, err
	}

	return &genero, nil
}

// CreateGenero is the resolver for the createGenero field.
func (r *mutationResolver) CreateGenero(ctx context.Context, nombre string) (*model.Genero, error) {
	log.Printf("CreateGenero called with nombre: %s", nombre)
	if nombre == "" {
		return nil, fmt.Errorf("el nombre es obligatorio")
	}
	_, err := r.DB.ExecContext(ctx, "INSERT INTO genero (nombre) VALUES (?)", nombre)
	if err != nil {
		log.Printf("CreateGenero DB error: %v", err)
		return nil, err
	}
	genero := &model.Genero{
		Nombre: nombre,
	}
	return genero, nil
}

// KillGenero is the resolver for the killGenero field.
func (r *mutationResolver) KillGenero(ctx context.Context, id string) (bool, error) {
	// log de la llamada para depuración
	log.Printf("KillCliente called with id: %s", id)
	if id == "" {
		return false, fmt.Errorf("el ID del genero es obligatorio para eliminar")
	}
	// ejecucion de la consulta
	_, err := r.DB.ExecContext(ctx, "DELETE FROM genero WHERE id_genero = ?", id)
	if err != nil {
		log.Printf("KillGenero DB error: %v", err)
		return false, err
	}
	return true, nil
}
