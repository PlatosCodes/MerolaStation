package api

import (
	"context"
	"database/sql"
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

type CollectionPaginationQuery struct {
	PageID   int `form:"page_id" binding:"required"`
	PageSize int `form:"page_size" binding:"required"`
}

type UserCollectionPath struct {
	UserID int64 `uri:"id" binding:"required"`
}

func (server *Server) listUserCollection(ctx *gin.Context) {
	var req CollectionPaginationQuery
	var path UserCollectionPath

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUserCollectionParams{
		UserID: path.UserID,
		Limit:  int32(req.PageSize),
		Offset: int32((req.PageID - 1) * req.PageSize),
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if arg.UserID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to view this collection"))
		return
	}

	collectionTrains, err := server.Store.ListUserCollection(context.Background(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	trains := []db.GetTrainDetailRow{}
	for _, ct := range collectionTrains {
		train, err := server.Store.GetTrainDetail(context.Background(), db.GetTrainDetailParams{ID: ct.TrainID, UserID: authPayload.UserID})
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		trains = append(trains, train)
	}

	// Fetch total collection value
	totalValue := int64(0)

	if len(trains) > 0 {
		totalValue, err = server.Store.GetTotalCollectionValue(context.Background(), authPayload.UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// Create a new struct for the response which includes both the list of trains and the total collection value
	type Response struct {
		Trains     []db.GetTrainDetailRow `json:"trains"`
		TotalValue int64                  `json:"totalValue"`
	}

	resp := Response{
		Trains:     trains,
		TotalValue: totalValue,
	}

	ctx.JSON(http.StatusOK, resp)
}

type getUserCollectionTrainRequest struct {
	ID      int64 `uri:"id" form:"required,min=1"`
	TrainID int64 `uri:"train_id" form:"required,min=1"`
}

func (server *Server) getUserCollectionTrain(ctx *gin.Context) {
	var path getUserCollectionTrainRequest

	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if path.ID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to view this collection"))
		return
	}

	collectionTrain, err := server.Store.GetCollectionTrain(context.Background(), db.GetCollectionTrainParams{
		UserID:  path.ID,
		TrainID: path.TrainID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusOK, gin.H{"isInCollection": false})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"isInCollection": true, "train": collectionTrain})
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
	fmt.Println(req.ID, authPayload.UserID)
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

type CollectionTrainWithWishlistStatus struct {
	CollectionTrain db.CollectionTrain
	IsInWishlist    bool `json:"is_in_wishlist"`
}

func (server *Server) getUserCollectionWithWishlistStatus(ctx *gin.Context) {
	var path UserCollectionPath

	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if path.UserID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to view this collection"))
		return
	}

	trains, err := server.Store.GetUserCollectionWithWishlistStatus(context.Background(), path.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, trains)
}
