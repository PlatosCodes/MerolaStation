package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/token"
	"github.com/gin-gonic/gin"
)

type CreateTradeTransactionParams struct {
	TradeOfferID int64 `json:"id" binding:"required,min=1"`
}

// createTrade uses TradeTX to create a new trade
func (server *Server) createTrade(ctx *gin.Context) {
	var req CreateTradeTransactionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	to, err := server.Store.GetTradeOfferByTradeID(ctx, req.TradeOfferID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if to.OfferedTrainOwner != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to make this trade offer"))
		return
	}

	ct1, err := server.Store.GetCollectionTrainforUpdate(ctx, db.GetCollectionTrainforUpdateParams{
		UserID:  to.OfferedTrainOwner,
		TrainID: to.OfferedTrain,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ct2, err := server.Store.GetCollectionTrainforUpdate(ctx, db.GetCollectionTrainforUpdateParams{
		UserID:  to.RequestedTrainOwner,
		TrainID: to.RequestedTrain,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.TradeTxParams{
		TradeOfferID:   to.ID,
		OfferedTrain:   ct1,
		RequestedTrain: ct2,
	}

	tradeResult, err := server.Store.TradeTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tradeResult)
}

type getTradeTransactionRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTradeTransaction(ctx *gin.Context) {
	var req getTradeTransactionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	id := req.ID

	tt, err := server.Store.GetTradeTransaction(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if tt.OfferedTrainOwner != authPayload.UserID && tt.RequestedTrainOwner != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to make this trade offer"))
		return
	}

	ctx.JSON(http.StatusOK, tt)
}
