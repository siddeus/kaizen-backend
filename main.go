package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.Default()

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
