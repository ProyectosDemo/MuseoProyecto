package graph

import (
	"database/sql"

	"github.com/gocql/gocql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/mongo" // NUEVO: Importamos el driver de MongoDB
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	DB      *sql.DB
	MongoDB *mongo.Database
	Cassandra *gocql.Session
	Neo4j     neo4j.DriverWithContext
}