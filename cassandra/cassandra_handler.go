package cassandra

import (
	"log"
	"time"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func ConectarCassandra() {
	// Configura el clúster usando tu localhost
	cluster := gocql.NewCluster("127.0.0.1")
	
	// Apuntamos al nuevo Keyspace oficial del museo
	cluster.Keyspace = "museo_cassandra" 
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5 * time.Second
	
	// Forzamos el protocolo nativo v4 para Cassandra 3.11
	cluster.ProtoVersion = 4 

	// Crear la sesión de conexión
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Error al conectar a Cassandra: %v", err)
	}

	Session = session
	log.Println("¡Conexión exitosa a Cassandra (museo_cassandra)!")
}

// GetCassandra devuelve la sesión activa para usarla en los resolvers
func GetCassandra() *gocql.Session {
	return Session
}