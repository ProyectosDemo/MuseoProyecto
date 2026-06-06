package main

import (
	"context"
	"log"
	"main/data_bases/cassandra"
	"main/data_bases/mongodb"
	"main/data_bases/mysql"
	"main/data_bases/neo4j"
	"main/graph"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	gqlparser "github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	ctx := context.Background() // Contexto necesario para el ciclo de vida del driver de Neo4j

	// Inicializar BDs relacionales y NoSQL operativas
	mysql.ConectarBD()
	if err := mysql.GetBD().Ping(); err != nil {
		log.Fatalf("Error al conectar a la base de datos MySQL: %v", err)
	}

	mongodb.ConectarMongo()
	
	cassandra.ConectarCassandra() 
	defer cassandra.GetCassandra().Close()

	// Inicializar y asegurar el cierre ordenado de Neo4j
	neo4j.ConectarNeo4j()
	defer neo4j.GetNeo4j().Close(ctx) 

	// Carga y parseo del esquema de GraphQL
	schemaBytes, err := os.ReadFile("graph/schema.graphqls")
	if err != nil {
		log.Fatalf("No se pudo leer schema.graphqls: %v", err)
	}
	schemaAst, err := gqlparser.LoadSchema(&ast.Source{Input: string(schemaBytes)})
	if err != nil {
		log.Fatalf("Error al parsear schema.graphqls: %v", err)
	}
	
	// Configuracion del servidor GraphQL inyectando todas las dbs
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			DB:        mysql.GetBD(),        // MySQL
			MongoDB:   mongodb.GetMongoDB(),   // MongoDB
			Cassandra: cassandra.GetCassandra(), // Cassandra
			Neo4j:     neo4j.GetNeo4j(),     // Neo4j
		},
		Schema: schemaAst,
	}))

	// Registrar transportes HTTP obligatorios para gqlgen
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	// Cache y extensiones del servidor de GraphQL
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	// Crear router Gin y exponer los recursos del Frontend
	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Error al configurar proxies: %v", err)
	}

	router.Static("/Frontend", "./Frontend")

	router.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL playground", "/query").ServeHTTP(c.Writer, c.Request)
		c.File("./Frontend/html/index.html")
	})

	// Middleware global para el manejo de CORS de llamadas fetch
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	router.Any("/query", gin.WrapH(srv))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor GraphQL corriendo en http://localhost:%s/", port)
	log.Fatal(router.Run(":" + port))
}