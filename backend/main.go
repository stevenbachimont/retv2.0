package main

import (
	"carbone-app/models"

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
	factors := models.CarbonFactors{}

	// Initialisation des données
	factors.Transports.Train = 0.014
	factors.Transports.Flight = 0.285
	factors.Transports.Car.Small = 0.1
	factors.Transports.Car.Medium = 0.2
	factors.Transports.Car.Big = 0.3

	factors.LogementElectromenagers.Electricity = 0.4
	factors.LogementElectromenagers.Gas = 0.2
	factors.LogementElectromenagers.Apartment = 15
	factors.LogementElectromenagers.House = 20
	factors.LogementElectromenagers.Appliance = 0.5
	factors.LogementElectromenagers.Electronic = 0.3

	factors.Alimentation.RedMeat = 27
	factors.Alimentation.WhiteMeat = 6.9
	factors.Alimentation.Pork = 7.2
	factors.Alimentation.BulkFoodPurchase.None = 1
	factors.Alimentation.BulkFoodPurchase.Partial = 0.9
	factors.Alimentation.BulkFoodPurchase.Total = 0.8

	factors.Vetements.Large = 15
	factors.Vetements.Small = 10
	factors.Vetements.Madein.France = 1
	factors.Vetements.Madein.Autre = 1.2

	c.JSON(200, factors)
}
