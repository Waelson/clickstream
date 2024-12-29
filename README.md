# Solução de Clickstream para E-commerce

## Introdução

Esta é uma solução de **Clickstream** projetada para plataformas de **e-commerce**, permitindo que as empresas acompanhem, analisem e otimizem a interação dos usuários em tempo real.

Os dados de Clickstream consistem em capturar os cliques e interações do usuário com os produtos exibidos no site, oferecendo informações valiosas sobre o comportamento do consumidor. Com a capacidade de processar esses dados em tempo real, essa solução fornece insights imediatos para ações estratégicas em campanhas de publicidade e personalização de experiência.


## Arquitetura
![Architecture](documentation/architecture.png)


## Tecnologias Utilizadas

Esta solução é baseada em uma arquitetura moderna que utiliza ferramentas robustas para processar e armazenar dados em tempo real:

- **Kafka**: Para ingestão e transporte de eventos de cliques.
- **KSQLDB**: Para agregação e processamento em tempo real dos dados capturados.
- **Kafka Connect**: Para integração e persistência dos dados processados em um banco de dados relacional.
- **PostgreSQL**: Para armazenamento e consulta de métricas agregadas.
- **Docker**: Para containerização de todos os serviços, facilitando o deployment e a escalabilidade.
- **Schema Registry**: Para gerenciamento de esquemas e validação de dados, garantindo a integridade das mensagens.

## Inicializando e Configurando a Solução

Siga os passos abaixo para inicializar e configurar a solução de Clickstream:

### Passo 1: Subir os Contêineres com Docker Compose
Acesse o diretório raiz da solução e execute o seguinte comando para construir e iniciar todos os contêineres:

```bash
$ docker-compose up --build
```
### Passo 2: Executar o Script de Configuração
Após o Docker Compose inicializar todos os contêineres, conceda permissão de execução ao script de configuração e execute-o:
```bash
$ chmod +x setup.sh
$ ./setup.sh
```
Este script irá configurar os tópicos Kafka, streams e tabelas necessários para o funcionamento da solução.

### Passo 3: Criando o Data Source no Grafana
Acesse o Grafana no link http://localhost:3000/ utilizando `admin` como usuário e senha. Crie um Data Source para PostgreSQL com os dados abaixo:
- **Host URL**: postgres:5432
- **Database Name**: db_metrics
- **Username**: user_metrics
- **Password**: password_metrics
- **TLS/SSL Mode**: Disable

Depois clique no botão `Save & Test`. Se tudo correr bem o DS será criado. Agora observe a URL do browser e verá algo similar a essa `http://localhost:3000/connections/datasources/edit/ee8bhesmqpczkb`. O último conjunto de informação (`ee8bhesmqpczkb`) é o UID do Data Source, faça uma cópia desse dado, ele será útil no próximo passo.

## Utilizando a Solução



## Entendendo o Fluxo de Dados
![Data Flow](documentation/data_flow.png)

A solução de Clickstream segue o fluxo de dados descrito abaixo, permitindo o processamento e análise em tempo real dos eventos de clique:


### 1. Captura do Evento de Clique
- **Origem**: Um evento de clique, representado por uma mensagem JSON, é gerado por uma aplicação de e-commerce.
- **Destino**: O evento é enviado para o tópico Kafka chamado `click_events`.


### 2. Kafka
- **Papel**: Atua como uma camada intermediária de armazenamento de eventos de clique. O tópico `click_events` recebe todos os eventos gerados pela aplicação de e-commerce.
- **Características**:
   - Os eventos são armazenados de forma distribuída.
   - Podem ser consumidos por diferentes sistemas downstream.


### 3. KsqlDB
#### **Stream (`click_events_stream`)**
- **Descrição**: O `click_events_stream` consome mensagens do tópico `click_events`.
- **Papel**:
   - Representa o fluxo contínuo e imutável dos eventos de clique.
   - Cada evento capturado pelo Kafka é processado pelo KSQLDB para extração de informações relevantes.

#### **Tabela (`click_counts_table`)**
- **Descrição**: A tabela `click_counts_table` agrega os eventos capturados pelo stream em janelas de 2 minutos.
- **Papel**:
   - Realiza a contagem de cliques por campanha (`campaignId`) dentro do intervalo de tempo especificado.
   - Os dados agregados são armazenados de forma estruturada para permitir consultas adicionais.
- **Detalhe Técnico**:
   - Essa tabela publica os dados agregados em um novo tópico Kafka chamado `click_counts_table_output`.

### 4. Kafka (Tópico `click_counts_table_output`)
- **Papel**: Recebe os dados agregados e sumarizados da tabela `click_counts_table`.
- **Formato dos Dados**:
   - **campaignId**: Identificação da campanha.
   - **click_count**: Total de cliques registrados na janela de 2 minutos.
   - **start_time** e **end_time**: Intervalo de tempo da agregação.

### 5. Kafka Connect
#### **Sink (`postgres-sink-connector`)**
- **Papel**:
   - Lê os dados do tópico `click_counts_table_output`.
   - Insere as informações agregadas no banco de dados PostgreSQL.
- **Transformação de Dados**:
   - Os campos de tempo (`start_time` e `end_time`) são convertidos de formato `BIGINT` para `timestamp` antes de serem inseridos no banco de dados.

### 6. PostgreSQL
- **Papel**: Armazena os dados agregados de cliques de forma estruturada e persistente.
- **Benefícios**:
   - Facilita a consulta e análise dos dados agregados.
   - Permite integração com ferramentas de visualização, como Grafana ou Looker, para criação de dashboards em tempo real.

### Resumo do Fluxo
1. Eventos de cliques são capturados e enviados para o Kafka.
2. O KSQLDB processa os eventos em tempo real:
   - Converte o fluxo contínuo (`click_events_stream`) em agregações sumarizadas (`click_counts_table`).
3. Os dados agregados são publicados no tópico Kafka `click_counts_table_output`.
4. O Kafka Connect lê os dados do tópico e insere no PostgreSQL para armazenamento persistente.
5. Ferramentas externas podem acessar o banco de dados para análises e visualizações.




