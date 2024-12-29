# Solução de Clickstream para E-commerce

## Introdução

Esta é uma solução de **Clickstream** projetada para plataformas de **e-commerce**, permitindo que as empresas acompanhem, analisem e otimizem a interação dos usuários em tempo real.

Os dados de Clickstream consistem em capturar os cliques e interações do usuário com os produtos exibidos no site, oferecendo informações valiosas sobre o comportamento do consumidor. Com a capacidade de processar esses dados em tempo real, essa solução fornece insights imediatos para ações estratégicas em campanhas de publicidade e personalização de experiência.

<br>

## Arquitetura
![Architecture](documentation/architecture.png)

### Componentes da Arquitetura

| **Componente**          | **Descrição**                                                                                                      |
|--------------------------|--------------------------------------------------------------------------------------------------------------------|
| items-frontend           | Interface usada por usuários de e-commerce para interagir com os itens disponíveis.                                |
| items-api                | Backend responsável por buscar itens e registrar cliques, chamando a **clickstream-api**.                          |
| clickstream-api          | Envia os eventos de cliques para o **Kafka**.                                                                      |
| Kafka                    | - Armazena eventos de cliques no tópico `click_events`. Serve como intermediário para sistemas downstream.         |
| KSQLDB                   | - Consome eventos do Kafka. Agrega dados e publica em um tópico de saída (`click_counts_table_output`).            |
| Schema Registry          | - Gerencia os esquemas das mensagens trafegadas no Kafka. Garante consistência e validação.                        |
| Kafka Connect            | Consome dados agregados do Kafka e insere no banco de dados PostgreSQL.                                            |
| Postgres                 | Armazena dados agregados para persistência.                                                                        |
| Grafana                  | Lê os dados do PostgreSQL e permite a visualização de métricas em dashboards acessados por usuários de marketing.  |
| Connect UI               | Interface para monitorar o Kafka Connect.                                                                          |
| Registry UI              | Interface para gerenciar esquemas no Schema Registry.                                                              |

<br>

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

Siga os passos abaixo para criar o Data Source no Grafana:

1. **Acesse o Grafana**
  - Link: [http://localhost:3000/](http://localhost:3000/)
  - **Usuário**: `admin`
  - **Senha**: `admin`

2. **Navegue até o menu Data Sources**  
   No menu lateral, vá até **`Connections > Data Sources`**.

3. **Adicionar um novo Data Source**  
   Clique no botão **`Add new data source`** e selecione a opção **`PostgreSQL`**.

4. **Configurar o Data Source**  
   Preencha os campos com as seguintes informações:
  - **Host URL**: `postgres:5432`
  - **Database Name**: `db_metrics`
  - **Username**: `user_metrics`
  - **Password**: `password_metrics`
  - **TLS/SSL Mode**: `Disable`

5. **Salvar e testar a conexão**  
   Após preencher os campos, clique no botão **`Save & Test`**. Se tudo estiver correto, o Data Source será criado com sucesso.

6. **Obter o UID do Data Source**  
   Após a criação, observe a URL no navegador. Ela será algo como:  
   `http://localhost:3000/connections/datasources/edit/ee8bhesmqpczkb`
  - O último conjunto de caracteres (`ee8bhesmqpczkb`) é o **UID do Data Source**.
  - Copie este valor, pois ele será necessário no próximo passo.


### Passo 4: Configurando e Importando o Dashboard
Siga os passos abaixo para configurar e importar o dashboard no Grafana:

1. **Editar o arquivo `dashboard.json`**  
   Localize o arquivo `dashboard.json` na raiz do projeto e substitua o texto `<datasource_uid_grafana>` pelo UID do Data Source criado no passo anterior.

2. **Acessar o menu Dashboards**  
   No menu lateral do Grafana, clique em **`Dashboards`**.

3. **Criar um novo dashboard**
  - Clique no botão **`New`**.
  - Em seguida, selecione a opção **`New dashboard`**.

4. **Importar o dashboard**
  - Clique no botão **`Import dashboard`**.
  - Abra o arquivo `dashboard.json`, copie o conteúdo e cole no campo **`Import via dashboard JSON model`**.

5. **Salvar o dashboard**
  - Por último, clique no botão **`Save`** para finalizar a importação.

<br>

## Utilizando a Solução
> **Nota:** Para acessar os recursos da solução, é necessário que todas as configurações tenham sido realizadas com sucesso.



<br>

## Entendendo o Fluxo de Dados
![Data Flow](documentation/data_flow.png)

O fluxo dos dados funciona de forma simples e eficiente. Toda vez que alguém clica em um item no site, um evento em formato JSON é gerado e enviado para o Kafka, que armazena esses eventos de forma organizada, pronta para ser consumida por outros sistemas.

O KSQLDB entra em ação consumindo esses eventos e processando as informações. Ele cria um fluxo contínuo (stream) que processa os cliques em tempo real e, em seguida, organiza esses cliques em uma tabela que conta o número de cliques por campanha a cada 2 minutos. Esses dados agregados são então enviados para um novo tópico no Kafka.

Posteriormente, o Kafka Connect lê os dados do tópico de agregação e insere no banco de dados PostgreSQL. Durante essa etapa, são aplicadas transformações nos dados, como a conversão dos campos de tempo (`start_time` e `end_time`) para o formato de timestamp.

Por fim, o PostgreSQL armazena os dados de forma estruturada e organizada. Esses dados podem ser acessados por ferramentas como Grafana ou Looker, que possibilitam a criação de dashboards atrativos e úteis para análises em tempo real.


