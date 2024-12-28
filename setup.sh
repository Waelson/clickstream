#!/bin/bash

# Nome do tópico Kafka
TOPIC="click_events"
PARTITIONS=3
REPLICATION_FACTOR=1

# Configuração do Kafka e KSQLDB
KAFKA_CONTAINER="clickstream-kafka-1"
KSQL_CONTAINER="clickstream-ksql-1"
KAFKA_SERVER="kafka:9092"

# Configuração do Kafka Connect
CONNECTOR_NAME="postgres-sink-connector"
CONNECTOR_CONFIG=$(cat <<EOF
{
  "name": "$CONNECTOR_NAME",
  "config": {
    "connector.class": "io.confluent.connect.jdbc.JdbcSinkConnector",
    "connection.url": "jdbc:postgresql://postgres:5432/db_metrics",
    "connection.user": "user_metrics",
    "connection.password": "password_metrics",
    "topics": "CLICK_COUNTS_TABLE_OUTPUT",
    "insert.mode": "insert",
    "auto.create": "true",
    "auto.evolve": "true",
    "tasks.max": "1",
    "key.converter": "org.apache.kafka.connect.storage.StringConverter",
    "value.converter": "io.confluent.connect.avro.AvroConverter",
    "value.converter.schema.registry.url": "http://schema-registry:8081",
    "transforms": "ConvertStartTime,ConvertEndTime",
    "transforms.ConvertStartTime.type": "org.apache.kafka.connect.transforms.TimestampConverter\$Value",
    "transforms.ConvertStartTime.field": "START_TIME",
    "transforms.ConvertStartTime.target.type": "Timestamp",
    "transforms.ConvertEndTime.type": "org.apache.kafka.connect.transforms.TimestampConverter\$Value",
    "transforms.ConvertEndTime.field": "END_TIME",
    "transforms.ConvertEndTime.target.type": "Timestamp"
  }
}
EOF
)

# Função para solicitar aprovação do usuário
prompt_user() {
  local message=$1
  echo "$message (y/n):"
  read -r response
  if [[ "$response" =~ ^[Yy]$ ]]; then
    return 0
  else
    echo "Ação cancelada pelo usuário."
    return 1
  fi
}

# Comando para criar tópico
create_topic() {
  if prompt_user "Deseja criar o tópico Kafka '$TOPIC'?"; then
    echo "Criando tópico: $TOPIC"
    docker exec "$KAFKA_CONTAINER" kafka-topics --create \
      --topic "$TOPIC" \
      --bootstrap-server "$KAFKA_SERVER" \
      --partitions "$PARTITIONS" \
      --replication-factor "$REPLICATION_FACTOR"
    if [ $? -eq 0 ]; then
      echo "Tópico '$TOPIC' criado com sucesso."
    else
      echo "Erro ao criar o tópico '$TOPIC'. Ele pode já existir."
    fi
  fi
}

# Comando para criar stream no KSQLDB
create_stream() {
  if prompt_user "Deseja criar o stream 'click_events_stream'?"; then
    echo "Criando stream no KSQLDB..."
    docker exec "$KSQL_CONTAINER" ksql -e "CREATE STREAM IF NOT EXISTS click_events_stream (itemId STRING, campaignId STRING, timestamp STRING) WITH (KAFKA_TOPIC='$TOPIC',VALUE_FORMAT='JSON');"

    if [ $? -eq 0 ]; then
      echo "Stream 'click_events_stream' criado com sucesso."
    else
      echo "Erro ao criar o stream 'click_events_stream'."
    fi
  fi
}

# Comando para criar tabela de agregação no KSQLDB
create_aggregation_table() {
  if prompt_user "Deseja criar a tabela de agregação 'CLICK_COUNTS_TABLE'?"; then
    echo "Criando tabela de agregação no KSQLDB..."
    docker exec "$KSQL_CONTAINER" ksql -e "CREATE TABLE IF NOT EXISTS CLICK_COUNTS_TABLE WITH (KAFKA_TOPIC='CLICK_COUNTS_TABLE_OUTPUT', VALUE_FORMAT='AVRO', PARTITIONS=3, REPLICAS=1) AS SELECT CLICK_EVENTS_STREAM.CAMPAIGNID AS CAMPAIGNID, AS_VALUE(CLICK_EVENTS_STREAM.CAMPAIGNID) AS ID, COUNT(*) AS CLICK_COUNT, WINDOWSTART AS START_TIME, WINDOWEND AS END_TIME FROM CLICK_EVENTS_STREAM CLICK_EVENTS_STREAM WINDOW TUMBLING (SIZE 2 MINUTES) GROUP BY CLICK_EVENTS_STREAM.CAMPAIGNID EMIT CHANGES;"

    if [ $? -eq 0 ]; then
      echo "Tabela 'CLICK_COUNTS_TABLE' criada com sucesso."
    else
      echo "Erro ao criar a tabela 'CLICK_COUNTS_TABLE'."
    fi
  fi
}

# Comando para criar conector no Kafka Connect
create_connector() {
  if prompt_user "Deseja configurar o conector Kafka Connect '$CONNECTOR_NAME'?"; then
    echo "Criando conector Kafka Connect..."
    curl --location --silent --show-error --fail \
      --request POST http://localhost:8083/connectors \
      --header "Content-Type: application/json" \
      --data "$CONNECTOR_CONFIG"
    if [ $? -eq 0 ]; then
      echo "Conector '$CONNECTOR_NAME' criado com sucesso."
    else
      echo "Erro ao criar o conector '$CONNECTOR_NAME'."
    fi
  fi
}

# Execução iterativa dos comandos
echo "Iniciando a configuração inicial da aplicação..."
create_topic
create_stream
create_aggregation_table
create_connector
echo "Configuração inicial concluída."
