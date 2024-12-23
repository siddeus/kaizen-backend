package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func gemini(question string) (string, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	geminiUrl := os.Getenv("GEMINI_API_URL")

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": question,
					},
				},
			},
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Erro ao marshalling JSON: %v", err)
	}

	url := geminiUrl + geminiKey

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erro ao ler resposta: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		log.Fatalf("Erro ao unmarshalling JSON: %v", err)
	}

	if candidates, ok := result["candidates"].([]interface{}); ok && len(candidates) > 0 {
		if candidate, ok := candidates[0].(map[string]interface{}); ok {
			if content, ok := candidate["content"].(map[string]interface{}); ok {
				if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
					if part, ok := parts[0].(map[string]interface{}); ok {
						if text, ok := part["text"].(string); ok {
							return text, nil
						} else {
							log.Println("O campo 'text' não foi encontrado dentro de parts")
						}
					}
				} else {
					log.Println("O campo 'parts' não foi encontrado ou está vazio")
				}
			}
		}
	}

	return "NC", fmt.Errorf("o campo 'text' não foi encontrado na resposta")
}
