package main

import (
	"carbone-app/config"
	"carbone-app/models"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("votre_clé_secrète_ici") // En production, utilisez une vraie clé secrète

type Claims struct {
	UserID string
	jwt.StandardClaims
}

func dbMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()
	r.Use(dbMiddleware(db))

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	// Routes pour les calculs carbone
	api := r.Group("/api")
	{
		// Routes publiques
		api.GET("/factors", getCarbonFactors)
		api.POST("/register", register)
		api.POST("/login", login)
		api.GET("/verify", verifyToken)

		// Routes protégées
		authorized := api.Group("")
		authorized.Use(authMiddleware())
		{
			authorized.POST("/calculate", calculateCarbon)
			authorized.POST("/results", saveResult)
			authorized.GET("/results", getResults)
		}
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

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization required"})
			c.Abort()
			return
		}
		// Vérifier le token JWT et ajouter l'userID au contexte
		userID, err := validateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}

func register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		log.Printf("Register - Erreur de binding: %v", err)
		c.JSON(400, gin.H{"error": "Email et mot de passe requis"})
		return
	}

	log.Printf("Register - Tentative d'inscription pour: %s", input.Email)
	log.Printf("Register - Mot de passe reçu (len=%d): %s", len(input.Password), input.Password)

	// Vérifier si l'email existe déjà
	db := c.MustGet("db").(*sql.DB)
	var existingID string
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", input.Email).Scan(&existingID)
	if err == nil {
		log.Printf("Register - Email déjà utilisé: %s", input.Email)
		c.JSON(400, gin.H{"error": "Cet email est déjà utilisé"})
		return
	}

	// Créer le nouvel utilisateur
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Register - Erreur de hachage: %v", err)
		c.JSON(500, gin.H{"error": "Erreur serveur"})
		return
	}

	userID := uuid.New().String()
	_, err = db.Exec(`
		INSERT INTO users (id, email, password)
		VALUES ($1, $2, $3)
	`, userID, input.Email, string(hashedPassword))

	if err != nil {
		log.Printf("Register - Erreur d'insertion: %v", err)
		c.JSON(500, gin.H{"error": "Impossible de créer l'utilisateur"})
		return
	}

	log.Printf("Register - Inscription réussie pour: %s (ID: %s)", input.Email, userID)

	token := generateToken(userID)
	c.JSON(200, gin.H{
		"token": token,
		"user": gin.H{
			"id":    userID,
			"email": input.Email,
		},
	})
}

func login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&credentials); err != nil {
		log.Printf("Login - Erreur de binding: %v", err)
		c.JSON(400, gin.H{"error": "Données invalides"})
		return
	}

	log.Printf("Login - Tentative de connexion pour: %s avec mot de passe: %s", credentials.Email, credentials.Password)

	db := c.MustGet("db").(*sql.DB)
	var user models.User
	var hashedPassword string

	// Ajout de logs pour déboguer la requête SQL
	query := "SELECT id, email, password FROM users WHERE email = $1"
	err := db.QueryRow(query, credentials.Email).Scan(&user.ID, &user.Email, &hashedPassword)
	if err != nil {
		log.Printf("Login - Utilisateur non trouvé: %v", err)
		c.JSON(401, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	log.Printf("Login - Hash stocké: %s", hashedPassword)
	log.Printf("Login - Mot de passe fourni: %s", credentials.Password)

	// Ajout de logs pour déboguer la comparaison des mots de passe
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password)); err != nil {
		log.Printf("Login - Erreur de comparaison: %v", err)
		log.Printf("Login - Hash stocké (len=%d): %x", len(hashedPassword), hashedPassword)
		log.Printf("Login - Mot de passe fourni (len=%d): %x", len(credentials.Password), credentials.Password)
		c.JSON(401, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	log.Printf("Login - Connexion réussie pour: %s (ID: %s)", credentials.Email, user.ID)
	token := generateToken(user.ID)
	c.JSON(200, gin.H{
		"token": token,
		"user":  user,
	})
}

func saveResult(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	userID, _ := c.Get("userID")
	var result models.Result
	if err := c.BindJSON(&result); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	result.ID = uuid.New().String()
	result.UserID = userID.(string)
	result.CreatedAt = time.Now()

	_, err := db.Exec(`
		INSERT INTO results (id, user_id, category, value, inputs, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, result.ID, result.UserID, result.Category, result.Value, result.Inputs, result.CreatedAt)

	if err != nil {
		c.JSON(500, gin.H{"error": "Could not save result"})
		return
	}

	c.JSON(200, result)
}

func getResults(c *gin.Context) {
	_, _ = c.Get("userID")
	// TODO: Récupérer les résultats depuis la base de données
	c.JSON(200, []models.Result{})
}

func validateToken(tokenStr string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		log.Printf("Erreur parsing token: %v", err)
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token invalide")
	}

	return claims.UserID, nil
}

func generateToken(userID string) string {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Erreur génération token: %v", err)
		return ""
	}

	return tokenString
}

func verifyToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{"error": "No token provided"})
		return
	}

	// Enlever le préfixe "Bearer " si présent
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	userID, err := validateToken(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	// Récupérer l'utilisateur depuis la base de données
	db := c.MustGet("db").(*sql.DB)
	var user models.User
	err = db.QueryRow("SELECT id, email FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Email)
	if err != nil {
		log.Printf("Erreur de recherche utilisateur: %v", err)
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, user)
}
