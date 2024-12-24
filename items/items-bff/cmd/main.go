package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Item struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CampaignID string `json:"campaignId"`
	ImageURL   string `json:"imageUrl"`
}

var items = []Item{
	{ID: "1", Name: "Item A", CampaignID: "promocao-natal", ImageURL: "https://via.placeholder.com/200?text=Item+A"},
	{ID: "2", Name: "Item B", CampaignID: "promocao-natal", ImageURL: "https://via.placeholder.com/200?text=Item+B"},
	{ID: "3", Name: "Item C", CampaignID: "promocao-natal", ImageURL: "https://via.placeholder.com/200?text=Item+C"},
	{ID: "4", Name: "Item D", CampaignID: "saldao-black-friday", ImageURL: "https://via.placeholder.com/200?text=Item+D"},
	{ID: "5", Name: "Item E", CampaignID: "saldao-black-friday", ImageURL: "https://via.placeholder.com/200?text=Item+E"},
	{ID: "6", Name: "Item F", CampaignID: "saldao-black-friday", ImageURL: "https://via.placeholder.com/200?text=Item+F"},
}

// Middleware para habilitar CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Permitir qualquer origem
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Responder requisições OPTIONS diretamente
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

const clickstreamAPI = "http://localhost:4000/api/click"

func main() {
	// Criação de um roteador
	mux := http.NewServeMux()

	// Rota para retornar itens
	mux.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})

	// Rota para registrar cliques
	mux.HandleFunc("/api/click", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Evento de clique recebido")
		var click map[string]string
		if err := json.NewDecoder(r.Body).Decode(&click); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Encaminha o evento para o clickstream-api
		if err := forwardToClickstreamAPI(click); err != nil {
			log.Printf("Failed to forward click event: %v", err)
			http.Error(w, "Failed to register click event", http.StatusInternalServerError)
			return
		} else {
			log.Printf("Click registered and forwarded: %v", click)
			w.WriteHeader(http.StatusNoContent)
		}
	})

	// Adiciona o middleware de CORS ao roteador
	handler := enableCORS(mux)

	// Inicia o servidor
	port := ":3000"
	log.Printf("items-bff running on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, handler))
}

// Função para encaminhar o evento ao clickstream-api
func forwardToClickstreamAPI(click map[string]string) error {
	// Serializa o evento como JSON
	data, err := json.Marshal(click)
	if err != nil {
		return err
	}

	// Faz a requisição POST ao clickstream-api
	resp, err := http.Post(clickstreamAPI, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verifica se a resposta foi bem-sucedida
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
