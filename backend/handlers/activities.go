package handlers

import (
	"carbone-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetActivities(c *gin.Context) {
	var activities []models.Activity
	// TODO: Implémenter la logique de récupération depuis la base de données

	c.JSON(http.StatusOK, activities)
}

func CreateActivity(c *gin.Context) {
	var activity models.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implémenter la logique de sauvegarde en base de données

	c.JSON(http.StatusCreated, activity)
}

func GetCarbonTotal(c *gin.Context) {
	// TODO: Implémenter le calcul du total carbone
	total := 0.0

	c.JSON(http.StatusOK, gin.H{"total": total})
}
