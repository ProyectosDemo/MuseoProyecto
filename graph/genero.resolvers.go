package graph

import (
	"context"
	"fmt"
	"log"
	"main/graph/model"
	"main/graph/models_mongodb"
	"main/mongodb"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Helper para mapear la estructura de MongoDB directamente al modelo de GraphQL
func toGraphQLGenero(mg *models_mongodb.GeneroMongo) *model.Genero {
	return &model.Genero{
		ID:     fmt.Sprintf("%d", mg.ID),
		Nombre: mg.Nombre,
	}
}

// Genero trae todos los generos usando paginacion basica de MongoDB
func (r *queryResolver) Genero(ctx context.Context, limit *int32, offset *int32) ([]*model.Genero, error) {
	findOptions := options.Find()
	if offset != nil { findOptions.SetSkip(int64(*offset)) }
	if limit != nil {  findOptions.SetLimit(int64(*limit)) }

	cursor, err := mongodb.GetMongoDB().Collection("genero").Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Printf("Generos DB error: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var resultados []models_mongodb.GeneroMongo
	if err := cursor.All(ctx, &resultados); err != nil {
		log.Printf("Generos cursor error: %v", err)
		return nil, err
	}

	var generos []*model.Genero
	for _, mg := range resultados {
		generos = append(generos, toGraphQLGenero(&mg))
	}
	return generos, nil
}

// FindGenero busca un solo genero por su ID numerico unico (_id)
func (r *queryResolver) FindGenero(ctx context.Context, id string) ([]*model.Genero, error) {
	log.Printf("FindGenero llamado con id: %s", id)

	idNumerico, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("el ID debe ser un numero valido")
	}

	var mg models_mongodb.GeneroMongo
	err = mongodb.GetMongoDB().Collection("genero").FindOne(ctx, bson.M{"_id": int32(idNumerico)}).Decode(&mg)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("genero no encontrado")
		}
		log.Printf("Genero DB error: %v", err)
		return nil, err
	}

	return []*model.Genero{toGraphQLGenero(&mg)}, nil
}

// UpdateGenero modifica el nombre de forma dinamica usando bson.M
func (r *mutationResolver) UpdateGenero(ctx context.Context, input model.UpdateGenero) (*model.Genero, error) {
	log.Printf("UpdateGenero llamado con input: %+v", input)

	if input.ID == "" {
		return nil, fmt.Errorf("el ID del genero es obligatorio para la actualizacion")
	}

	idNumerico, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("ID invalido")
	}

	updateFields := bson.M{}
	if input.Nombre != nil {
		updateFields["nombre"] = *input.Nombre
	}

	coll := mongodb.GetMongoDB().Collection("genero")
	if len(updateFields) > 0 {
		_, err := coll.UpdateOne(ctx, bson.M{"_id": int32(idNumerico)}, bson.M{"$set": updateFields})
		if err != nil {
			log.Printf("UpdateGenero DB error: %v", err)
			return nil, err
		}
	}

	// Recupera el documento modificado para retornar la informacion actualizada
	var mg models_mongodb.GeneroMongo
	if err := coll.FindOne(ctx, bson.M{"_id": int32(idNumerico)}).Decode(&mg); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("genero no encontrado")
		}
		log.Printf("Error al leer genero actualizado: %v", err)
		return nil, err
	}

	return toGraphQLGenero(&mg), nil
}

// CreateGenero registra un genero nuevo calculando el ID autoincremental en Mongo
func (r *mutationResolver) CreateGenero(ctx context.Context, nombre string) (*model.Genero, error) {
	log.Printf("CreateGenero called with nombre: %s", nombre)
	if nombre == "" {
		return nil, fmt.Errorf("el nombre es obligatorio")
	}

	db := mongodb.GetMongoDB()

	// Calculo del siguiente ID numerico autoincremental ordenando por el ultimo _id
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})
	var ultimoGenero struct{ ID int32 `bson:"_id"` }
	var nuevoID int32 = 1
	if err := db.Collection("genero").FindOne(ctx, bson.D{}, opts).Decode(&ultimoGenero); err == nil {
		nuevoID = ultimoGenero.ID + 1
	}

	nuevoGeneroDoc := models_mongodb.GeneroMongo{
		ID:     nuevoID,
		Nombre: nombre,
	}

	if _, err := db.Collection("genero").InsertOne(ctx, nuevoGeneroDoc); err != nil {
		log.Printf("CreateGenero DB error: %v", err)
		return nil, err
	}

	return toGraphQLGenero(&nuevoGeneroDoc), nil
}

// KillGenero elimina un documento de genero por completo mediante su clave unica _id
func (r *mutationResolver) KillGenero(ctx context.Context, id string) (bool, error) {
	log.Printf("KillGenero called with id: %s", id)
	if id == "" {
		return false, fmt.Errorf("el ID del genero es obligatorio para eliminar")
	}

	idNumerico, err := strconv.Atoi(id)
	if err != nil {
		return false, fmt.Errorf("ID invalido")
	}

	res, err := mongodb.GetMongoDB().Collection("genero").DeleteOne(ctx, bson.M{"_id": int32(idNumerico)})
	if err != nil {
		log.Printf("KillGenero DB error: %v", err)
		return false, err
	}

	return res.DeletedCount > 0, nil
}