package api

import (
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	// "trackit/token"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/time/rate"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func rateLimit() {}

// isAdminMiddleware checks if the user role is "admin"
func isAdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := getAuthorizationPayload(ctx)

		// check users table to see if the userId has a role of admin

		ctx.Set("user_id", userId)
	}
}

// authMiddleware authorizes/validates a user
// func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
// 		if len(authorizationHeader) == 0 {
// 			err := errors.New("authorization header not provided")
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse("unauthorized", err))
// 			return
// 		}

// 		fields := strings.Fields(authorizationHeader)
// 		if len(fields) < 2 {
// 			err := errors.New("invalid authorization header format")
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse("unauthorized", err))
// 			return
// 		}

// 		authorizationType := strings.ToLower(fields[0])
// 		if authorizationType != authorizationTypeBearer {
// 			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse("unauthorized", err))
// 			return
// 		}

// 		accessToken := fields[1]
// 		payload, err := tokenMaker.VerifyToken(accessToken)
// 		if err != nil {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse("invalid access token", err))
// 			return
// 		}

// 		// store payload in the context
// 		ctx.Set(authorizationPayloadKey, payload)
// 		ctx.Next()
// 	}
// }

func AuthenticatedMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userId := getAuthorizationPayload(ctx)

		ctx.Set("user_id", userId)
	}
}

// getAuthorizationPayload retrieves the authorization payload from the context
func getAuthorizationPayload(ctx *gin.Context) primitive.ObjectID {
	token := ctx.GetHeader("Authorization")

	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
		ctx.Abort()
		return primitive.ObjectID{}
	}

	splitToken := strings.Split(token, " ")

	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authentication token"})
		ctx.Abort()
		return primitive.ObjectID{}
	}

	userId, err := tokenController.VerifyToken(splitToken[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		ctx.Abort()
		return primitive.ObjectID{}
	}

	return userId
}

// setCorsHeaders sets the CORS headers
func setCorsHeaders(corsConfig *cors.Config) {
	corsConfig.AllowOrigins = []string{"https://localhost", "http://localhost", "http://localhost:3000", "https://localhost:3000", "https://trakkit.vercel.app", "http://trakkit.vercel.app"}

	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "Accept", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Allow-Credentials", "Access-Control-Max-Age", "Access-Control-Expose-Headers", "Access-Control-Request-Headers", "Access-Control-Request-Method", "X-Forwarded-For", "X-Forwarded-Host", "X-Forwarded-Port", "X-Forwarded-Proto", "X-Real-Ip", "X-Request-Id", "X-Scheme", "X-Forwarded-Proto", "X-Forwarded-Protocol", "X-Forwarded-Ssl", "X-Url-Scheme", "X-Forwarded-Host", "X-Forwarded-Server", "X-Forwarded-For", "withCredentials"}

	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS", "GET", "POST", "PUT", "DELETE", "PATCH")

	corsConfig.AllowCredentials = true
}

// rateLimit - IP-based rate limiting
func (srv *Server) rateLimit() gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// background goroutine to remove old entries from the clients map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				// check if the client hasn't been seen for the past 3 minutes
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return func(ctx *gin.Context) {
		if srv.config.Limiter.ENABLED {
			ip, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)

			if err != nil {
				log.Fatal("error splitting network address:", err)
				return
			}

			// lock the mutex to prevent concurrent execution
			mu.Lock()

			// check if the IP exists in the map, if it doesn't, initialize a new rate limiter and add the IP address and limiter to the map
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(srv.config.Limiter.RPS), srv.config.Limiter.BURST),
				}
			}

			// update the client's last seen
			clients[ip].lastSeen = time.Now()

			// if the request is not allowed, unlock the mutex and send 429 error
			if !clients[ip].limiter.Allow() {
				// fmt.Println("IP:", ip, "\nLast seen:", clients[ip].lastSeen.String(), "\nTokens:", clients[ip].limiter.Tokens(), "\n...")
				mu.Unlock()

				srv.rateLimitExceededResponse(ctx)
				return
			}
			// Very Important: unlock the mutex before calling the next handler in the chain.
			mu.Unlock()
		}

		ctx.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {

		// log.Println(c.Request.Header)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		// c.Writer.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			log.Println("got options and stopped")
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
