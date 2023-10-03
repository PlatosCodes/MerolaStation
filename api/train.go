package api

import (
	"database/sql"
	"net/http"

	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/token"
	"github.com/gin-gonic/gin"
)

type createTrainRequest struct {
	ModelNumber string `json:"model_number" binding:"required"`
	Name        string `json:"train_name" binding:"required"`
}

func (server *Server) createTrain(ctx *gin.Context) {
	var req createTrainRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateTrainParams{
		ModelNumber: req.ModelNumber,
		Name:        req.Name,
	}

	train, err := server.Store.CreateTrain(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, train)
}

type getTrainRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTrain(ctx *gin.Context) {
	var req getTrainRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	id := req.ID

	train, err := server.Store.GetTrain(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, train)
}

type getTrainByModelRequest struct {
	ModelNumber string `uri:"model_number" binding:"required,min=1"`
}

func (server *Server) getTrainByModel(ctx *gin.Context) {
	var req getTrainByModelRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	train, err := server.Store.GetTrainByModel(ctx, req.ModelNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, train)
}

type listTrainRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=25"`
}

func (server *Server) listTrain(ctx *gin.Context) {
	var req listTrainRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListTrainsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	trains, err := server.Store.ListTrains(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, trains)
}

type updateTrainValueRequest struct {
	ID    int64 `json:"id" binding:"required,min=1"`
	Value int64 `json:"value" binding:"required,min=1"`
}

func (server *Server) updateTrainValue(ctx *gin.Context) {
	var req updateTrainValueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload).UserID

	arg := db.UpdateTrainValueParams{
		ID:    req.ID,
		Value: req.Value,
	}

	err := server.Store.UpdateTrainValue(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, req.Value)
}
