package graph

import (
	"context"
	"fmt"
	"log"
	"main/graph/model"
	"main/mongodb"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//	Por ahora se queda asi para evitar conflictos con models
type ObraMongo struct {
	ID            int32        `bson:"_id"`
	Nombre        string       `bson:"nombre"`
	IDArtista     int32        `bson:"id_artista"`
	IDGenero      int32        `bson:"id_genero"`
	PrecioObra    int32        `bson:"precio_obra"`
	FechaCreacion string       `bson:"fecha_creacion"`
	Status        string       `bson:"status"`
	Foto          string       `bson:"foto"`
	ArtistaInfo   ArtistaMongo `bson:"artista_info"` 
	GeneroInfo    GeneroMongo  `bson:"genero_info"`  
}

type ArtistaMongo struct {
	ID              int32  `bson:"_id"`
	Nombre          string `bson:"nombre"`
	FechaNacimiento string `bson:"fecha_nacimiento"`
	Nacionalidad    string `bson:"nacionalidad"`
	Biografia       string `bson:"biografia"`
	Foto            string `bson:"foto"`
}

type GeneroMongo struct {
	ID     int32  `bson:"_id"`
	Nombre string `bson:"nombre"`
}

type ObraPlanaMongo struct {
	ID            int32  `bson:"_id"`
	Nombre        string `bson:"nombre"`
	IDArtista     int32  `bson:"id_artista"`
	IDGenero      int32  `bson:"id_genero"`
	PrecioObra    int32  `bson:"precio_obra"`
	FechaCreacion string `bson:"fecha_creacion"`
	Status        string `bson:"status"`
	Foto          string `bson:"foto"`
}

var validStatuses = map[string]bool{
	"DISPONIBLE": true,
	"RESERVADA":  true,
	"VENDIDA":    true,
}

// queries en nosql

func (r *queryResolver) Obras(ctx context.Context, limit *int32, offset *int32) ([]*model.Obra, error) {
	db := mongodb.GetMongoDB()
	collection := db.Collection("obra")

	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "artista"},
			{Key: "localField", Value: "id_artista"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "artista_info"},
		}}},
		{{Key: "$unwind", Value: "$artista_info"}},

		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "genero"},
			{Key: "localField", Value: "id_genero"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "genero_info"},
		}}},
		{{Key: "$unwind", Value: "$genero_info"}},
	}

	if offset != nil {
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: *offset}})
	}
	if limit != nil {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: *limit}})
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Obras Aggregation error: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var resultados []ObraMongo
	if err := cursor.All(ctx, &resultados); err != nil {
		return nil, err
	}

	var obras []*model.Obra
	for _, mo := range resultados {
		obras = append(obras, &model.Obra{
			ID:            fmt.Sprintf("%d", mo.ID),
			Nombre:        mo.Nombre,
			IDArtista:     fmt.Sprintf("%d", mo.IDArtista),
			IDGenero:      fmt.Sprintf("%d", mo.IDGenero),
			Precio:        mo.PrecioObra,
			FechaCreacion: mo.FechaCreacion,
			Status:        model.StatusObra(mo.Status),
			Foto:          mo.Foto,
			Artista: &model.Artista{
				ID:              fmt.Sprintf("%d", mo.ArtistaInfo.ID),
				Nombre:          mo.ArtistaInfo.Nombre,
				FechaNacimiento: mo.ArtistaInfo.FechaNacimiento,
				Nacionalidad:    mo.ArtistaInfo.Nacionalidad,
				Biografia:       mo.ArtistaInfo.Biografia,
				Foto:            mo.ArtistaInfo.Foto,
			},
			Genero: &model.Genero{
				ID:     fmt.Sprintf("%d", mo.GeneroInfo.ID),
				Nombre: mo.GeneroInfo.Nombre,
			},
		})
	}

	return obras, nil
}

func (r *queryResolver) FindObra(ctx context.Context, id string) ([]*model.Obra, error) {
	idNumerico, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("el ID debe ser un numero entero valido")
	}

	db := mongodb.GetMongoDB()
	collection := db.Collection("obra")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: int32(idNumerico)}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "artista"},
			{Key: "localField", Value: "id_artista"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "artista_info"},
		}}},
		{{Key: "$unwind", Value: "$artista_info"}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "genero"},
			{Key: "localField", Value: "id_genero"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "genero_info"},
		}}},
		{{Key: "$unwind", Value: "$genero_info"}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var resultados []ObraMongo
	if err := cursor.All(ctx, &resultados); err != nil {
		return nil, err
	}

	if len(resultados) == 0 {
		return nil, fmt.Errorf("obra no encontrada")
	}

	mo := resultados[0]
	obra := &model.Obra{
		ID:            fmt.Sprintf("%d", mo.ID),
		Nombre:        mo.Nombre,
		IDArtista:     fmt.Sprintf("%d", mo.IDArtista),
		IDGenero:      fmt.Sprintf("%d", mo.IDGenero),
		Precio:        mo.PrecioObra,
		FechaCreacion: mo.FechaCreacion,
		Status:        model.StatusObra(mo.Status),
		Foto:          mo.Foto,
		Artista: &model.Artista{
			ID:              fmt.Sprintf("%d", mo.ArtistaInfo.ID),
			Nombre:          mo.ArtistaInfo.Nombre,
			FechaNacimiento: mo.ArtistaInfo.FechaNacimiento,
			Nacionalidad:    mo.ArtistaInfo.Nacionalidad,
			Biografia:       mo.ArtistaInfo.Biografia,
			Foto:            mo.ArtistaInfo.Foto,
		},
		Genero: &model.Genero{
			ID:     fmt.Sprintf("%d", mo.GeneroInfo.ID),
			Nombre: mo.GeneroInfo.Nombre,
		},
	}

	return []*model.Obra{obra}, nil
}

