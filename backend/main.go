package main

import (
	"carbone-app/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
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
	var input struct {
		Category   string                 `json:"category"`
		UserInputs map[string]interface{} `json:"userInputs"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input data"})
		return
	}

	factors := models.CarbonFactors{}

	factors.Transports.Train = 0.014
	factors.Transports.Flight = 0.285
	factors.Transports.Car.Small = 0.1
	factors.Transports.Car.Medium = 0.2
	factors.Transports.Car.Big = 0.3

	factors.LogementElectromenagers.Electricity = 0.57
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

	factors.Alimentation.ShortCircuit.None = 1.0
	factors.Alimentation.ShortCircuit.Partial = 0.9
	factors.Alimentation.ShortCircuit.Majority = 0.8

	factors.Vetements.Large = 15
	factors.Vetements.Small = 10
	factors.Vetements.Madein.France = 1
	factors.Vetements.Madein.Autre = 1.2

	factors.Numerique.GoogleSearch = 0.0002
	factors.Numerique.ChatGPT = 0.000382
	factors.Numerique.SocialMedia = 0.000380
	factors.Numerique.Smartphone.Small = 35
	factors.Numerique.Smartphone.Large = 75
	factors.Numerique.Smartphone.Used = 0.5
	factors.Numerique.Smartphone.Old = 0.5

	factors.Consommation.Ecommerce.Amazon = 0.25
	factors.Consommation.Ecommerce.LeBonCoin = 0.05
	factors.Consommation.Ecommerce.Artisanat = 0.1
	factors.Consommation.Commerce.Brocante = 0.03
	factors.Consommation.Commerce.LocalShops = 0.08

	var result float64

	switch input.Category {
	case "Transports":
		if trainKm, ok := input.UserInputs["trainKm"].(float64); ok {
			result += trainKm * factors.Transports.Train
		}
		if flightKm, ok := input.UserInputs["flightKm"].(float64); ok {
			result += flightKm * factors.Transports.Flight
		}
		if carKm, ok := input.UserInputs["carKm"].(float64); ok {
			if carType, ok := input.UserInputs["carType"].(string); ok {
				if occupants, ok := input.UserInputs["carOccupants"].(float64); ok && occupants > 0 {
					switch carType {
					case "small":
						result += (carKm * factors.Transports.Car.Small) / occupants
					case "medium":
						result += (carKm * factors.Transports.Car.Medium) / occupants
					case "big":
						result += (carKm * factors.Transports.Car.Big) / occupants
					}
				}
			}
		}

	case "Logement_electromenagers":
		if electricityKwh, ok := input.UserInputs["electricityKwh"].(float64); ok {
			if homeOccupants, ok := input.UserInputs["homeOccupants"].(float64); ok && homeOccupants > 0 {
				result += (electricityKwh * factors.LogementElectromenagers.Electricity) / homeOccupants
			}
		}
		if gasKwh, ok := input.UserInputs["gasKwh"].(float64); ok {
			if homeOccupants, ok := input.UserInputs["homeOccupants"].(float64); ok && homeOccupants > 0 {
				result += (gasKwh * factors.LogementElectromenagers.Gas) / homeOccupants
			}
		}
		if housingType, ok := input.UserInputs["housingType"].(string); ok {
			if homeSize, ok := input.UserInputs["homeSize"].(float64); ok {
				if homeOccupants, ok := input.UserInputs["homeOccupants"].(float64); ok && homeOccupants > 0 {
					switch housingType {
					case "apartment":
						result += (factors.LogementElectromenagers.Apartment * homeSize) / homeOccupants
					case "house":
						result += (factors.LogementElectromenagers.House * homeSize) / homeOccupants
					}
				}
			}
		}
		if applianceCount, ok := input.UserInputs["applianceCount"].(float64); ok {
			if homeOccupants, ok := input.UserInputs["homeOccupants"].(float64); ok && homeOccupants > 0 {
				result += (applianceCount * factors.LogementElectromenagers.Appliance) / homeOccupants
			}
		}
		if electronicCount, ok := input.UserInputs["electronicCount"].(float64); ok {
			if homeOccupants, ok := input.UserInputs["homeOccupants"].(float64); ok && homeOccupants > 0 {
				result += (electronicCount * factors.LogementElectromenagers.Electronic) / homeOccupants
			}
		}

	case "Alimentation":
		if redMeatKg, ok := input.UserInputs["redMeatKg"].(float64); ok {
			result += redMeatKg * factors.Alimentation.RedMeat
		}
		if whiteMeatKg, ok := input.UserInputs["whiteMeatKg"].(float64); ok {
			result += whiteMeatKg * factors.Alimentation.WhiteMeat
		}
		if porkKg, ok := input.UserInputs["porkKg"].(float64); ok {
			result += porkKg * factors.Alimentation.Pork
		}
		if bulkPurchase, ok := input.UserInputs["bulkPurchase"].(string); ok {
			switch bulkPurchase {
			case "none":
				result *= factors.Alimentation.BulkFoodPurchase.None
			case "partial":
				result *= factors.Alimentation.BulkFoodPurchase.Partial
			case "total":
				result *= factors.Alimentation.BulkFoodPurchase.Total
			}
		}
		if shortCircuit, ok := input.UserInputs["shortCircuit"].(string); ok {
			switch shortCircuit {
			case "none":
				result *= factors.Alimentation.ShortCircuit.None
			case "partial":
				result *= factors.Alimentation.ShortCircuit.Partial
			case "majority":
				result *= factors.Alimentation.ShortCircuit.Majority
			}
		}

	case "Vetements":
		if largeItems, ok := input.UserInputs["largeItems"].(float64); ok {
			result += largeItems * factors.Vetements.Large
		}
		if smallItems, ok := input.UserInputs["smallItems"].(float64); ok {
			result += smallItems * factors.Vetements.Small
		}
		if origin, ok := input.UserInputs["origin"].(string); ok {
			switch origin {
			case "france":
				result *= factors.Vetements.Madein.France
			case "autre":
				result *= factors.Vetements.Madein.Autre
			}
		}

	case "Numerique":
		if googleSearches, ok := input.UserInputs["googleSearches"].(float64); ok {
			result += (googleSearches * 365) * factors.Numerique.GoogleSearch
		}
		if chatgptPrompts, ok := input.UserInputs["chatgptPrompts"].(float64); ok {
			result += (chatgptPrompts * 365) * factors.Numerique.ChatGPT
		}
		if smartphoneType, ok := input.UserInputs["smartphoneType"].(string); ok && smartphoneType != "" {
			var baseEmission float64
			switch smartphoneType {
			case "small":
				baseEmission = factors.Numerique.Smartphone.Small
			case "large":
				baseEmission = factors.Numerique.Smartphone.Large
			}

			if state, ok := input.UserInputs["smartphoneState"].(string); ok {
				switch state {
				case "used":
					baseEmission *= factors.Numerique.Smartphone.Used
				case "old":
					baseEmission *= factors.Numerique.Smartphone.Old
				}
			}

			result += baseEmission
		}
		if socialHours, ok := input.UserInputs["socialHours"].(float64); ok {
			result += (socialHours * 365) * factors.Numerique.SocialMedia
		}

	case "Consommation":
		if amazonOrders, ok := input.UserInputs["amazonOrders"].(float64); ok {
			result += amazonOrders * factors.Consommation.Ecommerce.Amazon
		}
		if leboncoinOrders, ok := input.UserInputs["leboncoinOrders"].(float64); ok {
			result += leboncoinOrders * factors.Consommation.Ecommerce.LeBonCoin
		}
		if artisanatOrders, ok := input.UserInputs["artisanatOrders"].(float64); ok {
			result += artisanatOrders * factors.Consommation.Ecommerce.Artisanat
		}
		if brocanteItems, ok := input.UserInputs["brocanteItems"].(float64); ok {
			result += brocanteItems * factors.Consommation.Commerce.Brocante
		}
		if localShopOrders, ok := input.UserInputs["localShopOrders"].(float64); ok {
			result += localShopOrders * factors.Consommation.Commerce.LocalShops
		}
	}

	c.JSON(200, gin.H{
		"category": input.Category,
		"result":   result,
	})
}

func getCarbonFactors(c *gin.Context) {
	factors := models.CarbonFactors{}

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

	factors.Alimentation.ShortCircuit.None = 1.0
	factors.Alimentation.ShortCircuit.Partial = 0.9
	factors.Alimentation.ShortCircuit.Majority = 0.8

	factors.Vetements.Large = 15
	factors.Vetements.Small = 10
	factors.Vetements.Madein.France = 1
	factors.Vetements.Madein.Autre = 1.2

	factors.Numerique.GoogleSearch = 0.0002
	factors.Numerique.ChatGPT = 0.000382
	factors.Numerique.SocialMedia = 0.000380
	factors.Numerique.Smartphone.Small = 35
	factors.Numerique.Smartphone.Large = 75
	factors.Numerique.Smartphone.Used = 0.5
	factors.Numerique.Smartphone.Old = 0.5

	factors.Consommation.Ecommerce.Amazon = 0.25
	factors.Consommation.Ecommerce.LeBonCoin = 0.05
	factors.Consommation.Ecommerce.Artisanat = 0.1
	factors.Consommation.Commerce.Brocante = 0.03
	factors.Consommation.Commerce.LocalShops = 0.08

	c.JSON(200, factors)
}
