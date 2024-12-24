package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// URL do KSQLDB
const ksqlDBURL = "http://localhost:8088"

// Definições do stream e tabela
const streamSQL = `
CREATE STREAM IF NOT EXISTS click_events_stream (
  itemId STRING,
  campaignId STRING,
  timestamp STRING
) WITH (
  KAFKA_TOPIC='click_events',
  VALUE_FORMAT='JSON'
);`

const tableSQL = `
CREATE TABLE IF NOT EXISTS click_counts_table AS
  SELECT
    campaignId,
    COUNT(*) AS click_count,
    WINDOWSTART AS window_start,
    WINDOWEND AS window_end
  FROM click_events_stream
  WINDOW TUMBLING (SIZE 2 MINUTES)
  GROUP BY campaignId
  EMIT CHANGES;`

// Função para criar o stream e a tabela no KSQLDB
func initializeKSQLObjects() {

	log.Println("Creating TABLE and STREAM if not exists")
	queries := []string{streamSQL, tableSQL}
	for _, query := range queries {
		//log.Println(fmt.Sprintf("Running on KSQL: %s", query))
		if err := executeKSQLQuery(query); err != nil {
			log.Fatalf("Failed to initialize KSQL objects: %v", err)
		}
		log.Printf("Successfully executed KSQL query: %s", query)
	}
}

// Executa uma query no KSQLDB
func executeKSQLQuery(query string) error {
	body := map[string]interface{}{
		"ksql":              query,
		"streamsProperties": map[string]interface{}{},
	}
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal KSQL query: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/ksql", ksqlDBURL), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send KSQL query: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("KSQL query failed: %s", string(responseBody))
	}

	return nil
}

// Estrutura das métricas agregadas
type CampaignMetrics struct {
	CampaignID  string `json:"campaignId"`
	ClickCount  int    `json:"click_count"`
	WindowStart string `json:"window_start"`
	WindowEnd   string `json:"window_end"`
}

// Consulta periódica ao KSQLDB
func queryKSQLDB() ([]CampaignMetrics, error) {
	query := `{
		"ksql": "SELECT campaignId, click_count, window_start, window_end FROM click_counts_table EMIT CHANGES LIMIT 10;"
	}`

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/query", ksqlDBURL), bytes.NewBuffer([]byte(query)))
	if err != nil {
		return nil, fmt.Errorf("failed to create KSQL request: %v", err)
	}

	req.Header.Set("Content-Type", "application/vnd.ksql.v1+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send KSQL request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("ksql query failed: %s", string(body))
	}

	// Decodifica a resposta como um array de objetos JSON
	var ksqlResponse []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ksqlResponse); err != nil {
		return nil, fmt.Errorf("failed to decode KSQL response: %v", err)
	}

	var metrics []CampaignMetrics

	for _, record := range ksqlResponse {
		// Processa apenas objetos do tipo "row"
		if row, ok := record["row"].(map[string]interface{}); ok {
			columns := row["columns"].([]interface{})
			metrics = append(metrics, CampaignMetrics{
				CampaignID:  columns[0].(string),
				ClickCount:  int(columns[1].(float64)),
				WindowStart: time.UnixMilli(int64(columns[2].(float64))).Format(time.RFC3339),
				WindowEnd:   time.UnixMilli(int64(columns[3].(float64))).Format(time.RFC3339),
			})
		}
	}

	return metrics, nil
}

func queryKSQLDBWithRetries() []CampaignMetrics {
	for {
		metrics, err := queryKSQLDB()
		if err != nil {
			log.Printf("Failed to query KSQLDB: %v", err)
			log.Println("Retrying in 10 seconds...")
			time.Sleep(10 * time.Second)
			continue
		}

		log.Println("Successfully connected to KSQLDB.")
		return metrics
	}
}

func main() {
	log.Println("Starting clickstream-processor")

	// Inicializa o stream e a tabela no KSQLDB
	initializeKSQLObjects()

	for {
		log.Println("Waiting for metrics")

		// Consulta o KSQLDB com resiliência
		metrics := queryKSQLDBWithRetries()

		// Processa as métricas obtidas
		for _, metric := range metrics {
			log.Printf(
				"Processed Metrics: CampaignID=%s, ClickCount=%d, WindowStart=%s, WindowEnd=%s",
				metric.CampaignID,
				metric.ClickCount,
				metric.WindowStart,
				metric.WindowEnd,
			)
			// Enviar para Grafana, ElasticSearch ou outros sistemas
		}

		// Aguardar até a próxima janela de tempo (2 minutos neste caso)
		time.Sleep(1 * time.Minute)
	}
}
