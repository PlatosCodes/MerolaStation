package api

import (
	"database/sql"
	"log"
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

type getTrainDetailRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTrainDetail(ctx *gin.Context) {
	var req getTrainDetailRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	id := req.ID

	train, err := server.Store.GetTrainDetail(ctx, db.GetTrainDetailParams{ID: id, UserID: authPayload.UserID})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Println(train)
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

type ListUserTrainsWithPagesResponse struct {
	TotalCount int64                  `json:"total_count"`
	Trains     []db.ListUserTrainsRow `json:"trains"`
}

func (server *Server) listUserTrainsWithPages(ctx *gin.Context) {
	var req listTrainRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListUserTrainsParams{
		UserID: authPayload.UserID, // <-- Include the UserID in the params
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	trains, err := server.Store.ListUserTrains(ctx, arg) // <-- Call the new ListUserTrains method
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Query total count of trains
	totalCount, err := server.Store.GetTotalTrainCount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := ListUserTrainsWithPagesResponse{
		TotalCount: totalCount,
		Trains:     trains,
	}

	ctx.JSON(http.StatusOK, response)
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

func (server *Server) listUserTrains(ctx *gin.Context) {
	var req listTrainRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListUserTrainsParams{
		UserID: authPayload.UserID, // <-- Include the UserID in the params
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	trains, err := server.Store.ListUserTrains(ctx, arg) // <-- Call the new ListUserTrains method
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

type updateTrainValuesBatchRequest struct {
	Updates []struct {
		ID    int64 `json:"id" binding:"required,min=1"`
		Value int64 `json:"value" binding:"required,min=1"`
	} `json:"updates" binding:"required,dive"`
}

func (server *Server) updateTrainsValuesBatch(ctx *gin.Context) {
	var req updateTrainValuesBatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload).UserID

	ids := make([]int64, len(req.Updates))
	values := make([]int64, len(req.Updates))
	for i, update := range req.Updates {
		ids[i] = update.ID
		values[i] = update.Value
	}

	err := server.Store.UpdateTrainsValuesBatch(ctx, ids, values)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, req)
}

type searchTrainByModelNumberSuggestionRequest struct {
	ModelNumber string `form:"model_number" binding:"required" default:""`
	PageSize    int    `form:"page_size" binding:"required" default:"10"`
	PageID      int    `form:"page_id" binding:"required" default:"1"`
}

type SearchTrainsByModelNumberResponse struct {
	Trains []db.SearchTrainsByModelNumberSuggestionsRow `json:"trains"`
}

func (server *Server) searchTrainsByModelNumberSuggestions(ctx *gin.Context) {

	var req searchTrainByModelNumberSuggestionRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Printf("Bind error: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	limit := int32(req.PageSize)

	// Build the ILIKE pattern for model number
	modelNumberPattern := sql.NullString{
		String: req.ModelNumber + "%",
		Valid:  true,
	}

	// Fetch trains based on model number
	searchParams := db.SearchTrainsByModelNumberSuggestionsParams{
		Column1: modelNumberPattern,
		Limit:   limit,
		Offset:  int32((req.PageID - 1) * req.PageSize),
	}

	trains, err := server.Store.SearchTrainsByModelNumberSuggestions(ctx, searchParams)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := SearchTrainsByModelNumberResponse{
		Trains: trains,
	}

	ctx.JSON(http.StatusOK, response)

}

type searchTrainByNameSuggestionRequest struct {
	Name     string `form:"name" binding:"required"`
	PageSize int    `form:"page_size" binding:"required" default:"10"`
	PageID   int    `form:"page_id" binding:"required" default:"1"`
}

// This can be the same response type as search by model number
type SearchTrainsByNameResponse struct {
	Trains []db.SearchTrainsByNameSuggestionsRow `json:"trains"`
}

func (server *Server) searchTrainsByNameSuggestions(ctx *gin.Context) {
	var req searchTrainByNameSuggestionRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Printf("Bind error: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	limit := int32(req.PageSize)

	namePattern := sql.NullString{
		String: "%" + req.Name + "%",
		Valid:  true,
	}

	log.Println("rsp:", namePattern)

	searchParams := db.SearchTrainsByNameSuggestionsParams{
		Column1: namePattern,
		Limit:   limit,
		Offset:  int32((req.PageID - 1) * req.PageSize),
	}

	trains, err := server.Store.SearchTrainsByNameSuggestions(ctx, searchParams)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := SearchTrainsByNameResponse{
		Trains: trains,
	}

	ctx.JSON(http.StatusOK, response)
}
