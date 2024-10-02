package auction

import (
	"context"
	"fmt"
	"github.com/betonetotbo/go-expert-labs-auctions/configuration/logger"
	"github.com/betonetotbo/go-expert-labs-auctions/internal/entity/auction_entity"
	"github.com/betonetotbo/go-expert-labs-auctions/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
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
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	go ar.autocloseAuction(auctionEntityMongo)

	return nil
}

func (ar *AuctionRepository) autocloseAuction(auction *AuctionEntityMongo) {
	select {
	case <-time.After(GetAuctionInterval()):
		logger.Info(fmt.Sprintf("Completing auction (expired): %s", auction.Id))

		filter := bson.D{
			primitive.E{
				Key:   "_id",
				Value: auction.Id,
			},
		}
		update := bson.D{
			primitive.E{
				Key: "$set",
				Value: bson.D{
					primitive.E{
						Key:   "status",
						Value: auction_entity.Completed,
					},
				},
			},
		}

		_, err := ar.Collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			logger.Error(fmt.Sprintf("Error trying to complete auction: %s", auction.Id), err)
		}
	}
}

func GetAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}

	return duration
}
