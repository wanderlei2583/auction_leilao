package auction

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
	mutex      sync.Mutex
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repo := &AuctionRepository{
		Collection: database.Collection("auctions"),
	}

	go repo.checkExpiredAuctions(context.Background())

	return repo
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction,
) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}

	ar.mutex.Lock()
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	ar.mutex.Unlock()

	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError(
			"Error trying to insert auction",
		)
	}

	return nil
}

func (ar *AuctionRepository) checkExpiredAuctions(ctx context.Context) {
	interval := getCheckInterval()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ar.closeExpiredAuctions(ctx)
		}
	}
}

func (ar *AuctionRepository) closeExpiredAuctions(ctx context.Context) {
	expirationTime := time.Now().Add(-getAuctionDuration()).Unix()

	filter := bson.M{
		"status":    0,
		"timestamp": bson.M{"$lt": expirationTime},
	}

	update := bson.M{
		"$set": bson.M{
			"status": 1,
		},
	}

	ar.mutex.Lock()
	result, err := ar.Collection.UpdateMany(ctx, filter, update)
	ar.mutex.Unlock()

	if err != nil {
		logger.Error("Error trying to close expired auctions", err)
		return
	}

	if result.ModifiedCount > 0 {
		logger.Info(
			fmt.Sprintf("Closed %d expired auctions", result.ModifiedCount),
		)
	}
}

func getAuctionDuration() time.Duration {
	durationStr := os.Getenv("AUCTION_DURATION_MINUTES")
	if durationStr == "" {
		return 5 * time.Minute
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		logger.Error("Invalid auction duration format, using default 5m", err)
		return 5 * time.Minute
	}

	return duration
}

func getCheckInterval() time.Duration {
	intervalStr := os.Getenv("AUCTION_INTERVAL")
	if intervalStr == "" {
		return 1 * time.Minute
	}

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		logger.Error("Invalid check interval format, using default 1m", err)
		return 1 * time.Minute
	}

	return interval
}
