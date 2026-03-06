package graph

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/graph/model"
)
// TarjetasCliente lista de todas las tarjetas con el cliente incluido
func (r *queryResolver) TarjetasCliente(ctx context.Context, limit *int32, offset *int32) ([]*model.TarjetaCliente, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT t.id_tarjeta, t.id_cliente, t.numero_tarjeta, t.tipo,
		       c.nombre, c.email, c.telefono, c.login, c.password, c.codigo_seguridad
		FROM tarjeta_cliente t
		LEFT JOIN cliente c ON t.id_cliente = c.id_cliente
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		log.Printf("TarjetasCliente DB error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tarjetas []*model.TarjetaCliente
	for rows.Next() {
		var tarjeta model.TarjetaCliente
		var cliente model.Cliente
		err := rows.Scan(
			&tarjeta.IDTarjeta,
			&tarjeta.IDCliente,
			&tarjeta.NumeroTarjeta,
			&tarjeta.Tipo,
			&cliente.Nombre,
			&cliente.Email,
			&cliente.Telefono,
			&cliente.Login,
			&cliente.Password,
			&cliente.CodigoSeguridad,
		)
		if err != nil {
			log.Printf("TarjetasCliente scan error: %v", err)
			return nil, err
		}
		tarjeta.Cliente = &cliente
		tarjetas = append(tarjetas, &tarjeta)
	}
	if err = rows.Err(); err != nil {
		log.Printf("TarjetasCliente rows error: %v", err)
		return nil, err
	}
	return tarjetas, nil
}

// FindTarjetaCliente busca una tarjeta por su ID con cliente
func (r *queryResolver) FindTarjetaCliente(ctx context.Context, id string) (*model.TarjetaCliente, error) {
	log.Printf("FindTarjetaCliente called with id: %s", id)

	var tarjeta model.TarjetaCliente
	var cliente model.Cliente
	query := `
		SELECT t.id_tarjeta, t.id_cliente, t.numero_tarjeta, t.tipo,
		       c.nombre, c.email, c.telefono, c.login, c.password, c.codigo_seguridad
		FROM tarjeta_cliente t
		LEFT JOIN cliente c ON t.id_cliente = c.id
		WHERE t.id_tarjeta = ?`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&tarjeta.IDTarjeta,
		&tarjeta.IDCliente,
		&tarjeta.NumeroTarjeta,
		&tarjeta.Tipo,
		&cliente.Nombre,
		&cliente.Email,
		&cliente.Telefono,
		&cliente.Login,
		&cliente.Password,
		&cliente.CodigoSeguridad,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tarjeta no encontrada")
		}
		log.Printf("FindTarjetaCliente DB error: %v", err)
		return nil, err
	}
	tarjeta.Cliente = &cliente
	return &tarjeta, nil
}

// UpdateTarjetaCliente actualiza los datos de una tarjeta
func (r *mutationResolver) UpdateTarjetaCliente(ctx context.Context, input model.UpdateTarjetaCliente) (*model.TarjetaCliente, error) {
	log.Printf("UpdateTarjetaCliente called with input: %+v", input)

	if input.IDTarjeta == "" {
		return nil, fmt.Errorf("el ID de la tarjeta es obligatorio para la actualización")
	}

	query := `UPDATE tarjeta_cliente SET 
		id_cliente = COALESCE(?, id_cliente), 
		numero_tarjeta = COALESCE(?, numero_tarjeta), 
		tipo = COALESCE(?, tipo)
	WHERE id_tarjeta = ?`

	_, err := r.DB.ExecContext(ctx, query, input.IDCliente, input.NumeroTarjeta, input.Tipo, input.IDTarjeta)
	if err != nil {
		log.Printf("UpdateTarjetaCliente DB error: %v", err)
		return nil, err
	}

	var tarjeta model.TarjetaCliente
	selectQuery := `SELECT id_tarjeta, id_cliente, numero_tarjeta, tipo
	FROM tarjeta_cliente WHERE id_tarjeta = ?`

	err = r.DB.QueryRowContext(ctx, selectQuery, input.IDTarjeta).Scan(
		&tarjeta.IDTarjeta, &tarjeta.IDCliente, &tarjeta.NumeroTarjeta, &tarjeta.Tipo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tarjeta no encontrada")
		}
		log.Printf("Error al leer tarjeta actualizada: %v", err)
		return nil, err
	}

	return &tarjeta, nil
}

// CreateTarjetaCliente crea una nueva tarjeta
func (r *mutationResolver) CreateTarjetaCliente(ctx context.Context, input model.NewTarjetaCliente) (*model.TarjetaCliente, error) {
	log.Printf("CreateTarjetaCliente called with input: %+v", input)

	if input.IDCliente == "" {
		return nil, fmt.Errorf("id_cliente es obligatorio")
	}
	if input.NumeroTarjeta == "" || input.Tipo == "" {
		return nil, fmt.Errorf("numero_tarjeta y tipo son obligatorios")
	}

	res, err := r.DB.ExecContext(ctx, `
		INSERT INTO tarjeta_cliente (id_cliente, numero_tarjeta, tipo)
		VALUES (?, ?, ?)`,
		input.IDCliente, input.NumeroTarjeta, input.Tipo)
	if err != nil {
		log.Printf("CreateTarjetaCliente DB error: %v", err)
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener el id_tarjeta: %v", err)
	}

	tarjeta := &model.TarjetaCliente{
		IDTarjeta:     fmt.Sprintf("%d", id),
		IDCliente:     input.IDCliente,
		NumeroTarjeta: input.NumeroTarjeta,
		Tipo:          input.Tipo,
	}

	return tarjeta, nil
}

// KillTarjetaCliente elimina una tarjeta
func (r *mutationResolver) KillTarjetaCliente(ctx context.Context, id string) (bool, error) {
	log.Printf("KillTarjetaCliente called with id: %s", id)
	if id == "" {
		return false, fmt.Errorf("el ID de la tarjeta es obligatorio para eliminar")
	}

	_, err := r.DB.ExecContext(ctx, "DELETE FROM tarjeta_cliente WHERE id_tarjeta = ?", id)
	if err != nil {
		log.Printf("KillTarjetaCliente DB error: %v", err)
		return false, err
	}
	return true, nil
}