package graph

import (
	"context"
	"database/sql"
	"fmt"
	"log"


	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SincronizarGrafoSprint3 extrae datos de MySQL y MongoDB manejando subdocumentos embebidos
func SincronizarGrafoSprint3(ctx context.Context, dbMySQL *sql.DB, dbMongo *mongo.Database, driverNeo4j neo4j.DriverWithContext) error {
	session := driverNeo4j.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "museoproyecto",
		AccessMode:   neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	log.Println("[NEO4J] Iniciando sincronizacion ")

	// 1. TRAER LOS CLIENTES
	rowsClientes, err := dbMySQL.QueryContext(ctx, "SELECT id_cliente, nombre FROM cliente")
	if err != nil {
		return err
	}
	defer rowsClientes.Close()

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		for rowsClientes.Next() {
			var idCliente int
			var nombreCliente string
			if err := rowsClientes.Scan(&idCliente, &nombreCliente); err != nil {
				continue
			}

			query := "MERGE (c:Comprador {id: $id}) SET c.nombre = $nombre"
			tx.Run(ctx, query, map[string]interface{}{
				"id":     idCliente,
				"nombre": nombreCliente,
			})
		}
		return nil, nil
	})
	if err != nil {
		log.Printf("[NEO4J] Error procesando clientes: %v", err)
	}

	// 2. TRAER LAS OBRAS
	collObra := dbMongo.Collection("obra_ultimate")
	cursorObras, err := collObra.Find(ctx, bson.M{})
	if err == nil {
		defer cursorObras.Close(ctx)
		_, _ = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			for cursorObras.Next(ctx) {
				var obra bson.M
				if err := cursorObras.Decode(&obra); err != nil {
					continue
				}

				// EXTRACCION SEGURA DEL OBJETO EMBEBIDO "artista"
				var nombreArtista string
				if artistaDoc, ok := obra["artista"].(bson.M); ok {
					nombreArtista = fmt.Sprintf("%v", artistaDoc["nombre"])
				} else if artistaDocNormal, ok := obra["artista"].(map[string]interface{}); ok {
					nombreArtista = fmt.Sprintf("%v", artistaDocNormal["nombre"])
				}

				// EXTRACCION SEGURA DEL OBJETO EMBEBIDO "genero"
				var nombreGenero string
				if generoDoc, ok := obra["genero"].(bson.M); ok {
					nombreGenero = fmt.Sprintf("%v", generoDoc["nombre"])
				} else if generoDocNormal, ok := obra["genero"].(map[string]interface{}); ok {
					nombreGenero = fmt.Sprintf("%v", generoDocNormal["nombre"])
				}

				// CORRECCIÓN CYPHER: Guardamos la propiedad status en el nodo Obra
				queryMapeoCatalogo := `
					MERGE (g:Género {id: $idGenero})
					SET g.nombre = $nombreGenero

					MERGE (a:Artista {id: $idArtista})
					SET a.nombre = $nombreArtista

					MERGE (o:Obra {id: $idObra})
					SET o.nombre = $nombreObra, o.foto = $fotoObra, o.status = $statusObra

					MERGE (a)-[:TRABAJA_EN]->(g)
					MERGE (a)-[:CREÓ]->(o)
				`

				tx.Run(ctx, queryMapeoCatalogo, map[string]interface{}{
					"idGenero":     obra["id_genero"],
					"nombreGenero": nombreGenero, 
					"idArtista":    obra["id_artista"],
					"nombreArtista": nombreArtista, 
					"idObra":       obra["_id"],
					"nombreObra":   obra["nombre"],
					"fotoObra":     obra["foto"],
					"statusObra":   obra["status"], // <-- CORRECCIÓN: Le pasamos el status de Mongo
				})
			}
			return nil, nil
		})
	}

	// 3. RELACIONES DE COMPRA
	rowsCompras, err := dbMySQL.QueryContext(ctx, "SELECT id_cliente, id_obra FROM orden")
	if err == nil {
		defer rowsCompras.Close()
		_, _ = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			for rowsCompras.Next() {
				var idCliente, idObra int
				if err := rowsCompras.Scan(&idCliente, &idObra); err != nil {
					continue
				}

				// El MATCH busca los nodos que ya creaste en PASO A y PASO B
				queryRelacion := `
					MATCH (c:Comprador {id: $idCliente})
					MATCH (o:Obra {id: $idObra})
					MERGE (c)-[:COMPRÓ]->(o)
				`
				tx.Run(ctx, queryRelacion, map[string]interface{}{
					"idCliente": idCliente,
					"idObra":    idObra,
				})
			}
			return nil, nil
		})
	}

	// 4. GENERAR RECOMENDACIONES DE GÉNERO
    log.Println("[NEO4J] Generando máximo 25 sugerencias por cliente...")
    _, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        querySugerencias := `
            // Buscamos clientes que compraron y sus géneros preferidos
            MATCH (c:Comprador)-[:COMPRÓ]->(obraComprada:Obra)<-[:CREÓ]-(:Artista)-[:TRABAJA_EN]->(g:Género)
            
            // Buscamos obras sugeridas del mismo género
            MATCH (g)<-[:TRABAJA_EN]-(:Artista)-[:CREÓ]->(obraSugerida:Obra)
            WHERE NOT (c)-[:COMPRÓ]->(obraSugerida)
            
            // Agrupamos por cliente para poder limitar el resultado
            WITH c, obraSugerida
            ORDER BY rand() // Mezcla las obras aleatoriamente
            WITH c, collect(obraSugerida)[0..25] AS sugerencias // Toma solo 25
            
            UNWIND sugerencias AS obra
            MERGE (c)-[:SUGERIDA_POR_GENERO]->(obra)
        `
        _, txErr := tx.Run(ctx, querySugerencias, nil)
        return nil, txErr
    })
	if err != nil {
		log.Printf("[NEO4J] Error al generar sugerencias: %v", err)
	}

	log.Println("[NEO4J] Sincronización finalizada con exito sin duplicados.")
	return nil
}
