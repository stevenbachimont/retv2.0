package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Configuration CORS plus permissive
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Permet toutes les origines
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	// Routes pour les calculs carbone
	api := r.Group("/api")
	{
		api.POST("/calculate", calculateCarbon)
		api.GET("/factors", getCarbonFactors)
	}

	r.Run(":8080")
}

func calculateCarbon(c *gin.Context) {
	// TODO: Implémenter le calcul
	c.JSON(200, gin.H{
		"message": "Calcul en cours de développement",
	})
}

func getCarbonFactors(c *gin.Context) {
	// Facteurs d'émission de base
	factors := []gin.H{
		{"id": "voiture", "name": "Voiture", "factor": 0.2},
		{"id": "train", "name": "Train", "factor": 0.02},
		{"id": "avion", "name": "Avion", "factor": 0.285},
		{"id": "viande_rouge", "name": "Viande Rouge", "factor": 60.0},
	}

	c.JSON(200, factors)
}
