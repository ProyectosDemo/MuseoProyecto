package main

import (
	"log"
	"main/graph"
	"main/mysql"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	gqlparser "github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

func PrintFuncionesDisponibles(schemaAst *ast.Schema) {

	// Loguear las operaciones encontradas en el schema para depuración
	var mutationNames []string
	if m := schemaAst.Types["Mutation"]; m != nil {
		for _, f := range m.Fields {
			mutationNames = append(mutationNames, f.Name)
		}
	}
	var queryNames []string
	if q := schemaAst.Types["Query"]; q != nil {
		for _, f := range q.Fields {
			queryNames = append(queryNames, f.Name)
		}
	}
	log.Printf("Schema Mutation fields: %v", mutationNames)
	log.Printf("Schema Query fields: %v", queryNames)
}

func main() {
	// inicializar BD
	mysql.ConectarBD()
	if err := mysql.GetBD().Ping(); err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Cargar y parsear esquema GraphQL desde archivo para asegurar que incluya createCliente
	schemaBytes, err := os.ReadFile("graph/schema.graphqls")
	if err != nil {
		log.Fatalf("No se pudo leer schema.graphqls: %v", err)
	}
	schemaAst, err := gqlparser.LoadSchema(&ast.Source{Input: string(schemaBytes)})
	if err != nil {
		log.Fatalf("Error al parsear schema.graphqls: %v", err)
	}
	PrintFuncionesDisponibles(schemaAst)
	// crear servidor GraphQL usando el schema parseado
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: mysql.GetBD()}, Schema: schemaAst}))

	// Registrar transportes HTTP que aceptarás (POST, GET, OPTIONS)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	// Caché y extensiones recomendadas
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	// Crear router Gin y exponer endpoints
	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Error al configurar proxies: %v", err)
	}

	// Home de la web
	router.GET("/", func(c *gin.Context) {
		//playground.Handler("GraphQL playground", "/query").ServeHTTP(c.Writer, c.Request)
		c.File("./Frontend/html/index.html")
	})

	// las paginas nuevas las agregas aqui para servirlas
	paginas := []string{
		"index",
		"login",
		"exposiciones",
		"artistas",
		"obra",
		"admin",
		"admin-clientes",
		"admin-trabajadores",
		"admin-obras",
		"admin-membresias",
		"admin-artistas",
		"admin-generos",
		"admin-tarjetas",
		"artista-detalle",
		"consultas",
		"cuenta-cliente",
		"login-trabajador",
		"registro",
		"reservas"}

	for _, pagina := range paginas {
		router.GET("/"+pagina+".html", func(c *gin.Context) {
			c.File("./Frontend/html/" + pagina + ".html")
		})
	}

	router.Static("/static", "./Frontend")

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