//mutations en nosql
func (r *mutationResolver) CreateObra(ctx context.Context, input model.NewObra) (*model.Obra, error) {
	log.Printf("CreateObra llamado con input: %+v", input)

	if input.Nombre == "" || input.Foto == "" || input.IDArtista == "" || input.IDGenero == "" {
		return nil, fmt.Errorf("hay campos obligatorios vacios")
	}

	statusStr := string(input.Status)
	if !validStatuses[statusStr] {
		return nil, fmt.Errorf("status invalido")
	}

	db := mongodb.GetMongoDB()
	collection := db.Collection("obra")

	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})
	var ultimaObra struct {
		ID int32 `bson:"_id"`
	}
	var nuevoID int32 = 1

	err := collection.FindOne(ctx, bson.D{}, opts).Decode(&ultimaObra)
	if err == nil {
		nuevoID = ultimaObra.ID + 1
	}

	artistaID, _ := strconv.Atoi(input.IDArtista)
	generoID, _ := strconv.Atoi(input.IDGenero)

	nuevaObraDoc := bson.D{
		{Key: "_id", Value: nuevoID},
		{Key: "nombre", Value: input.Nombre},
		{Key: "id_artista", Value: int32(artistaID)},
		{Key: "id_genero", Value: int32(generoID)},
		{Key: "precio_obra", Value: int32(input.Precio)},
		{Key: "fecha_creacion", Value: input.FechaCreacion},
		{Key: "status", Value: statusStr},
		{Key: "foto", Value: input.Foto},
	}

	_, err = collection.InsertOne(ctx, nuevaObraDoc)
	if err != nil {
		log.Printf("CreateObra Mongo error: %v", err)
		return nil, err
	}

	return &model.Obra{
		ID:            fmt.Sprintf("%d", nuevoID),
		Nombre:        input.Nombre,
		IDArtista:     input.IDArtista,
		IDGenero:      input.IDGenero,
		Precio:        input.Precio,
		FechaCreacion: input.FechaCreacion,
		Status:        input.Status,
		Foto:          input.Foto,
	}, nil
}

func (r *mutationResolver) UpdateObra(ctx context.Context, input model.UpdateObra) (*model.Obra, error) {
	log.Printf("UpdateObra llamado con input: %+v", input)

	if input.ID == "" {
		return nil, fmt.Errorf("el ID es obligatorio")
	}

	idNumerico, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("ID invalido")
	}

	db := mongodb.GetMongoDB()
	collection := db.Collection("obra")
	updateFields := bson.D{}

	if input.Nombre != nil {
		updateFields = append(updateFields, bson.E{Key: "nombre", Value: *input.Nombre})
	}
	if input.IDArtista != nil {
		artID, _ := strconv.Atoi(*input.IDArtista)
		updateFields = append(updateFields, bson.E{Key: "id_artista", Value: int32(artID)})
	}
	if input.IDGenero != nil {
		genID, _ := strconv.Atoi(*input.IDGenero)
		updateFields = append(updateFields, bson.E{Key: "id_genero", Value: int32(genID)})
	}
	if input.Precio != nil {
		updateFields = append(updateFields, bson.E{Key: "precio_obra", Value: int32(*input.Precio)})
	}
	if input.FechaCreacion != nil && *input.FechaCreacion != "" {
		updateFields = append(updateFields, bson.E{Key: "fecha_creacion", Value: *input.FechaCreacion})
	}
	if input.Status != nil {
		statusStr := string(*input.Status)
		if !validStatuses[statusStr] {
			return nil, fmt.Errorf("status invalido")
		}
		updateFields = append(updateFields, bson.E{Key: "status", Value: statusStr})
	}
	if input.Foto != nil {
		updateFields = append(updateFields, bson.E{Key: "foto", Value: *input.Foto})
	}

	if len(updateFields) > 0 {
		_, err = collection.UpdateOne(
			ctx,
			bson.D{{Key: "_id", Value: int32(idNumerico)}},
			bson.D{{Key: "$set", Value: updateFields}},
		)
		if err != nil {
			log.Printf("UpdateObra Mongo error: %v", err)
			return nil, err
		}
	}

	var mo ObraPlanaMongo
	err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: int32(idNumerico)}}).Decode(&mo)
	if err != nil {
		return nil, fmt.Errorf("obra no encontrada tras actualizar")
	}

	return &model.Obra{
		ID:            fmt.Sprintf("%d", mo.ID),
		Nombre:        mo.Nombre,
		IDArtista:     fmt.Sprintf("%d", mo.IDArtista),
		IDGenero:      fmt.Sprintf("%d", mo.IDGenero),
		Precio: 		mo.PrecioObra,
		FechaCreacion: mo.FechaCreacion,
		Status:        model.StatusObra(mo.Status),
		Foto:          mo.Foto,
	}, nil
}

func (r *mutationResolver) KillObra(ctx context.Context, id string) (bool, error) {
	log.Printf("KillObra llamado con id: %s", id)
	if id == "" {
		return false, fmt.Errorf("el ID es obligatorio")
	}

	idNumerico, err := strconv.Atoi(id)
	if err != nil {
		return false, fmt.Errorf("ID invalido")
	}

	db := mongodb.GetMongoDB()
	collection := db.Collection("obra")
	res, err := collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: int32(idNumerico)}})
	if err != nil {
		log.Printf("KillObra Mongo error: %v", err)
		return false, err
	}

	return res.DeletedCount > 0, nil
}