package graph

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/graph/model"
)

// Clientes lista de los clientes
func (r *queryResolver) Clientes(ctx context.Context, limit *int32, offset *int32) ([]*model.Cliente, error) {

	rows, err := r.DB.QueryContext(ctx, "SELECT id_cliente, nombre, email, telefono, login, password, codigo_seguridad FROM cliente LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		log.Printf("Clientes DB error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var clientes []*model.Cliente
	for rows.Next() {
		var cliente model.Cliente
		err := rows.Scan(&cliente.ID, &cliente.Nombre, &cliente.Email, &cliente.Telefono, &cliente.Login, &cliente.Password, &cliente.CodigoSeguridad)
		if err != nil {
			log.Printf("Clientes scan error: %v", err)
			return nil, err
		}
		clientes = append(clientes, &cliente)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Clientes rows error: %v", err)
		return nil, err
	}
	return clientes, nil
}

// FindCliente is the resolver for the findCliente field.
func (r *queryResolver) FindCliente(ctx context.Context, id string) ([]*model.Cliente, error) {
	log.Printf("FindCliente called with id: %s", id)

	var cliente model.Cliente
	query := `SELECT id_cliente, nombre, email, telefono, login, password, codigo_seguridad FROM cliente WHERE id_cliente = ?`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&cliente.ID, &cliente.Nombre, &cliente.Email, &cliente.Telefono,
		&cliente.Login, &cliente.Password, &cliente.CodigoSeguridad,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cliente no encontrado")
		}
		log.Printf("FindCliente DB error: %v", err)
		return nil, err
	}
	return []*model.Cliente{&cliente}, nil
}

// UpdateCliente is the resolver for the updateCliente field.
func (r *mutationResolver) UpdateCliente(ctx context.Context, input model.UpdateCliente) (*model.Cliente, error) {
	log.Printf("UpdateCliente called with input: %+v", input)

	if input.ID == "" {
		return nil, fmt.Errorf("el ID del cliente es obligatorio para la actualización")
	}

	query := `UPDATE cliente SET 
		nombre = COALESCE(?, nombre), 
		email = COALESCE(?, email), 
		telefono = COALESCE(?, telefono), 
		login = COALESCE(?, login), 
		password = COALESCE(?, password), 
		codigo_seguridad = COALESCE(?, codigo_seguridad) 
	WHERE id_cliente = ?`

	_, err := r.DB.ExecContext(ctx, query, input.Nombre, input.Email, input.Telefono, input.Login, input.Password, input.CodigoSeguridad, input.ID)
	if err != nil {
		log.Printf("UpdateCliente DB error: %v", err)
		return nil, err
	}

	var cliente model.Cliente
	selectQuery := `SELECT id_cliente, nombre, email, telefono, login, password, codigo_seguridad 
	FROM cliente WHERE id_cliente = ?`

	err = r.DB.QueryRowContext(ctx, selectQuery, input.ID).Scan(
		&cliente.ID, &cliente.Nombre, &cliente.Email, &cliente.Telefono,
		&cliente.Login, &cliente.Password, &cliente.CodigoSeguridad,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cliente no encontrado")
		}
		log.Printf("Error al leer cliente actualizado: %v", err)
		return nil, err
	}

	return &cliente, nil
}

// CreateCliente is the resolver for the createCliente field.
func (r *mutationResolver) CreateCliente(ctx context.Context, input model.NewCliente) (*model.Cliente, error) {
    log.Printf("CreateCliente called with input: %+v", input)

    if input.Nombre == "" || input.Email == "" || input.Telefono == "" || input.Login == "" || input.Password == "" || input.CodigoSeguridad == "" {
        return nil, fmt.Errorf("todos los campos de input son obligatorios")
    }

    // Ejecutar el INSERT y capturar resultado
    res, err := r.DB.ExecContext(ctx, "INSERT INTO cliente (nombre, email, telefono, login, password, codigo_seguridad) VALUES (?, ?, ?, ?, ?, ?)",
        input.Nombre, input.Email, input.Telefono, input.Login, input.Password, input.CodigoSeguridad)
    if err != nil {
        log.Printf("CreateCliente DB error: %v", err)
        return nil, err
    }

    // Obtener el ID generado automáticamente
    id, err := res.LastInsertId()
    if err != nil {
        log.Printf("Error obteniendo ID del cliente: %v", err)
        return nil, err
    }

    // Devolver cliente con ID
    cliente := &model.Cliente{
        ID:              fmt.Sprintf("%d", id), // importante: el ID como string
        Nombre:          input.Nombre,
        Email:           input.Email,
        Telefono:        input.Telefono,
        Login:           input.Login,
        Password:        input.Password,
        CodigoSeguridad: input.CodigoSeguridad,
    }

    return cliente, nil
}

func (r *queryResolver) LoginCliente(ctx context.Context, email string, password string) (*model.LoginResponseCliente, error) {
    var cliente model.Cliente
    query := `SELECT id_cliente, nombre, password FROM cliente WHERE email = ?`
    err := r.DB.QueryRowContext(ctx, query, email).Scan(&cliente.ID, &cliente.Nombre, &cliente.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            return &model.LoginResponseCliente{Success: false}, nil
        }
        return nil, err
    }

    if cliente.Password != password {
        return &model.LoginResponseCliente{Success: false}, nil
    }

    // Crear punteros
    idPtr := &cliente.ID
    nombrePtr := &cliente.Nombre

    return &model.LoginResponseCliente{
        Success: true,
        ID:      idPtr,
        Nombre:  nombrePtr,
    }, nil
}