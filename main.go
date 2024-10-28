package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type formData struct {
	Name      string `json:"name"`
	Vacancy   string `json:"vacancy"`
	Level     string `json:"level"`
	HardSkill string `json:"hardskill"`
	SoftSkill string `json:"softskill"`
	Pcd       string `json:"pcd"`
	Former    string `json:"former"`
	Comment   string `json:"comment"`
}

func getPostData(w http.ResponseWriter, r *http.Request) {
	var answer formData

	err := json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	//fmt.Printf("Name: %s", answer.Name)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success!"))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Render!")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/answer", getPostData).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
