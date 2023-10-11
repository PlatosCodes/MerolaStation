package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/token"
	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	UserID            int64     `json:"id"`
	Username          string    `json:"username"`
	CreatedAt         time.Time `json:"created_at"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		UserID:            user.ID,
		Username:          user.Username,
		CreatedAt:         user.CreatedAt,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}

// createUser uses RegisterTX to create a new user
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}

	registerResult, err := server.Store.RegisterTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Println("register result: ", registerResult)
	activationToken, activationPayload, err := server.tokenMaker.CreateToken(
		registerResult.User.ID,
		registerResult.User.Username,
		server.config.RefreshTokenDuration*7,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	activationInfo, err := server.Store.InsertActivationToken(ctx, db.InsertActivationTokenParams{
		UserID:          registerResult.User.ID,
		ActivationToken: activationToken,
		ExpiresAt:       activationPayload.ExpiresAt.Time,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	dbs := map[string]interface{}{
		"activationToken": activationInfo.ActivationToken,
		"userID":          activationPayload.UserID,
		"username":        registerResult.User.Username,
	}
	log.Println("dbs ", registerResult)

	// can make async later
	err = server.mailer.Send(registerResult.User.Email, "user_welcome.tmpl", dbs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(registerResult.User)

	log.Println("rsp: ", rsp)

	ctx.JSON(http.StatusOK, rsp)
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getUserResponse struct {
	ID                int64
	CreatedAt         time.Time
	Username          string
	Email             string
	PasswordChangedAt time.Time
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if req.ID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, fmt.Errorf("you do not have permission to view this collection"))
		return
	}

	user, err := server.Store.GetUser(ctx, authPayload.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := getUserResponse{
		ID:                user.ID,
		CreatedAt:         user.CreatedAt,
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type activateUserRequest struct {
	UserID          int64  `json:"user_id" binding:"required"`
	ActivationToken string `json:"activation_token" binding:"required"`
}

func (server *Server) activateUser(ctx *gin.Context) {
	var req activateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println(1, "req:", req)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.tokenMaker.VerifyToken(req.ActivationToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	storedToken, err := server.Store.GetActivationToken(ctx, req.ActivationToken)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	if storedToken.UserID == req.UserID {
		err = server.Store.ActivateUser(ctx, storedToken.UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	err = server.Store.DeleteActivationToken(ctx, storedToken.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "User successfully activated.")
}

type listUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUser(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.Store.ListUsers(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse((err)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if payload.UserID != 1 {
		err := fmt.Errorf("user does not have permission to list users")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userResponses := []userResponse{}
	for _, ct := range users {
		userResponse := newUserResponse(ct)
		userResponses = append(userResponses, userResponse)
	}

	ctx.JSON(http.StatusOK, userResponses)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	fmt.Println(req)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Use this if want to send a more generic message in response:
	// if err := ctx.ShouldBindJSON(&req); err != nil {
	// 	log.Printf("Error binding login request: %v", err) // Log the detailed error
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error. Please try again later."})
	// 	return
	// }

	user, err := server.Store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !user.Activated {
		ctx.JSON(http.StatusUnauthorized, "Please use the link that was sent to your email in order to activate your account.")
		log.Println("yoho")
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.Store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     refreshPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
		User:                  newUserResponse(user),
	}

	// Set the refresh token as an HttpOnly cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  refreshPayload.ExpiresAt.Time,
		HttpOnly: true,
		Path:     "/",
	})

	ctx.JSON(http.StatusOK, rsp)

}

func (server *Server) CheckUserSession(ctx *gin.Context) {
	// Extract the token from the Authorization header
	authorizationHeader := ctx.GetHeader("Authorization")
	token := strings.TrimPrefix(authorizationHeader, "Bearer ")

	// Verify the token
	_, err := server.tokenMaker.VerifyToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"isAuthenticated": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"isAuthenticated": true})
}
