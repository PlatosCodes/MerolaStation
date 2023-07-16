package api

import (
	"fmt"

	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/token"
	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	Store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		Store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/:id", server.getUser)
	authRoutes.GET("/users", server.listUser)

	authRoutes.POST("/trains", server.createTrain)
	authRoutes.GET("/trains/:id", server.getTrain)
	authRoutes.GET("/trains/model/:model_number", server.getTrainByModel)

	authRoutes.GET("/trains", server.listTrain)
	authRoutes.PUT("/trains", server.updateTrainValue)

	authRoutes.GET("/users/:id/collection", server.getUserCollection)

	authRoutes.POST("/users/:id/collection/:train_id", server.createCollectionTrain)
	authRoutes.DELETE("/users/:id/collection/:train_id", server.deleteCollectionTrain)

	authRoutes.GET("/users/:id/wishlist", server.getUserWishlist)

	authRoutes.POST("/users/:id/wishlist/:train_id", server.createWishlistTrain)
	authRoutes.DELETE("/users/:id/wishlist/:train_id", server.deleteWishlistTrain)

	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
