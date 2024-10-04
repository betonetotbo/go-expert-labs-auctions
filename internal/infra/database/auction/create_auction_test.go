package auction

import (
	"context"
	"github.com/betonetotbo/go-expert-labs-auctions/internal/entity/auction_entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"time"
)

func TestCreateAuction(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("AutoClose", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateSuccessResponse())

		repo := &AuctionRepository{
			Collection: t.Coll,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		t.Setenv("AUCTION_INTERVAL", "500ms")

		entity := &auction_entity.Auction{
			Id:          uuid.New().String(),
			Status:      auction_entity.Active,
			Description: "some long description",
			Category:    "auction category",
			ProductName: "the product name",
			Timestamp:   time.Now(),
		}
		err := repo.CreateAuction(ctx, entity)

		if assert.Nil(t, err, "Auction creation expected to be succeed") {
			if ev := t.GetSucceededEvent(); assert.NotNil(t, ev) && assert.Equal(t, "insert", ev.CommandName) {
				assert.Nil(t, t.GetSucceededEvent(), "No more events expected")

				t.AddMockResponses(mtest.CreateSuccessResponse())

				// give up a chance to run autoclose goroutine
				time.Sleep(time.Second)

				if ev = t.GetSucceededEvent(); assert.NotNil(t, ev, "Update event expected for complete auction") {
					assert.Equal(t, "update", ev.CommandName, "Update event expected for complete auction")
				}
			}
		}
	})
}
