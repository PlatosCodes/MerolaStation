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

type WishlistPaginationQuery struct {
	PageID   int `form:"page_id" binding:"required"`
	PageSize int `form:"page_size" binding:"required"`
}

type UserWishlistPath struct {
	UserID int64 `uri:"id" binding:"required"`
}

func (server *Server) listUserWishlist(ctx *gin.Context) {
	var req WishlistPaginationQuery
	var path UserWishlistPath

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUserWishlistParams{
		UserID: path.UserID,
		Limit:  int32(req.PageSize),
		Offset: int32((req.PageID - 1) * req.PageSize),
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if arg.UserID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to view this wishlist"))
		return
	}

	wishlistTrains, err := server.Store.ListUserWishlist(context.Background(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	trains := []db.Train{}
	for _, ct := range wishlistTrains {
		train, err := server.Store.GetTrain(context.Background(), ct.TrainID)
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

	ctx.JSON(http.StatusOK, trains)
}

type getUserWishlistTrainRequest struct {
	ID      int64 `uri:"id" form:"required,min=1"`
	TrainID int64 `uri:"train_id" form:"required,min=1"`
}

func (server *Server) getUserWishlistTrain(ctx *gin.Context) {
	var path getUserWishlistTrainRequest

	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if path.ID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to view this wishlist"))
		return
	}

	wishlistTrain, err := server.Store.GetWishlistTrain(context.Background(), db.GetWishlistTrainParams{
		UserID:  path.ID,
		TrainID: path.TrainID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusOK, gin.H{"isInWishlist": false})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"isInWishlist": true, "train": wishlistTrain})
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

type WishlistTrainWithCollectionStatus struct {
	WishlistTrain  db.WishlistTrain
	IsInCollection bool `json:"is_in_collection"`
}

func (server *Server) getUserWishlistWithCollectionStatus(ctx *gin.Context) {
	var path UserCollectionPath // Using the same structure since it only contains UserID

	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if path.UserID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to view this wishlist"))
		return
	}

	trains, err := server.Store.GetUserWishlistWithCollectionStatus(context.Background(), path.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, trains)
}
