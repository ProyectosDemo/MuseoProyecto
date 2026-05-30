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

// Helper clave para no repetir codigo, mapea la estructura de base de datos directamente al modelo de GraphQL
func toGraphQLEscultura(me *models_mongodb.EsculturaMongo, mo *models_mongodb.ObraMongo) *model.Escultura {
	return &model.Escultura{
		IDObra:      fmt.Sprintf("%d", me.ID), // Su _id es el ID de la obra
		Material:    me.Material,
		Peso:        me.Peso,
		Dimensiones: me.Dimensiones,
		Obra: &model.Obra{
			ID:            fmt.Sprintf("%d", mo.ID),
			Nombre:        mo.Nombre,
			IDArtista:     fmt.Sprintf("%d", mo.IDArtista),
			IDGenero:      fmt.Sprintf("%d", mo.IDGenero),
			Precio:        mo.PrecioObra,
			FechaCreacion: mo.FechaCreacion,
			Status:        model.StatusObra(mo.Status),
			Foto:          mo.Foto,
			Artista: &model.Artista{
				ID:              fmt.Sprintf("%d", mo.Artista.ID),
				Nombre:          mo.Artista.Nombre,
				FechaNacimiento: mo.Artista.FechaNacimiento,
				Nacionalidad:    mo.Artista.Nacionalidad,
				Biografia:       mo.Artista.Biografia,
				Foto:            mo.Artista.Foto,
			},
			Genero: &model.Genero{
				ID:     fmt.Sprintf("%d", mo.Genero.ID),
				Nombre: mo.Genero.Nombre,
			},
		},
	}
}

// Esculturas lista todas las esculturas cruzando los datos con obra ULTIMATE HUMONGOSAURIO
func (r *queryResolver) Esculturas(ctx context.Context, limit *int32, offset *int32) ([]*model.Escultura, error) {
	findOptions := options.Find()
	if offset != nil { findOptions.SetSkip(int64(*offset)) }
	if limit != nil {  findOptions.SetLimit(int64(*limit)) }

	db := mongodb.GetMongoDB()
	cursor, err := db.Collection("escultura").Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Printf("Esculturas DB error: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var resultadosEsculturas []models_mongodb.EsculturaMongo
	if err := cursor.All(ctx, &resultadosEsculturas); err != nil {
		return nil, err
	}

	var esculturas []*model.Escultura
	for _, me := range resultadosEsculturas {
		// Por cada escultura, buscamos recursivamente su obra correspondientE
		var mo models_mongodb.ObraMongo
		err := db.Collection("obra_ultimate").FindOne(ctx, bson.M{"_id": me.ID}).Decode(&mo)
		if err != nil {
			// para evitar romper saltamos en caso de no encontrar
			log.Printf("Advertencia: No se encontro la obra para la escultura con ID %d", me.ID)
			continue
		}
		esculturas = append(esculturas, toGraphQLEscultura(&me, &mo))
	}

	return esculturas, nil
}

// FindEscultura busca una escultura
func (r *queryResolver) FindEscultura(ctx context.Context, id string) (*model.Escultura, error) {
	idNumerico, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("el ID debe ser valido")
	}

	db := mongodb.GetMongoDB()

	// Buscar los datos de la escultura
	var me models_mongodb.EsculturaMongo
	err = db.Collection("escultura").FindOne(ctx, bson.M{"_id": int32(idNumerico)}).Decode(&me)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("escultura no encontrada")
		}
		return nil, err
	}

	// Buscar los datos cruzados de la obra
	var mo models_mongodb.ObraMongo
	err = db.Collection("obra_ultimate").FindOne(ctx, bson.M{"_id": int32(idNumerico)}).Decode(&mo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("se encontro la escultura pero los datos de su obra no existen en obra_ultimate")
		}
		return nil, err
	}

	return toGraphQLEscultura(&me, &mo), nil
}

// CreateEscultura inserta una nueva escultura vinculada al ID de la obra
func (r *mutationResolver) CreateEscultura(ctx context.Context, input model.NewEscultura) (*model.Escultura, error) {
	if input.IDObra == "" {
		return nil, fmt.Errorf("el id_obra es obligatorio")
	}
	if input.Material == "" || input.Peso <= 0 || input.Dimensiones == "" {
		return nil, fmt.Errorf("material, peso y dimensiones son obligatorios; peso debe ser > 0")
	}

	idNumerico, err := strconv.Atoi(input.IDObra)
	if err != nil {
		return nil, fmt.Errorf("ID de obra invalido")
	}

	db := mongodb.GetMongoDB()

	// Validar que la obra exista
	var mo models_mongodb.ObraMongo
	err = db.Collection("obra_ultimate").FindOne(ctx, bson.M{"_id": int32(idNumerico)}).Decode(&mo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no se puede crear la escultura porque no existe la obra con id %s en obra_ultimate", input.IDObra)
		}
		return nil, err
	}

	nuevaEsculturaDoc := models_mongodb.EsculturaMongo{
		ID:          int32(idNumerico), // Su _id coincide con la obra
		Material:    input.Material,
		Peso:        input.Peso,
		Dimensiones: input.Dimensiones,
	}

	_, err = db.Collection("escultura").InsertOne(ctx, nuevaEsculturaDoc)
	if err != nil {
		log.Printf("CreateEscultura DB error: %v", err)
		return nil, err
	}

	return toGraphQLEscultura(&nuevaEsculturaDoc, &mo), nil
}

// UpdateEscultura actualiza las propiedades
func (r *mutationResolver) UpdateEscultura(ctx context.Context, input model.UpdateEscultura) (*model.Escultura, error) {
	log.Printf("UpdateEscultura llamado con input: %+v", input)

	if input.IDObra == "" {
		return nil, fmt.Errorf("el ID de la obra es obligatorio para la actualización")
	}

	idNumerico, err := strconv.Atoi(input.IDObra)
	if err != nil {
		return nil, fmt.Errorf("ID inválido")
	}

	updateFields := bson.M{}
	if input.Material != nil {      updateFields["material"] = *input.Material }
	if input.Peso != nil {          updateFields["peso"] = *input.Peso }
	if input.Dimensiones != nil {  updateFields["dimensiones"] = *input.Dimensiones }

	db := mongodb.GetMongoDB()
	collEscultura := db.Collection("escultura")

	if len(updateFields) > 0 {
		_, err = collEscultura.UpdateOne(ctx, bson.M{"_id": int32(idNumerico)}, bson.M{"$set": updateFields})
		if err != nil {
			log.Printf("UpdateEscultura DB error: %v", err)
			return nil, err
		}
	}

	// Recuperar los datos actualizados de la escultura
	var me models_mongodb.EsculturaMongo
	err = collEscultura.FindOne(ctx, bson.M{"_id": int32(idNumerico)}).Decode(&me)
	if err != nil {
		if err == mongo.ErrNoDocuments { return nil, fmt.Errorf("escultura no encontrada") }
		return nil, err
	}

	// Recuperar los datos estables de la obra para construir la respuesta completa de GraphQL
	var mo models_mongodb.ObraMongo
	err = db.Collection("obra_ultimate").FindOne(ctx, bson.M{"_id": int32(idNumerico)}).Decode(&mo)
	if err != nil {
		return nil, fmt.Errorf("error al recuperar los datos de la obra vinculada: %v", err)
	}

	return toGraphQLEscultura(&me, &mo), nil
}