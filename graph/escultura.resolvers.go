package graph

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/graph/model"
)

// Esculturas lista todas las esculturas
func (r *queryResolver) Esculturas(ctx context.Context, limit *int32, offset *int32) ([]*model.Escultura, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT e.id_obra, e.material, e.peso, e.dimensiones,
		        o.nombre, o.precio_obra, o.fecha_creacion, o.status, o.foto,
		        a.id_artista, a.nombre,
		        g.id_genero, g.nombre
		 FROM escultura e
		 JOIN obra o ON e.id_obra = o.id_obra
		 JOIN artista a ON o.id_artista = a.id_artista
		 JOIN genero g ON o.id_genero = g.id_genero
		 LIMIT ? OFFSET ?`,
		limit, offset)
	if err != nil {
		log.Printf("Esculturas DB error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var esculturas []*model.Escultura
	for rows.Next() {
		var esc model.Escultura
		var obra model.Obra
		var artista model.Artista
		var genero model.Genero

		err := rows.Scan(&esc.IDObra, &esc.Material, &esc.Peso, &esc.Dimensiones,
			&obra.Nombre, &obra.Precio, &obra.FechaCreacion, &obra.Status, &obra.Foto,
			&artista.ID, &artista.Nombre,
			&genero.ID, &genero.Nombre)
		if err != nil {
			log.Printf("Escultura scan error: %v", err)
			return nil, err
		}

		obra.ID = esc.IDObra
		obra.Artista = &artista
		obra.Genero = &genero

		esc.Obra = &obra

		esculturas = append(esculturas, &esc)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Esculturas rows error: %v", err)
		return nil, err
	}

	return esculturas, nil
}

// FindEscultura busca una escultura por id_obra
func (r *queryResolver) FindEscultura(ctx context.Context, id string) (*model.Escultura, error) {
	var esc model.Escultura
	var obra model.Obra
	var artista model.Artista
	var genero model.Genero

	query := `
	SELECT e.id_obra, e.material, e.peso, e.dimensiones,
	       o.nombre, o.precio_obra, o.fecha_creacion, o.status, o.foto,
	       a.id_artista, a.nombre,
	       g.id_genero, g.nombre
	FROM escultura e
	JOIN obra o ON e.id_obra = o.id_obra
	JOIN artista a ON o.id_artista = a.id_artista
	JOIN genero g ON o.id_genero = g.id_genero
	WHERE e.id_obra = ?`

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&esc.IDObra, &esc.Material, &esc.Peso, &esc.Dimensiones,
		&obra.Nombre, &obra.Precio, &obra.FechaCreacion, &obra.Status, &obra.Foto,
		&artista.ID, &artista.Nombre,
		&genero.ID, &genero.Nombre,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("escultura no encontrada")
		}
		return nil, err
	}

	obra.ID = esc.IDObra
	obra.Artista = &artista
	obra.Genero = &genero
	esc.Obra = &obra

	return &esc, nil
}

// CreateEscultura crea una escultura asociada a una obra existente
func (r *mutationResolver) CreateEscultura(ctx context.Context, input model.NewEscultura) (*model.Escultura, error) {
	if input.IDObra == "" {
		return nil, fmt.Errorf("el id_obra es obligatorio")
	}
	if input.Material == "" || input.Peso <= 0 || input.Dimensiones == "" {
		return nil, fmt.Errorf("material, peso y dimensiones son obligatorios y peso debe ser > 0")
	}

	// Verificar que la obra exista
	var exists int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM obra WHERE id_obra = ?", input.IDObra).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, fmt.Errorf("no existe la obra con id %s", input.IDObra)
	}

	insertQuery := `INSERT INTO escultura (id_obra, material, peso, dimensiones) VALUES (?, ?, ?, ?)`
	_, err = r.DB.ExecContext(ctx, insertQuery, input.IDObra, input.Material, input.Peso, input.Dimensiones)
	if err != nil {
		log.Printf("CreateEscultura DB error: %v", err)
		return nil, err
	}

	return &model.Escultura{
		IDObra:     input.IDObra,
		Material:   input.Material,
		Peso:       input.Peso,
		Dimensiones: input.Dimensiones,
	}, nil
}

func (r *mutationResolver) UpdateEscultura(ctx context.Context, input model.UpdateEscultura) (*model.Escultura, error) {
	log.Printf("UpdateEscultura llamado con input: %+v", input)

	if input.IDObra == "" {
		return nil, fmt.Errorf("el ID de la obra es obligatorio para la actualización")
	}

	query := `UPDATE escultura SET 
		material = COALESCE(?, material),
		peso = COALESCE(?, peso),
		dimensiones = COALESCE(?, dimensiones)
		WHERE id_obra = ?`

	_, err := r.DB.ExecContext(ctx, query, input.Material, input.Peso, input.Dimensiones, input.IDObra)
	if err != nil {
		log.Printf("UpdateEscultura DB error: %v", err)
		return nil, err
	}

	// Leer la escultura actualizada
	var esc model.Escultura
	var obra model.Obra
	var artista model.Artista
	var genero model.Genero

	selectQuery := `
	SELECT e.id_obra, e.material, e.peso, e.dimensiones,
	       o.nombre, o.precio_obra, o.fecha_creacion, o.status, o.foto,
	       a.id_artista, a.nombre,
	       g.id_genero, g.nombre
	FROM escultura e
	JOIN obra o ON e.id_obra = o.id_obra
	JOIN artista a ON o.id_artista = a.id_artista
	JOIN genero g ON o.id_genero = g.id_genero
	WHERE e.id_obra = ?`

	err = r.DB.QueryRowContext(ctx, selectQuery, input.IDObra).Scan(
		&esc.IDObra, &esc.Material, &esc.Peso, &esc.Dimensiones,
		&obra.Nombre, &obra.Precio, &obra.FechaCreacion, &obra.Status, &obra.Foto,
		&artista.ID, &artista.Nombre,
		&genero.ID, &genero.Nombre,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("escultura no encontrada")
		}
		log.Printf("Error al leer escultura actualizada: %v", err)
		return nil, err
	}

	obra.ID = esc.IDObra
	obra.Artista = &artista
	obra.Genero = &genero
	esc.Obra = &obra

	return &esc, nil
}