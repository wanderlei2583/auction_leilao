package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fullcycle-auction_go/internal/entity/auction_entity"
)

func TestAutomaticAuctionClosure(t *testing.T) {
	ctx := context.Background()

	mongoURL := "mongodb://admin:admin@localhost:27017/auctions?authSource=admin"

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		t.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		t.Fatalf("Erro ao pingar MongoDB: %v", err)
	}
	t.Log("Conectado ao MongoDB com sucesso")

	database := client.Database("auctions")
	repository := NewAuctionRepository(database)

	os.Setenv("AUCTION_DURATION", "1m")
	os.Setenv("AUCTION_CHECK_INTERVAL", "10s")

	t.Log("Criando leilão de teste...")

	auctionEntity, internalErr := auction_entity.CreateAuction(
		"Test Product Name",
		"Test Category",
		"This is a test product description",
		auction_entity.New,
	)
	if internalErr != nil {
		t.Fatalf("Erro ao criar entidade do leilão: %s", internalErr.Error())
	}
	t.Log("Entidade do leilão criada com sucesso")

	auctionEntity.Timestamp = time.Now().Add(-2 * time.Minute)

	if internalErr := repository.CreateAuction(ctx, auctionEntity); internalErr != nil {
		t.Fatalf("Erro ao inserir leilão no MongoDB: %s", internalErr.Error())
	}
	t.Log("Leilão inserido com sucesso no MongoDB")

	var initialAuction AuctionEntityMongo
	err = repository.Collection.FindOne(ctx, bson.M{"_id": auctionEntity.Id}).
		Decode(&initialAuction)
	if err != nil {
		t.Fatalf("Erro ao verificar leilão inicial: %v", err)
	}
	t.Logf("Status inicial do leilão: %v", initialAuction.Status)

	t.Log("Aguardando processamento automático...")
	time.Sleep(15 * time.Second)

	var finalAuction AuctionEntityMongo
	err = repository.Collection.FindOne(ctx, bson.M{"_id": auctionEntity.Id}).
		Decode(&finalAuction)
	if err != nil {
		t.Fatalf("Erro ao verificar leilão final: %v", err)
	}
	t.Logf("Status final do leilão: %v", finalAuction.Status)

	if finalAuction.Status != auction_entity.Completed {
		t.Errorf(
			"Status esperado %v, obtido %v",
			auction_entity.Completed,
			finalAuction.Status,
		)
	}
}
