package graph

import (
	"context"
	"fmt"
	"log"
	"main/graph/model"
	"main/graph/models_mongodb"
	"main/data_bases/mongodb"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Helper para mapear la estructura de MongoDB directamente al modelo de GraphQL
func toGraphQLArtista(ma *models_mongodb.ArtistaMongo) *model.Artista {
	return &model.Artista{
		ID:              fmt.Sprintf("%d", ma.ID),
		Nombre:          ma.Nombre,
		FechaNacimiento: ma.FechaNacimiento,
		Nacionalidad:    ma.Nacionalidad,
		Biografia:       ma.Biografia,
		Foto:            ma.Foto,
	}
}

// Artistas trae todos los artistas usando paginacion basica de MongoDB
func (r *queryResolver) Artistas(ctx context.Context, limit *int32, offset *int32) ([]*model.Artista, error) {
	findOptions := options.Find()
	if offset != nil { findOptions.SetSkip(int64(*offset)) }
	if limit != nil {  findOptions.SetLimit(int64(*limit)) }

	cursor, err := mongodb.GetMongoDB().Collection("artista").Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Printf("Artistas DB error: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var resultados []models_mongodb.ArtistaMongo
	if err := cursor.All(ctx, &resultados); err != nil {
		log.Printf("Artistas cursor error: %v", err)
		return nil, err
	}

	var artistas []*model.Artista
	for _, ma := range resultados {
		artistas = append(artistas, toGraphQLArtista(&ma))
	}
	return artistas, nil
}

// FindArtista busca un solo artista por su ID numerico unico (_id)
func (r *queryResolver) FindArtista(ctx context.Context, id string) ([]*model.Artista, error) {
	log.Printf("FindArtista called with id: %s", id)

	idNumerico, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("el ID debe ser un numero valido")
	}

	var ma models_mongodb.ArtistaMongo
	err = mongodb.GetMongoDB().Collection("artista").FindOne(ctx, bson.D{{Key: "_id", Value: int32(idNumerico)}}).Decode(&ma)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("artista no encontrado")
		}
		log.Printf("FindArtista DB error: %v", err)
		return nil, err
	}

	return []*model.Artista{toGraphQLArtista(&ma)}, nil
}

// UpdateArtista modifica campos de forma dinamica usando bson.M de manera compacta
func (r *mutationResolver) UpdateArtista(ctx context.Context, input model.UpdateArtista) (*model.Artista, error) {
	log.Printf("UpdateArtista called with input: %+v", input)

	if input.ID == "" {
		return nil, fmt.Errorf("el ID del artista es obligatorio para la actualizacion")
	}

	idNumerico, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("ID invalido")
	}

	// Creamos un mapa vacio y agregamos propiedades solo si vienen informadas desde GraphQL
	updateFields := bson.M{}
	if input.Nombre != nil {          updateFields["nombre"] = *input.Nombre }
	if input.FechaNacimiento != nil { updateFields["fecha_nacimiento"] = *input.FechaNacimiento }
	if input.Nacionalidad != nil {    updateFields["nacionalidad"] = *input.Nacionalidad }
	if input.Biografia != nil {       updateFields["biografia"] = *input.Biografia }
	if input.Foto != nil {            updateFields["foto"] = *input.Foto }

	coll := mongodb.GetMongoDB().Collection("artista")
	if len(updateFields) > 0 {
		_, err := coll.UpdateOne(ctx, bson.M{"_id": int32(idNumerico)}, bson.M{"$set": updateFields})
		if err != nil {
			log.Printf("UpdateArtista DB error: %v", err)
			return nil, err
		}
	}

	// Recupera el documento modificado para retornar la informacion actualizada
	var ma models_mongodb.ArtistaMongo
	if err := coll.FindOne(ctx, bson.M{"_id": int32(idNumerico)}).Decode(&ma); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("artista no encontrado")
		}
		log.Printf("Error al leer artista actualizado: %v", err)
		return nil, err
	}

	return toGraphQLArtista(&ma), nil
}

// CreateArtista registra un artista nuevo calculando el ID autoincremental en Mongo
func (r *mutationResolver) CreateArtista(ctx context.Context, input model.NewArtista) (*model.Artista, error) {
	log.Printf("CreateArtista called with input: %+v", input)

	if input.Nombre == "" || input.FechaNacimiento == "" || input.Nacionalidad == "" || input.Biografia == "" || input.Foto == "" {
		return nil, fmt.Errorf("todos los campos de input son obligatorios")
	}

	db := mongodb.GetMongoDB()

	// Calculo del siguiente ID numerico autoincremental ordenando por el ultimo _id
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})
	var ultimoArtista struct{ ID int32 `bson:"_id"` }
	var nuevoID int32 = 1
	if err := db.Collection("artista").FindOne(ctx, bson.D{}, opts).Decode(&ultimoArtista); err == nil {
		nuevoID = ultimoArtista.ID + 1
	}

	nuevoArtistaDoc := models_mongodb.ArtistaMongo{
		ID:              nuevoID,
		Nombre:          input.Nombre,
		FechaNacimiento: input.FechaNacimiento,
		Nacionalidad:    input.Nacionalidad,
		Biografia:       input.Biografia,
		Foto:            input.Foto,
	}

	if _, err := db.Collection("artista").InsertOne(ctx, nuevoArtistaDoc); err != nil {
		log.Printf("CreateArtista DB error: %v", err)
		return nil, err
	}

	return toGraphQLArtista(&nuevoArtistaDoc), nil
}

// KillArtista elimina un documento de artista por completo mediante su clave unica _id
func (r *mutationResolver) KillArtista(ctx context.Context, id string) (bool, error) {
	log.Printf("KillArtista called with id: %s", id)
	if id == "" {
		return false, fmt.Errorf("el ID del artista es obligatorio para eliminar")
	}

	idNumerico, err := strconv.Atoi(id)
	if err != nil {
		return false, fmt.Errorf("ID invalido")
	}

	res, err := mongodb.GetMongoDB().Collection("artista").DeleteOne(ctx, bson.M{"_id": int32(idNumerico)})
	if err != nil {
		log.Printf("KillArtista DB error: %v", err)
		return false, err
	}

	return res.DeletedCount > 0, nil
}