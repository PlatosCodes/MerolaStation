package api

import (
	"fmt"

	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/token"
	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our train collector service.
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

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // This should be your frontend's address
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	config.AllowCredentials = true // Important!
	config.ExposeHeaders = []string{"Authorization"}

	router.Use(cors.New(config))

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/renew_access", server.RenewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/:id", server.getUser)
	authRoutes.GET("/users", server.listUser)

	authRoutes.POST("/trains", server.createTrain)
	authRoutes.GET("/trains/:id", server.getTrain)
	authRoutes.GET("/trains/model/:model_number", server.getTrainByModel)

	authRoutes.GET("/trains", server.listTrain)
	authRoutes.PUT("/trains", server.updateTrainValue)

	authRoutes.GET("/users/:id/collection", server.listUserCollection)

	authRoutes.POST("/users/:id/collection/:train_id", server.createCollectionTrain)
	authRoutes.DELETE("/users/:id/collection/:train_id", server.deleteCollectionTrain)

	authRoutes.POST("/trade_offer", server.createTradeOffer)
	authRoutes.GET("/trade_offer/:id", server.getTradeOfferByID)

	authRoutes.GET("/users/trade_offers/offered/:offered_train_owner", server.listUserTradeOffers)
	authRoutes.GET("/users/trade_offers/requests/:requested_train_owner", server.listUserTradeRequests)
	authRoutes.GET("/users/trade_offers/all/:offered_train_owner", server.listAllUserTradeOffers)
	authRoutes.DELETE("/users/trade_offer", server.deleteTradeOffer)

	authRoutes.POST("/trade", server.createTrade)
	authRoutes.GET("/trade/:id", server.getTradeTransaction)

	authRoutes.GET("/collection_trains/trade_offers/:id", server.listCollectionTrainTradeOffers)

	authRoutes.GET("/users/:id/wishlist", server.listUserWishlist)

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
