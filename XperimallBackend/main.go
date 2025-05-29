package main

import (
	"fmt"
	"os"

	"XperimallBackend/database"
	"XperimallBackend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	database.ConnectDB()

	r.Use(cors.Default())

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🟢 Server starting on port: " + port)
	err := r.Run("0.0.0.0:" + port)
	if err != nil {
		fmt.Println("🔴 Failed to start server:", err)
	}
}
