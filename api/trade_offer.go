package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/token"
	"github.com/gin-gonic/gin"
)

type createTradeOfferRequest struct {
	OfferedTrain        int64 `json:"offered_train" binding:"required,min=1"`
	OfferedTrainOwner   int64 `json:"offered_train_owner" binding:"required,min=1"`
	RequestedTrain      int64 `json:"requested_train" binding:"required,min=1"`
	RequestedTrainOwner int64 `json:"requested_train_owner" binding:"required,min=1"`
}

func (server *Server) createTradeOffer(ctx *gin.Context) {
	var req createTradeOfferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if req.OfferedTrainOwner != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to make this trade offer"))
		return
	}

	arg := db.CreateTradeOfferParams{
		OfferedTrain:        req.OfferedTrain,
		OfferedTrainOwner:   req.OfferedTrainOwner,
		RequestedTrain:      req.RequestedTrain,
		RequestedTrainOwner: req.RequestedTrainOwner,
	}

	tradeOffer, err := server.Store.CreateTradeOffer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tradeOffer)
}

type getTradeOfferByIDRequest struct {
	ID int64 `uri:"id" form:"required,min=1"`
}

func (server *Server) getTradeOfferByID(ctx *gin.Context) {
	var req getTradeOfferByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tradeOffer, err := server.Store.GetTradeOfferByTradeID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if tradeOffer.OfferedTrainOwner != authPayload.UserID && tradeOffer.RequestedTrainOwner != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to view this trade offer"))
		return
	}

	ctx.JSON(http.StatusOK, tradeOffer)
}

type UserTradeOffersPaginationQuery struct {
	PageID   int `form:"page_id" binding:"required"`
	PageSize int `form:"page_size" binding:"required"`
}

type UserTradeOffersPath struct {
	OfferedTrainOwner int64 `uri:"offered_train_owner" binding:"required,min=1"`
}

func (server *Server) listUserTradeOffers(ctx *gin.Context) {
	var req UserTradeOffersPaginationQuery
	var path UserTradeOffersPath

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUserTradeOffersParams{
		OfferedTrainOwner: path.OfferedTrainOwner,
		Limit:             int32(req.PageSize),
		Offset:            int32((req.PageID - 1) * req.PageSize),
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if arg.OfferedTrainOwner != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to view this collection"))
		return
	}

	tradeOffers, err := server.Store.ListUserTradeOffers(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tradeOffers)
}

type UserTradeRequestsPaginationQuery struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

type UserTradeRequestsPath struct {
	RequestedTrainOwner int64 `uri:"requested_train_owner" binding:"required,min=1"`
}

func (server *Server) listUserTradeRequests(ctx *gin.Context) {
	var req UserTradeRequestsPaginationQuery
	var path UserTradeRequestsPath

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if path.RequestedTrainOwner != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to view this collection"))
		return
	}

	tradeRequests, err := server.Store.ListUserTradeRequests(ctx, db.ListUserTradeRequestsParams{
		RequestedTrainOwner: path.RequestedTrainOwner,
		Limit:               req.PageSize,
		Offset:              (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tradeRequests)
}

type AllUserTradeOffersPaginationQuery struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

type AllUserTradeOffersPath struct {
	OfferedTrainOwner int64 `uri:"offered_train_owner" binding:"required,min=1"`
}

func (server *Server) listAllUserTradeOffers(ctx *gin.Context) {
	var req AllUserTradeOffersPaginationQuery
	var path AllUserTradeOffersPath

	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if path.OfferedTrainOwner != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to view this collection"))
		return
	}

	tradeOffers, err := server.Store.ListAllUserTradeOffers(ctx, db.ListAllUserTradeOffersParams{
		OfferedTrainOwner: path.OfferedTrainOwner,
		Limit:             req.PageSize,
		Offset:            (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tradeOffers)
}

// ListCollectionTrainTradeOffers
type listCollectionTrainTradeOffersQueryParams struct {
	Limit  int32 `form:"limit,default=10,min=1,max=100"`
	Offset int32 `form:"offset,default=0,min=0"`
}

type listCollectionTrainTradeOffersPath struct {
	TrainID int64 `uri:"train_id" form:"required,min=1"`
}

func (server *Server) listCollectionTrainTradeOffers(ctx *gin.Context) {
	var req listCollectionTrainTradeOffersQueryParams
	var path listCollectionTrainTradeOffersPath

	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	tradeOffers, err := server.Store.ListCollectionTrainTradeOffers(ctx, db.ListCollectionTrainTradeOffersParams{
		RequestedTrain: path.TrainID,
		Limit:          req.Limit,
		Offset:         req.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tradeOffers)
}

// // ListTradeOffers
// type listTradeOffersRequest struct {
// 	TrainID int64 `uri:"train_id" form:"required,min=1"`
// 	Limit   int32 `form:"limit,default=10,min=1,max=100"`
// 	Offset  int32 `form:"offset,default=0,min=0"`
// }

// func (server *Server) listTradeOfferRequests(ctx *gin.Context) {
// 	var req listTrainTradeOffersRequest
// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	tradeOffers, err := server.Store.ListUserTradeRequests(ctx, db.ListUserTradeRequestsParams{
// 		RequestedTrainOwner: req,
// 		Limit:        req.Limit,
// 		Offset:       req.Offset,
// 	})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, tradeOffers)
// }

type deleteTradeOfferRequest struct {
	ID                int64 `json:"id" binding:"required,min=1"`
	OfferedTrainOwner int64 `json:"offered_train_owner" binding:"required,min=1"`
}

func (server *Server) deleteTradeOffer(ctx *gin.Context) {
	var req deleteTradeOfferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if req.OfferedTrainOwner != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to view this collection"))
		return
	}

	err := server.Store.DeleteTradeOffer(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Trade offer successfully deleted.")
}
