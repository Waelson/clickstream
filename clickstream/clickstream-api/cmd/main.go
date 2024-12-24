package main

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
)

const (
	brokerAddress = "localhost:9092" // Endereço do Kafka Broker
	topic         = "click_events"   // Tópico para eventos de clique
)

// Estrutura do evento de clique
type ClickEvent struct {
	ItemID     string `json:"itemId"`
	CampaignID string `json:"campaignId"`
	Timestamp  string `json:"timestamp"`
}

// Função para registrar eventos de clique
func handleClickEvent(w http.ResponseWriter, r *http.Request) {
	var event ClickEvent

	// Decodifica o corpo da requisição JSON
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Publica o evento no Kafka
	if err := publishToKafka(event); err != nil {
		log.Printf("Failed to publish event: %v", err)
		http.Error(w, "Failed to register event", http.StatusInternalServerError)
		return
	}

	log.Printf("Event registered: %+v", event)
	w.WriteHeader(http.StatusNoContent)
}

// Funcao para publicar o evento no Kafka
func publishToKafka(event ClickEvent) error {
	// Cria um writer para o Kafka
	writer := kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}
	defer writer.Close()

	// Serializa o evento como JSON
	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Envia a mensagem para o Kafka
	return writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(event.ItemID), // Usando ItemID como chave
		Value: message,
	})
}

func main() {
	http.HandleFunc("/api/click", handleClickEvent)

	port := ":4000"
	log.Printf("clickstream-api running on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
