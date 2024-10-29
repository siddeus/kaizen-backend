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
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	// Cria uma nova instância do Gin
	router := gin.Default()

	// Define o endpoint POST
	router.POST("/submit", func(c *gin.Context) {
		var jsonData Data

		// Valida e vincula o JSON recebido à estrutura Data
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Retorna uma resposta de sucesso com os dados recebidos
		c.JSON(http.StatusOK, gin.H{
			"message": "Dados recebidos com sucesso!",
			"data":    jsonData,
		})
	})

	// Inicia o servidor na porta 8080
	router.Run(":" + port)
}
