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
MONGODB_DATABASE=auctions

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

3. Aguarde o tempo definido em AUCTION_DURATION para ver o leilão ser fechado automaticamente

4. Verificando os logs da aplicação:
```bash
docker-compose logs -f app
```

## ⚙️ Configurações Importantes

### Variáveis de Ambiente

- `AUCTION_DURATION`: Define a duração do leilão (ex: 30s, 5m, 1h)
- `AUCTION_CHECK_INTERVAL`: Define o intervalo de verificação para fechamento dos leilões
- `MONGODB_URL`: URL de conexão com o MongoDB
- `MONGODB_DATABASE`: Nome do banco de dados

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

