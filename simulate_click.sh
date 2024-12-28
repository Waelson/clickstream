#!/bin/bash

# URL base do endpoint
URL="http://localhost:8079/api/click"

# Declaração dos dados possíveis
declare -a requests=(
  '{"itemId": "1", "campaignId": "promocao-natal"}'
  '{"itemId": "2", "campaignId": "promocao-natal"}'
  '{"itemId": "3", "campaignId": "saldao-black-friday"}'
  '{"itemId": "4", "campaignId": "saldao-black-friday"}'
)

# Função para realizar chamadas aleatórias
send_random_request() {
  # Seleciona um índice aleatório
  random_index=$((RANDOM % ${#requests[@]}))

  # Recupera o corpo da requisição aleatória
  request_data=${requests[$random_index]}

  # Realiza a chamada usando curl
  echo "Enviando requisição: $request_data"
  curl --location --request GET "$URL" \
    --header 'Content-Type: application/json' \
    --data "$request_data"

  echo -e "\nRequisição enviada com sucesso."
}

# Loop indefinido
while true; do
  send_random_request

  # Pausa entre as requisições (ajustável, aqui 1 segundo)
  #sleep 1
done
