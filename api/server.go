package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vapor05/financeview/graph"
	"github.com/vapor05/financeview/graph/generated"
	"github.com/vapor05/financeview/pkg/store"
)

// Defining the Graphql handler
func graphqlHandler() (gin.HandlerFunc, error) {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	db, err := store.NewDatabase(context.Background(), os.Getenv("DBURL"))
	if err != nil {
		return func(c *gin.Context) {}, fmt.Errorf("failed to connect to database, %w", err)
	}
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Db: db}}))
	ghf := func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
	return ghf, nil
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func LogRequest(c *gin.Context) {
	b, err := io.ReadAll(ioutil.NopCloser(c.Request.Body))
	if err != nil {
		log.Printf("error reading request body, %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	log.Printf("Request: %s", b)
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(b))
	c.Next()
}

func main() {
	// Setting up Gin
	r := gin.Default()
	r.Use(LogRequest)
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"OPTIONS", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))
	gh, err := graphqlHandler()
	if err != nil {
		log.Fatalf("failed to setup graphql handler, %v", err)
	}
	r.POST("/query", gh)
	r.GET("/", playgroundHandler())
	r.Run()
}
