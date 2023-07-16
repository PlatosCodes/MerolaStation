package api

import (
	"fmt"
	"net/http"

	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/token"
	"github.com/gin-gonic/gin"
)

type createWishlistTrainRequest struct {
	ID      int64 `uri:"id" form:"required,min=1"`
	TrainID int64 `uri:"train_id" form:"required,min=1"`
}

func (server *Server) createWishlistTrain(ctx *gin.Context) {
	var req createWishlistTrainRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateWishlistTrainParams{
		UserID:  req.ID,
		TrainID: req.TrainID,
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if arg.UserID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to add to this wishlist"))
		return
	}

	wishlistTrain, err := server.Store.CreateWishlistTrain(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, wishlistTrain)
}

type deleteWishlistTrainRequest struct {
	ID      int64 `uri:"id" form:"required,min=1"`
	TrainID int64 `uri:"train_id" form:"required,min=1"`
}

func (server *Server) deleteWishlistTrain(ctx *gin.Context) {
	var req deleteWishlistTrainRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if req.ID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to add to this wishlist"))
		return
	}

	arg := db.DeleteWishlistTrainParams{
		UserID:  authPayload.UserID,
		TrainID: req.TrainID,
	}

	err := server.Store.DeleteWishlistTrain(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Train successfully removed from wishlist.")
}
