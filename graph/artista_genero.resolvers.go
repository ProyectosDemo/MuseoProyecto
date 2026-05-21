package graph

import (
	"context"
	"fmt"
	"log"
	"main/graph/model"
)

// ArtistaGenero lista todas las relaciones artista-genero con detalles
func (r *queryResolver) ArtistaGenero(ctx context.Context, limit *int32, offset *int32) ([]*model.ArtistaGenero, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT ag.id_artista, ag.id_genero,
		       a.nombre, a.fecha_nacimiento, a.nacionalidad, a.biografia, a.foto,
		       g.nombre
		FROM artista_genero ag
		LEFT JOIN artista a ON ag.id_artista = a.id_artista
		LEFT JOIN genero g ON ag.id_genero = g.id_genero
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		log.Printf("ArtistaGenero DB error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var rels []*model.ArtistaGenero
	for rows.Next() {
		var ag model.ArtistaGenero
		var artista model.Artista
		var genero model.Genero
		err := rows.Scan(
			&ag.IDArtista,
			&ag.IDGenero,
			&artista.Nombre,
			&artista.FechaNacimiento,
			&artista.Nacionalidad,
			&artista.Biografia,
			&artista.Foto,
			&genero.Nombre,
		)
		if err != nil {
			log.Printf("ArtistaGenero scan error: %v", err)
			return nil, err
		}
		ag.Artista = &artista
		ag.Genero = &genero
		rels = append(rels, &ag)
	}
	if err = rows.Err(); err != nil {
		log.Printf("ArtistaGenero rows error: %v", err)
		return nil, err
	}
	return rels, nil
}

// CreateArtistaGenero inserta una nueva relación artista-genero
func (r *mutationResolver) CreateArtistaGenero(ctx context.Context, input model.NewArtistaGenero) (*model.ArtistaGenero, error) {
	if input.IDArtista == "" || input.IDGenero == "" {
		return nil, fmt.Errorf("id_artista y id_genero son obligatorios")
	}

	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO artista_genero (id_artista, id_genero) VALUES (?, ?)`,
		input.IDArtista, input.IDGenero)
	if err != nil {
		log.Printf("CreateArtistaGenero DB error: %v", err)
		return nil, err
	}

	ag := &model.ArtistaGenero{
		IDArtista: input.IDArtista,
		IDGenero:  input.IDGenero,
	}

	return ag, nil
}

// KillArtistaGenero elimina una relación artista-genero
func (r *mutationResolver) KillArtistaGenero(ctx context.Context, idArtista string, idGenero string) (bool, error) {
	if idArtista == "" || idGenero == "" {
		return false, fmt.Errorf("id_artista e id_genero son obligatorios")
	}

	_, err := r.DB.ExecContext(ctx, `
		DELETE FROM artista_genero WHERE id_artista = ? AND id_genero = ?`,
		idArtista, idGenero)
	if err != nil {
		log.Printf("KillArtistaGenero DB error: %v", err)
		return false, err
	}
	return true, nil
}

func (r *queryResolver) FindArtistaGeneroByGenero(ctx context.Context, id_genero string) ([]*model.ArtistaGenero, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT ag.id_artista, ag.id_genero,
		       a.id_artista, a.nombre, a.fecha_nacimiento, a.nacionalidad, a.biografia, a.foto,
		       g.nombre
		FROM artista_genero ag
		LEFT JOIN artista a ON ag.id_artista = a.id_artista
		LEFT JOIN genero g ON ag.id_genero = g.id_genero
		WHERE ag.id_genero = ?`, id_genero)
	if err != nil {
		log.Printf("FindArtistaGeneroByGenero DB error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var rels []*model.ArtistaGenero
	for rows.Next() {
		var ag model.ArtistaGenero
		var artista model.Artista
		var genero model.Genero

		err := rows.Scan(
			&ag.IDArtista,
			&ag.IDGenero,
			&artista.ID,
			&artista.Nombre,
			&artista.FechaNacimiento,
			&artista.Nacionalidad,
			&artista.Biografia,
			&artista.Foto,
			&genero.Nombre,
		)
		if err != nil {
			log.Printf("FindArtistaGeneroByGenero scan error: %v", err)
			return nil, err
		}

		ag.Artista = &artista
		ag.Genero = &genero
		rels = append(rels, &ag)
	}

	if err = rows.Err(); err != nil {
		log.Printf("FindArtistaGeneroByGenero rows error: %v", err)
		return nil, err
	}

	return rels, nil
}