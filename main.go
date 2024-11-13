package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/joho/godotenv"
	"github.com/mike-jacks/neo/db"
	"github.com/mike-jacks/neo/generated"
	"github.com/mike-jacks/neo/resolver"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

// LRUQueryCache is a custom cache that implements graphql.Cache[*ast.QueryDocument]
type LRUQueryCache struct {
	cache *lru.Cache[string, *ast.QueryDocument]
}

// NewLRUQueryCache creates a new LRUQueryCache with the given size
func NewLRUQueryCache(size int) (*LRUQueryCache, error) {
	cache, err := lru.New[string, *ast.QueryDocument](size)
	if err != nil {
		return nil, err
	}
	return &LRUQueryCache{cache: cache}, nil
}

// Add adds a query document to the cache
func (c *LRUQueryCache) Add(ctx context.Context, key string, value *ast.QueryDocument) {
	c.cache.Add(key, value)
}

// Get retrieves a query document from the cache
func (c *LRUQueryCache) Get(ctx context.Context, key string) (*ast.QueryDocument, bool) {
	return c.cache.Get(key)
}

// LRUStringCache is a custom cache that implements graphql.Cache[string]
type LRUStringCache struct {
	cache *lru.Cache[string, string]
}

// NewLRUStringCache creates a new LRUStringCache with the given size
func NewLRUStringCache(size int) (*LRUStringCache, error) {
	cache, err := lru.New[string, string](size)
	if err != nil {
		return nil, err
	}
	return &LRUStringCache{cache: cache}, nil
}

// Add adds a string to the cache
func (c *LRUStringCache) Add(ctx context.Context, key string, value string) {
	c.cache.Add(key, value)
}

// Get retrieves a string from the cache
func (c *LRUStringCache) Get(ctx context.Context, key string) (string, bool) {
	return c.cache.Get(key)
}

func setupGraphQLServer(db db.Database) *handler.Server {
	resolver := resolver.NewResolver(db)
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: resolver})
	server := handler.New(schema)

	server.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	server.AddTransport(transport.Options{})
	server.AddTransport(transport.GET{})
	server.AddTransport(transport.POST{})
	server.AddTransport(transport.MultipartForm{})

	// Create a custom LRU cache for query documents
	queryCache, err := NewLRUQueryCache(1000)
	if err != nil {
		log.Fatal("Error creating query cache:", err)
	}
	server.SetQueryCache(queryCache)

	// Create a custom LRU cache for persisted queries
	persistedQueryCache, err := NewLRUStringCache(100)
	if err != nil {
		log.Fatal("Error creating persisted query cache:", err)
	}

	server.Use(extension.Introspection{})
	server.Use(extension.AutomaticPersistedQuery{
		Cache: persistedQueryCache, // Use the custom LRUStringCache
	})

	return server
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	driver, err := db.SetupNeo4jDriver()
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close(context.Background())

	neo4jdb := &db.Neo4jDatabase{Driver: driver}

	srv := setupGraphQLServer(neo4jdb)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	loggingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "" {
			proto = "http"
		}
		log.Printf("Request received: Method: %s, Path: %s, Protocol: %s", r.Method, r.URL.Path, proto)

		if websocket.IsWebSocketUpgrade(r) {
			log.Printf("WebSocket Upgrade Detected. Origin: %s", r.Header.Get("Origin"))
			srv.ServeHTTP(w, r)
			return
		} else {
			log.Printf("Non-WebSocket Request Detected. Origin: %s", r.Header.Get("Origin"))
			corsHandler.Handler(srv).ServeHTTP(w, r)
		}
	})

	http.Handle("/graphql", corsHandler.Handler(playground.Handler("GraphQL Playground", "/query")))
	http.Handle("/query", loggingHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	url := os.Getenv("URL")

	var websocketUrl string
	if url == "" {
		url = "http://localhost" + ":" + port
		websocketUrl = "ws://localhost" + ":" + port
	} else {
		websocketUrl = strings.Replace(url, "https://", "wss://", 1)
	}

	log.Printf("Connect to %s/graphql for GraphQL Playground", url)
	log.Printf("GraphQL WebSocket endpoint: %s/query", websocketUrl)
	log.Printf("Connect to %s/query for GraphQL API", url)
	log.Printf("Connect to https://console.neo4j.io for Neo4j Browser Console")
	
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
