package graph

import (
	"database/sql"

	"go.mongodb.org/mongo-driver/mongo" // NUEVO: Importamos el driver de MongoDB
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	DB      *sql.DB         // Tu conexión de MySQL
	MongoDB *mongo.Database // NUEVO: Inyectamos la conexión a MongoDB Atlas
}