# Auction System - Go

Este é um sistema de leilões desenvolvido em Go, utilizando MongoDB como banco de dados. O sistema permite a criação de leilões e gerencia automaticamente o fechamento dos mesmos após um período determinado.

## 🚀 Tecnologias Utilizadas

- Go 1.23
- MongoDB
- Docker
- Docker Compose

## 🔧 Configuração do Ambiente de Desenvolvimento

### Pré-requisitos

- Docker
- Docker Compose
- Git

### 🛠️ Configuração Local

1. Clone o repositório:
```bash
git clone https://github.com/wanderlei2583/auction_leilao.git
cd auction_leilao
```

2. Configure as variáveis de ambiente:
Crie um arquivo `cmd/auction/.env` com as seguintes configurações:
```env
BATCH_INSERT_INTERVAL=20s
MAX_BATCH_SIZE=4
AUCTION_INTERVAL=20s

MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=admin
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
MONGODB_DB=auctions

AUCTION_DURATION=30s        # Duração do leilão
AUCTION_CHECK_INTERVAL=10s  # Intervalo de verificação
```

### 🚀 Executando o Projeto

1. Inicie os containers:
```bash
docker-compose up -d
```

2. Verifique se os serviços estão rodando:
```bash
docker-compose ps
```

### 📝 Testando a Aplicação

1. Criando um novo leilão:
```bash
curl -X POST http://localhost:8080/auction \
-H "Content-Type: application/json" \
-d '{
    "product_name": "iPhone 16 Pro Max",
    "category": "Eletrônicos",
    "description": "iPhone 16 Pro Max novo, na caixa lacrada",
    "condition": 0
}'
```

2. Verificando o status do leilão (substitua {auction_id} pelo ID retornado na criação):
```bash
curl http://localhost:8080/auction/{auction_id}
```
- O leilão deverá estar com o status "0"

3. Aguarde o tempo definido em AUCTION_DURATION para ver o leilão ser fechado automaticamente
* Após o tempo de leilão, verifique o status do leilão:
```bash
curl http://localhost:8080/auction/{auction_id}
```
- O leilão deverá ser fechado com o status "1"

4. Verificando os logs da aplicação:
```bash
docker-compose logs -f app
```

## 🧪Testes Automatizados
1. Caso não esteja startado, inicie os serviços usando Docker Compose:
```bash
docker compose up -d
```

2. Execute o teste específico do fechamento automático de leilões:
```bash
go test -v ./internal/infra/database/auction -run TestAutomaticAuctionClosure
```

#### O que o Teste Verifica

O teste `TestAutomaticAuctionClosure` valida o seguinte fluxo:
1. Conexão com o MongoDB
2. Criação de um leilão com status Active
3. Configuração de um timestamp passado para simular um leilão expirado
4. Verificação do fechamento automático do leilão após o período de expiração
5. Confirmação da mudança de status de Active para Completed

#### Configurações de Tempo
O teste utiliza as seguintes configurações de ambiente que podem ser ajustadas:
```env
AUCTION_DURATION=1m        # Duração do leilão
AUCTION_CHECK_INTERVAL=10s # Intervalo de verificação
```

#### Resultados Esperados
Um teste bem-sucedido mostrará logs indicando:
- Conexão bem-sucedida com o MongoDB
- Criação do leilão
- Status inicial do leilão
- Status final do leilão (deve ser Completed)

##### Exeplo de Resultados
![Teste fechamento automático de leilão](imagens/img_test.jpg)

Em caso de falha, o teste fornecerá informações detalhadas sobre qual etapa falhou e por quê.

## ⚙️ Configurações Importantes

### Variáveis de Ambiente

- `AUCTION_DURATION`: Define a duração do leilão (ex: 30s, 5m, 1h)
- `AUCTION_CHECK_INTERVAL`: Define o intervalo de verificação para fechamento dos leilões
- `MONGODB_URL`: URL de conexão com o MongoDB
- `MONGODB_DB`: Nome do banco de dados

### Docker Compose

O projeto utiliza dois serviços principais:
- `app`: Aplicação Go
- `mongodb`: Banco de dados MongoDB

## 🔒 Segurança

- Autenticação MongoDB configurada
- Rede Docker isolada para os serviços
- Variáveis sensíveis em arquivo .env

## 📝 Licença

Este projeto está sob a licença [MIT](LICENSE).

