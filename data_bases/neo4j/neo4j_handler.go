package neo4j

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var driver neo4j.DriverWithContext

// ConectarNeo4j inicializa el driver para comunicarse con Neo4j Desktop
func ConectarNeo4j() {
	ctx := context.Background()
	uri := "bolt://localhost:7687"
	// cambia por las tuyas, pero tambien lo podemos dejar igual xd
	usuario := "neo4j"
	contrasena := "123456789"

	var err error
	driver, err = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(usuario, contrasena, ""))
	if err != nil {
		log.Fatalf("Error al crear el driver de Neo4j: %v", err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		log.Fatalf("No se pudo establecer conexion con Neo4j: %v", err)
	}

	log.Println("Instancia de Neo4j conectada exitosamente en el puerto :7687")
}

// GetNeo4j retorna el driver activo para ser utilizado en los Resolvers
func GetNeo4j() neo4j.DriverWithContext {
	return driver
}