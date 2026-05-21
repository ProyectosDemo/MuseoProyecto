package graph

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/graph/model"
)

func (r *queryResolver) Ordenes(ctx context.Context, limit *int32, offset *int32) ([]*model.Orden, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id_orden, id_obra, id_cliente, id_trabajador, fecha, status, direccion
		 FROM orden
		 LIMIT ? OFFSET ?`,
		limit, offset)
	if err != nil {
		log.Printf("Ordenes DB error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var ordenes []*model.Orden
	for rows.Next() {
		var orden model.Orden
		var idTrabajador, direccion sql.NullString

		err := rows.Scan(&orden.ID, &orden.IDObra, &orden.IDCliente, &idTrabajador, &orden.Fecha, &orden.Status, &direccion)
		if err != nil {
			log.Printf("Orden scan error: %v", err)
			return nil, err
		}

		if idTrabajador.Valid {
			orden.IDTrabajador = idTrabajador.String
		} else {
			orden.IDTrabajador = ""
		}

		// Solo asignar direccion si la orden está concretada
		if orden.Status == "CONCRETADA" && direccion.Valid {
			orden.Direccion = &direccion.String
		}

		ordenes = append(ordenes, &orden)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Orden rows error: %v", err)
		return nil, err
	}

	return ordenes, nil
}

func (r *queryResolver) FindOrden(ctx context.Context, id string) (*model.Orden, error) {
	var orden model.Orden
	var idTrabajador, direccion sql.NullString

	query := `SELECT id_orden, id_obra, id_cliente, id_trabajador, fecha, status, direccion
	          FROM orden WHERE id_orden = ?`

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&orden.ID, &orden.IDObra, &orden.IDCliente, &idTrabajador, &orden.Fecha, &orden.Status, &direccion,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("orden no encontrada")
		}
		return nil, err
	}

	if idTrabajador.Valid {
		orden.IDTrabajador = idTrabajador.String
	} else {
		orden.IDTrabajador = ""
	}

	// Solo asignar direccion si la orden está concretada
	if orden.Status == "CONCRETADA" && direccion.Valid {
		orden.Direccion = &direccion.String
	}

	return &orden, nil
}

var validOrdenStatus = map[string]bool{
	"PENDIENTE":   true,
	"CONCRETADA":  true,
	"CANCELADA":   true,
}

func (r *mutationResolver) CreateOrden(ctx context.Context, input model.NewOrden) (*model.Orden, error) {
	if input.IDObra == "" || input.IDCliente == "" || input.Fecha == "" {
		return nil, fmt.Errorf("faltan campos obligatorios")
	}

	var idTrabajador interface{}
	if input.IDTrabajador != nil {
		idTrabajador = *input.IDTrabajador
	} else {
		idTrabajador = nil
	}

	statusStr := string(input.Status)
	if !validOrdenStatus[statusStr] {
		return nil, fmt.Errorf("status inválido, debe ser PENDIENTE, CONCRETADA o CANCELADA")
	}

	insertQuery := `INSERT INTO orden (id_obra, id_cliente, id_trabajador, fecha, status, direccion)
	                VALUES (?, ?, ?, ?, ?, ?)`

	var direccion interface{}
	if statusStr == "CONCRETADA" && input.Direccion != nil {
		direccion = *input.Direccion
	} else {
		direccion = nil
	}

	result, err := r.DB.ExecContext(ctx, insertQuery,
		input.IDObra, input.IDCliente, idTrabajador, input.Fecha, statusStr, direccion)
	if err != nil {
		log.Printf("CreateOrden DB error: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var idTrabajadorStr string
	if input.IDTrabajador != nil {
		idTrabajadorStr = *input.IDTrabajador
	}

	var direccionStr *string
	if statusStr == "CONCRETADA" && input.Direccion != nil {
		direccionStr = input.Direccion
	}

	return &model.Orden{
		ID:           fmt.Sprintf("%d", id),
		IDObra:       input.IDObra,
		IDCliente:    input.IDCliente,
		IDTrabajador: idTrabajadorStr,
		Fecha:        input.Fecha,
		Status:       input.Status,
		Direccion:    direccionStr,
	}, nil
}

func (r *mutationResolver) UpdateOrden(ctx context.Context, input model.UpdateOrden) (*model.Orden, error) {
	log.Printf("UpdateOrden llamado con input: %+v", input)
	if input.ID == "" {
		return nil, fmt.Errorf("el ID de la orden es obligatorio para la actualizacion")
	}

	var statusValue interface{}
	if input.Status != nil {
		statusStr := string(*input.Status)
		if !validOrdenStatus[statusStr] {
			return nil, fmt.Errorf("status inválido, debe ser PENDIENTE, CONCRETADA o CANCELADA")
		}
		statusValue = statusStr
	}

	query := `UPDATE orden SET 
        id_obra = COALESCE(?, id_obra),
        id_cliente = COALESCE(?, id_cliente),
        id_trabajador = COALESCE(?, id_trabajador),
        fecha = COALESCE(?, fecha),
        status = COALESCE(?, status),
        direccion = COALESCE(?, direccion)
        WHERE id_orden = ?`

	var direccionValue interface{}
	if input.Direccion != nil {
		direccionValue = *input.Direccion
	} else {
		direccionValue = nil
	}

	_, err := r.DB.ExecContext(ctx, query,
		input.IDObra, input.IDCliente, input.IDTrabajador, input.Fecha, statusValue, direccionValue, input.ID,
	)
	if err != nil {
		log.Printf("UpdateOrden DB error: %v", err)
		return nil, err
	}

	var orden model.Orden
	selectQuery := `SELECT id_orden, id_obra, id_cliente, id_trabajador, fecha, status, direccion
                    FROM orden WHERE id_orden = ?`
	err = r.DB.QueryRowContext(ctx, selectQuery, input.ID).Scan(
		&orden.ID, &orden.IDObra, &orden.IDCliente, &orden.IDTrabajador, &orden.Fecha, &orden.Status, &orden.Direccion,
	)
	if err != nil {
		return nil, err
	}

	// Solo asignar direccion si la orden está concretada
	if orden.Status != "CONCRETADA" {
		orden.Direccion = nil
	}

	obraStatus := "RESERVADA"
	switch orden.Status {
	case "CONCRETADA":
		obraStatus = "VENDIDA"
	case "CANCELADA":
		obraStatus = "DISPONIBLE"
	case "PENDIENTE":
		obraStatus = "RESERVADA"
	}

	_, err = r.DB.ExecContext(ctx, "UPDATE obra SET status = ? WHERE id_obra = ?", obraStatus, orden.IDObra)
	if err != nil {
		log.Printf("Error actualizando status de la obra: %v", err)
	}

	return &orden, nil
}

func (r *mutationResolver) KillOrden(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("el ID de la orden es obligatorio")
	}

	_, err := r.DB.ExecContext(ctx, "DELETE FROM orden WHERE id_orden = ?", id)
	if err != nil {
		log.Printf("KillOrden DB error: %v", err)
		return false, err
	}

	return true, nil
}