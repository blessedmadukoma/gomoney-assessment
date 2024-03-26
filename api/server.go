package api

import (
	"fmt"
	"net/http"

	"github.com/blessedmadukoma/gomoney-assessment/token"
	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

var tokenController *token.JWTToken

type Server struct {
	Config      utils.Config
	Collections map[string]*mongo.Collection
	redisclient *redis.Client
	Router      *gin.Engine
}

func healthy(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Healthy")
	return
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config utils.Config, collections map[string]*mongo.Collection, redisClient *redis.Client) (*Server, error) {

	tokenController = token.NewJWTToken(&config)

	server := &Server{
		Collections: collections,
		Config:      config,
		redisclient: redisClient,
	}

	gin.SetMode(config.GinMode)

	router := gin.Default()

	router.Use(CORS())
	router.Use(server.rateLimit())

	router.SetTrustedProxies(nil)
	router.TrustedPlatform = gin.PlatformCloudflare

	server.Routes(router)

	server.Router = router

	return server, nil
}

// StartServer runs the HTTP server on a specific address
func (srv *Server) StartServer(address string) error {
	fmt.Printf("Server starting on address: %s\n", address)
	return srv.Router.Run(fmt.Sprintf(":%s", address))
}
