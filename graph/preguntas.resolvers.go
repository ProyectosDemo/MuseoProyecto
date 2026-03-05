package graph

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/graph/model"
)

// Preguntas lista todas las preguntas
func (r *queryResolver) Preguntas(ctx context.Context, limit *int32, offset *int32) ([]*model.Preguntas, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT p.id, p.id_cliente, p.pregunta, p.respuesta,
		        c.id, c.nombre, c.email, c.telefono, c.login, c.password, c.codigo_seguridad
		 FROM preguntas p
		 JOIN cliente c ON p.id_cliente = c.id_cliente
		 LIMIT ? OFFSET ?`,
		limit, offset)
	if err != nil {
		log.Printf("Preguntas DB error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var preguntasList []*model.Preguntas
	for rows.Next() {
		var preg model.Preguntas
		var cliente model.Cliente

		err := rows.Scan(&preg.ID, &preg.IDCliente, &preg.Pregunta, &preg.Respuesta,
			&cliente.ID, &cliente.Nombre, &cliente.Email, &cliente.Telefono, &cliente.Login, &cliente.Password, &cliente.CodigoSeguridad)
		if err != nil {
			log.Printf("Preguntas scan error: %v", err)
			return nil, err
		}

		preg.Cliente = &cliente
		preguntasList = append(preguntasList, &preg)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Preguntas rows error: %v", err)
		return nil, err
	}

	return preguntasList, nil
}

// FindPreguntas busca una pregunta por ID
func (r *queryResolver) FindPreguntas(ctx context.Context, id string) (*model.Preguntas, error) {
	var preg model.Preguntas
	var cliente model.Cliente

	query := `
	SELECT p.id, p.id_cliente, p.pregunta, p.respuesta,
	       c.id, c.nombre, c.email, c.telefono, c.login, c.password, c.codigo_seguridad
	FROM preguntas p
	JOIN cliente c ON p.id_cliente = c.id_cliente
	WHERE p.id = ?`

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&preg.ID, &preg.IDCliente, &preg.Pregunta, &preg.Respuesta,
		&cliente.ID, &cliente.Nombre, &cliente.Email, &cliente.Telefono, &cliente.Login, &cliente.Password, &cliente.CodigoSeguridad,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pregunta no encontrada")
		}
		return nil, err
	}

	preg.Cliente = &cliente
	return &preg, nil
}

// FindPreguntasByCliente busca todas las preguntas de un cliente
func (r *queryResolver) FindPreguntasByCliente(ctx context.Context, id_cliente string) ([]*model.Preguntas, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id_pregunta, id_cliente, pregunta, respuesta
		 FROM preguntas
		 WHERE id_cliente = ?`, id_cliente)
	if err != nil {
		log.Printf("FindPreguntasByCliente DB error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var preguntasList []*model.Preguntas
	for rows.Next() {
		var preg model.Preguntas
		preg.Cliente = &model.Cliente{ID: id_cliente}

		err := rows.Scan(&preg.ID, &preg.IDCliente, &preg.Pregunta, &preg.Respuesta)
		if err != nil {
			log.Printf("FindPreguntasByCliente scan error: %v", err)
			return nil, err
		}

		preguntasList = append(preguntasList, &preg)
	}

	return preguntasList, nil
}

// CreatePreguntas inserta una nueva pregunta
func (r *mutationResolver) CreatePreguntas(ctx context.Context, input model.NewPreguntas) (*model.Preguntas, error) {
	if input.IDCliente == "" || input.Pregunta == "" || input.Respuesta == "" {
		return nil, fmt.Errorf("todos los campos son obligatorios")
	}

	insertQuery := `INSERT INTO preguntas (id_cliente, pregunta, respuesta) VALUES (?, ?, ?)`
	result, err := r.DB.ExecContext(ctx, insertQuery, input.IDCliente, input.Pregunta, input.Respuesta)
	if err != nil {
		log.Printf("CreatePreguntas DB error: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	preg := &model.Preguntas{
		ID:        fmt.Sprintf("%d", id),
		IDCliente: input.IDCliente,
		Pregunta:  input.Pregunta,
		Respuesta: input.Respuesta,
	}

	return preg, nil
}

// UpdatePreguntas actualiza una pregunta existente
func (r *mutationResolver) UpdatePreguntas(ctx context.Context, input model.UpdatePreguntas) (*model.Preguntas, error) {
	if input.ID == "" {
		return nil, fmt.Errorf("el ID es obligatorio para actualizar")
	}

	updateQuery := `UPDATE preguntas SET 
		pregunta = COALESCE(?, pregunta),
		respuesta = COALESCE(?, respuesta)
		WHERE id = ?`

	_, err := r.DB.ExecContext(ctx, updateQuery, input.Pregunta, input.Respuesta, input.ID)
	if err != nil {
		log.Printf("UpdatePreguntas DB error: %v", err)
		return nil, err
	}

	// Leer la pregunta actualizada
	var preg model.Preguntas
	err = r.DB.QueryRowContext(ctx, `SELECT id, id_cliente, pregunta, respuesta FROM preguntas WHERE id = ?`, input.ID).
		Scan(&preg.ID, &preg.IDCliente, &preg.Pregunta, &preg.Respuesta)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pregunta no encontrada")
		}
		log.Printf("Error al leer pregunta actualizada: %v", err)
		return nil, err
	}

	return &preg, nil
}

// KillPreguntas elimina una pregunta por ID
func (r *mutationResolver) KillPreguntas(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("el ID es obligatorio para eliminar")
	}
	_, err := r.DB.ExecContext(ctx, "DELETE FROM preguntas WHERE id = ?", id)
	if err != nil {
		log.Printf("KillPreguntas DB error: %v", err)
		return false, err
	}
	return true, nil
}