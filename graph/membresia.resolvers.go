package graph

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/graph/model"
)

func (r *queryResolver) Membresias(ctx context.Context, limit *int32, offset *int32) ([]*model.Membresia, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT m.id_membresia, m.id_cliente, c.nombre, m.id_tarjeta, t.numero_tarjeta, m.fecha
		 FROM membresia m
		 JOIN cliente c ON m.id_cliente = c.id_cliente
		 JOIN tarjeta_cliente t ON m.id_tarjeta = t.id_tarjeta
		 LIMIT ? OFFSET ?`,
		limit, offset)
	if err != nil {
		log.Printf("Membresias DB error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var membresias []*model.Membresia
	for rows.Next() {
		var m model.Membresia
		var c model.Cliente
		var t model.TarjetaCliente

		err := rows.Scan(&m.IDMembresia, &m.IDCliente, &c.Nombre, &m.IDTarjeta, &t.NumeroTarjeta, &m.Fecha)
		if err != nil {
			log.Printf("Membresia scan error: %v", err)
			return nil, err
		}

		m.Cliente = &c
		m.Tarjeta = &t
		membresias = append(membresias, &m)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Membresia rows error: %v", err)
		return nil, err
	}

	return membresias, nil
}

func (r *queryResolver) FindMembresia(ctx context.Context, id string) (*model.Membresia, error) {
	var m model.Membresia
	var c model.Cliente
	var t model.TarjetaCliente

	query := `SELECT m.id_membresia, m.id_cliente, c.nombre, m.id_tarjeta, t.numero_tarjeta, m.fecha
			  FROM membresia m
			  JOIN cliente c ON m.id_cliente = c.id_cliente
			  JOIN tarjeta_cliente t ON m.id_tarjeta = t.id_tarjeta
			  WHERE m.id_membresia = ?`

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&m.IDMembresia, &m.IDCliente, &c.Nombre, &m.IDTarjeta, &t.NumeroTarjeta, &m.Fecha,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("membresia no encontrada")
		}
		return nil, err
	}

	m.Cliente = &c
	m.Tarjeta = &t
	return &m, nil
}

func (r *mutationResolver) CreateMembresia(ctx context.Context, input model.NewMembresia) (*model.Membresia, error) {
	if input.IDCliente == "" || input.IDTarjeta == "" || input.Fecha == "" {
		return nil, fmt.Errorf("todos los campos son obligatorios")
	}

	insertQuery := `INSERT INTO membresia (id_cliente, id_tarjeta, fecha) VALUES (?, ?, ?)`
	result, err := r.DB.ExecContext(ctx, insertQuery, input.IDCliente, input.IDTarjeta, input.Fecha)
	if err != nil {
		log.Printf("CreateMembresia DB error: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	m := &model.Membresia{
		IDMembresia: fmt.Sprintf("%d", id),
		IDCliente:   input.IDCliente,
		IDTarjeta:   input.IDTarjeta,
		Fecha:       input.Fecha,
	}

	return m, nil
}

func (r *mutationResolver) UpdateMembresia(ctx context.Context, input model.UpdateMembresia) (*model.Membresia, error) {
	log.Printf("UpdateMembresia llamado con input: %+v", input)

	if input.IDMembresia == "" {
		return nil, fmt.Errorf("el ID de la membresia es obligatorio para la actualizacion")
	}

	query := `UPDATE membresia SET
		id_cliente = COALESCE(?, id_cliente),
		id_tarjeta = COALESCE(?, id_tarjeta),
		fecha = COALESCE(?, fecha)
		WHERE id_membresia = ?`

	_, err := r.DB.ExecContext(ctx, query, input.IDCliente, input.IDTarjeta, input.Fecha, input.IDMembresia)
	if err != nil {
		log.Printf("UpdateMembresia DB error: %v", err)
		return nil, err
	}

	var m model.Membresia
	var c model.Cliente
	var t model.TarjetaCliente

	selectQuery := `SELECT m.id_membresia, m.id_cliente, c.nombre, m.id_tarjeta, t.numero_tarjeta, m.fecha
					FROM membresia m
					JOIN cliente c ON m.id_cliente = c.id_cliente
					JOIN tarjeta_cliente t ON m.id_tarjeta = t.id_tarjeta
					WHERE m.id_membresia = ?`

	err = r.DB.QueryRowContext(ctx, selectQuery, input.IDMembresia).Scan(
		&m.IDMembresia, &m.IDCliente, &c.Nombre, &m.IDTarjeta, &t.NumeroTarjeta, &m.Fecha,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("membresia no encontrada")
		}
		log.Printf("Error al leer membresia actualizada: %v", err)
		return nil, err
	}

	m.Cliente = &c
	m.Tarjeta = &t

	return &m, nil
}

func (r *mutationResolver) KillMembresia(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("el ID de la membresia es obligatorio")
	}

	_, err := r.DB.ExecContext(ctx, "DELETE FROM membresia WHERE id_membresia = ?", id)
	if err != nil {
		log.Printf("KillMembresia DB error: %v", err)
		return false, err
	}
	return true, nil
}