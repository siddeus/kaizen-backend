package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
)

type Data struct {
	Name       string `json:"name" binding:"required"`
	Vacancy    string `json:"vacancy" binding:"required"`
	Hardskills string `json:"hardskills" binding:"required"`
	Softskills string `json:"softskills" binding:"required"`
	Pcd        string `json:"pcd"`
	Former     string `json:"former"`
	Comment    string `json:"comment"`
}

type Question struct {
	Message string `json:"message"`
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	// Configuração do CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                             // Permite todas as origens (substitua "*" por uma origem específica, como "http://localhost:3000")
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},        // Métodos permitidos
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // Cabeçalhos permitidos
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // Permitir cookies (se necessário)
		MaxAge:           12 * time.Hour, // Cache de CORS por 12 horas
	}))

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router.POST("/getNewQuestion", func(c *gin.Context) {

		var jsonQuestion Question
		var questionBase = "Gere uma questão curta para entrevista sobre Magento 2. O nível da pergunta é " + jsonQuestion.Message

		if err := c.ShouldBindJSON(&jsonQuestion); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		answerGemini, err := gemini(questionBase)
		if err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"question": answerGemini,
		})

		jsonQuestionData, err := json.Marshal(jsonQuestion)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println(jsonQuestionData)
	})

	router.POST("/submit", func(c *gin.Context) {

		var jsonData Data

		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Dados recebidos com sucesso!",
			"data":    jsonData,
		})

		jsonDatas, err := json.Marshal(jsonData)
		if err != nil {
			log.Fatal(err)
		}

		jsonString := string(jsonDatas)

		saveData(jsonString)
	})

	router.Run(":" + port)

}

func saveData(jsonString string) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongodbUser := os.Getenv("MONGODB_USER")
	mongodbPass := os.Getenv("MONGODB_PASSWORD")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://" + mongodbUser + ":" + mongodbPass + "@cluster0.eenfx.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("rh").Collection("answers")

	var doc Data

	err = json.Unmarshal([]byte(jsonString), &doc)
	if err != nil {
		log.Fatal(err)
	}

	insertResult, err := collection.InsertOne(ctx, bson.M{
		"name":       doc.Name,
		"vacancy":    doc.Vacancy,
		"hardskills": doc.Hardskills,
		"softskills": doc.Softskills,
		"pcd":        doc.Pcd,
		"former":     doc.Former,
		"comment":    doc.Comment,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Documento inserido com ID: %v\n", insertResult.InsertedID)
}
