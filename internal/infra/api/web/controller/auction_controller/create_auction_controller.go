package auction_controller

import (
	"context"
	"github.com/betonetotbo/go-expert-labs-auctions/configuration/rest_err"
	"github.com/betonetotbo/go-expert-labs-auctions/internal/infra/api/web/validation"
	"github.com/betonetotbo/go-expert-labs-auctions/internal/usecase/auction_usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuctionController struct {
	auctionUseCase auction_usecase.AuctionUseCaseInterface
}

func NewAuctionController(auctionUseCase auction_usecase.AuctionUseCaseInterface) *AuctionController {
	return &AuctionController{
		auctionUseCase: auctionUseCase,
	}
}

func (u *AuctionController) CreateAuction(c *gin.Context) {
	var auctionInputDTO auction_usecase.AuctionInputDTO

	if err := c.ShouldBindJSON(&auctionInputDTO); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	out, err := u.auctionUseCase.CreateAuction(context.Background(), auctionInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusCreated, out)
}
