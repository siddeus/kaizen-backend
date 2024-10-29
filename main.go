package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Estrutura para receber os dados da requisição POST
type Data struct {
	Name       string `json:"name" binding:"required"`
	Vacancy    string `json:"vacancy" binding:"required"`
	Hardskills string `json:"hardskills" binding:"required"`
	Softskills string `json:"softskills" binding:"required"`
	Pcd        string `json:"pcd"`
	Former     string `json:"former"`
	Comment    string `json:"comment"`
}

func main() {
	//load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//get port variable
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	//Create new instance Gin
	router := gin.Default()

	//Set endpoint POST
	router.POST("/submit", func(c *gin.Context) {
		var jsonData Data

		//JSON data
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Return after response
		c.JSON(http.StatusOK, gin.H{
			"message": "Dados recebidos com sucesso!",
			"data":    jsonData,
		})
	})

	//Start server
	router.Run(":" + port)
}
