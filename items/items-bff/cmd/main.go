package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Item struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CampaignID string `json:"campaignId"`
	ImageURL   string `json:"imageUrl"` // Novo campo
}

var items = []Item{
	{ID: "1", Name: "Item A", CampaignID: "promo1", ImageURL: "https://via.placeholder.com/200?text=Item+A"},
	{ID: "2", Name: "Item B", CampaignID: "promo1", ImageURL: "https://via.placeholder.com/200?text=Item+B"},
	{ID: "3", Name: "Item C", CampaignID: "promo2", ImageURL: "https://via.placeholder.com/200?text=Item+C"},
	{ID: "4", Name: "Item D", CampaignID: "promo2", ImageURL: "https://via.placeholder.com/200?text=Item+D"},
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

func main() {
	// Criação de um roteador
	mux := http.NewServeMux()

	// Rotas principais
	mux.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})

	mux.HandleFunc("/api/click", func(w http.ResponseWriter, r *http.Request) {
		var click map[string]string
		if err := json.NewDecoder(r.Body).Decode(&click); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		log.Printf("Click registered: %v", click)
		w.WriteHeader(http.StatusNoContent)
	})

	// Adiciona o middleware de CORS ao roteador
	handler := enableCORS(mux)

	// Inicia o servidor
	port := ":3000"
	log.Printf("items-bff running on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
