package api

import (
	"fmt"
	"net/http"

	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/token"
	"github.com/gin-gonic/gin"
)

type createCollectionTrainRequest struct {
	ID      int64 `uri:"id" form:"required,min=1"`
	TrainID int64 `uri:"train_id" form:"required,min=1"`
}

func (server *Server) createCollectionTrain(ctx *gin.Context) {
	var req createCollectionTrainRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateCollectionTrainParams{
		UserID:  req.ID,
		TrainID: req.TrainID,
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if arg.UserID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to add to this collection"))
		return
	}

	collectionTrain, err := server.Store.CreateCollectionTrain(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, collectionTrain)
}

type deleteCollectionTrainRequest struct {
	ID      int64 `uri:"id" form:"required,min=1"`
	TrainID int64 `uri:"train_id" form:"required,min=1"`
}

func (server *Server) deleteCollectionTrain(ctx *gin.Context) {
	var req deleteCollectionTrainRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if req.ID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to add to this collection"))
		return
	}

	arg := db.DeleteCollectionTrainParams{
		UserID:  authPayload.UserID,
		TrainID: req.TrainID,
	}

	err := server.Store.DeleteCollectionTrain(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Train successfully removed from collection.")
}
