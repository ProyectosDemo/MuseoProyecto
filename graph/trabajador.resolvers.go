package graph

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/graph/model"
)

// Trabajadores lista de los trabajadores
func (r *queryResolver) Trabajadores(ctx context.Context, limit *int32, offset *int32) ([]*model.Trabajador, error) {

	rows, err := r.DB.QueryContext(ctx, "SELECT id_trabajador, nombre, login, password, admin FROM trabajador LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		log.Printf("Trabajadores DB error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var clientes []*model.Trabajador
	for rows.Next() {
		var cliente model.Trabajador
		err := rows.Scan(&cliente.ID, &cliente.Nombre, &cliente.Login, &cliente.Password, &cliente.Admin)
		if err != nil {
			log.Printf("Trabajadores scan error: %v", err)
			return nil, err
		}
		clientes = append(clientes, &cliente)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Trabajadores rows error: %v", err)
		return nil, err
	}
	return clientes, nil
}

// FindTrabajador is the resolver for the findTrabajador field.
func (r *queryResolver) FindTrabajador(ctx context.Context, id string) ([]*model.Trabajador, error) {
	log.Printf("FindTrabajador llamado con id: %s", id)

	var trabajador model.Trabajador
	query := `SELECT id_trabajador, nombre, login, password, admin FROM trabajador WHERE id_trabajador = ?`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&trabajador.ID, &trabajador.Nombre, &trabajador.Login, &trabajador.Password, &trabajador.Admin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trabajador no encontrado")
		}
		log.Printf("TrabajadorCliente DB error: %v", err)
		return nil, err
	}
	return []*model.Trabajador{&trabajador}, nil
}

// UpdateTrabajador is the resolver for the updateTrabajador field.
func (r *mutationResolver) UpdateTrabajador(ctx context.Context, input model.UpdateTrabajador) (*model.Trabajador, error) {
	log.Printf("UpdateTrabajador llamado con input: %+v", input)

	if input.ID == "" {
		return nil, fmt.Errorf("el ID del trabajador es obligatorio para la actualización")
	}

	query := `UPDATE trabajador SET nombre = COALESCE(?, nombre), login = COALESCE(?, login), password = COALESCE(?, password), admin = COALESCE(?, admin) WHERE id_trabajador = ?`

	_, err := r.DB.ExecContext(ctx, query, input.Nombre, input.Login, input.Password, input.Admin, input.ID)
	if err != nil {
		log.Printf("UpdateTrabajador DB error: %v", err)
		return nil, err
	}

	var trabajador model.Trabajador
	selectQuery := `SELECT id_trabajador, nombre, login, password, admin
	FROM trabajador WHERE id_trabajador = ?`

	err = r.DB.QueryRowContext(ctx, selectQuery, input.ID).Scan(
		&trabajador.ID, &trabajador.Nombre, &trabajador.Login, &trabajador.Password, &trabajador.Admin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trabajador no encontrado")
		}
		log.Printf("Error al leer trabajador actualizado: %v", err)
		return nil, err
	}

	return &trabajador, nil
}

// CreateTrabajador is the resolver for the createTrabajador field.
func (r *mutationResolver) CreateTrabajador(ctx context.Context, input model.NewTrabajador) (*model.Trabajador, error) {
	log.Printf("CreateTrabajador called with input: %+v", input)
	if input.Nombre == "" || input.Login == "" || input.Password == "" {
		return nil, fmt.Errorf("todos los campos de input son obligatorios")
	}
	_, err := r.DB.ExecContext(ctx, "INSERT INTO trabajador (nombre, login, password, admin) VALUES (?, ?, ?, ?)", input.Nombre, input.Login, input.Password, input.Admin)
	if err != nil {
		log.Printf("CreateTrabajador DB error: %v", err)
		return nil, err
	}
	trab := &model.Trabajador{
		Nombre:   input.Nombre,
		Login:    input.Login,
		Password: input.Password,
		Admin:    input.Admin,
	}
	return trab, nil
}

// KillTrabajador is the resolver for the killTrabajador field.
func (r *mutationResolver) KillTrabajador(ctx context.Context, id string) (bool, error) {
	// log de la llamada para depuración
	log.Printf("KillCliente called with id: %s", id)
	if id == "" {
		return false, fmt.Errorf("el ID del trabajador es obligatorio para eliminar")
	}
	// ejecucion de la consulta
	_, err := r.DB.ExecContext(ctx, "DELETE FROM trabajador WHERE id_trabajador = ?", id)
	if err != nil {
		log.Printf("KillTrabajador DB error: %v", err)
		return false, err
	}
	return true, nil
}

func (r *queryResolver) LoginTrabajador(ctx context.Context, login string, password string) (*model.LoginResponseTrabajador, error) {
    var trabajador model.Trabajador
    query := `SELECT id_trabajador, nombre, password, admin FROM trabajador WHERE login = ?`
    err := r.DB.QueryRowContext(ctx, query, login).Scan(&trabajador.ID, &trabajador.Nombre, &trabajador.Password, &trabajador.Admin)
    if err != nil {
        if err == sql.ErrNoRows {
            return &model.LoginResponseTrabajador{Success: false}, nil
        }
        return nil, err
    }

    if trabajador.Password != password {
        return &model.LoginResponseTrabajador{Success: false}, nil
    }

    idPtr := &trabajador.ID
    nombrePtr := &trabajador.Nombre
    adminPtr := &trabajador.Admin

    return &model.LoginResponseTrabajador{
        Success: true,
        ID:      idPtr,
        Nombre:  nombrePtr,
        Admin:   adminPtr,
    }, nil
}