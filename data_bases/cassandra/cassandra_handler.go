package cassandra

import (
	"log"
	"time"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func ConectarCassandra() {
	cluster := gocql.NewCluster("127.0.0.1")
	
	// Apuntamos al nuevo Keyspace oficial del museo
	cluster.Keyspace = "museo_cassandra" 
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5 * time.Second
	
	// Forzamos a la estupida de cassandra
	cluster.ProtoVersion = 4 

	// Crear la sesion de conexion
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Error al conectar a Cassandra: %v", err)
	}

	Session = session
	log.Println("¡Conexion exitosa a Cassandra!")
}

// devuelve la session, capaz lo usamos en resolvers
func GetCassandra() *gocql.Session {
	return Session
}