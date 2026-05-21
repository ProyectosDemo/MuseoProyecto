package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var base_datos *mongo.Database

func ConectarMongo() {
	// sin los 30 segundos explota, dejalo asi
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Tu URL limpia (Apuntando a tu base de datos)
	mongoURI := "mongodb+srv://proyectomuseo.igaaxvr.mongodb.net/proyectomuseo?retryWrites=true"

	// Cambias por los tuyos
	credential := options.Credential{
		Username: "fabian",
		Password: "XSzlsRlx9KAL9YEM",
	}

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetAuth(credential)

	var err error
	mongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error al crear cliente de MongoDB: %v", err)
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error de Ping a MongoDB Atlas: %v", err)
	}

	log.Println("¡Conectado exitosamente a MongoDB Atlas!")
	base_datos = mongoClient.Database("Museo")
}

func GetMongoDB() *mongo.Database {
	return base_datos
}