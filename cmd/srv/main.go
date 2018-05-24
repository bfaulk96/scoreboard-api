package main

import (
	"github.com/rs/cors"
	"log"
	"github.com/bfaulk96/scoreboard-api/pkg/configurations"
	"github.com/bfaulk96/scoreboard-api/pkg/database"
	"github.com/bfaulk96/scoreboard-api/pkg/router"
	"net/http"
	"fmt"
)

func main() {
	// Read config file and set up logging
	config := configurations.LoadConfig()

	collection := database.InitializeMongoDatabase(config)

	// Create API routes
	r := router.New()
	r.CreateRoutes(collection, config)

	// Start Web Server on given port with CORS enabled
	fmt.Printf("Server listening on port %s\n", config.Port)
	log.Fatal(http.ListenAndServe(":" + config.Port, cors.New(
		cors.Options{
			AllowedOrigins: []string{config.AllowedOrigins},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
			AllowedMethods: []string{"GET","POST","PUT","DELETE", "OPTIONS"},
			AllowCredentials: true,
		},
	).Handler(r)))
}
