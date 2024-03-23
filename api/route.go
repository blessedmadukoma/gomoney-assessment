package api

import (
	"github.com/gin-gonic/gin"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes(router *gin.Engine, srv *Server) {
	// router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// routes := router.Group("/api", srv.rateLimit())
	routes := router.Group("/api")
	{

		// health check
		routes.GET("/health", healthy)

		// auth routes
		authRoute := routes.Group("/auth")
		{
			// authRoute.GET("/current_user", srv.getCurrentUserBySession)
			authRoute.POST("/login", srv.login)
			authRoute.POST("/register", srv.register)
			authRoute.POST("/refresh", srv.refresh)
			authRoute.POST("/logout", srv.logout)
		}

		// // user routes
		// userRoute := routes.Group("/user")
		// {
		// 	// userRoute.POST("/", srv.CreateUser)
		// 	userRoute.GET("/", srv.GetUsers)
		// 	userRoute.GET("/:id", srv.GetUserByID)
		// 	// userRoute.GET("/:email", srv.GetUserByEmail)

		// 	userRoute.POST("/expense", srv.CreateExpense)
		// 	userRoute.GET("/expense", srv.GetExpenses)
		// }

		// teams routes
		teamsRoute := routes.Group("/teams")
		{
			teamsRoute.POST("/", srv.createTeam, AuthenticatedMiddleware())
			teamsRoute.GET("/", srv.getTeams)
			teamsRoute.GET("/:id", srv.getTeam)
			teamsRoute.PATCH("/:id", srv.editTeam, isAdminMiddleware())
			teamsRoute.DELETE("/:id", srv.removeTeam, isAdminMiddleware())
		}

		// fixtures routes
		fixturesRoute := routes.Group("/fixtures")
		{
			fixturesRoute.POST("/", srv.createFixture, isAdminMiddleware())
			fixturesRoute.GET("/", srv.getFixtures)
			fixturesRoute.GET("/:id", srv.getFixture)
			fixturesRoute.PATCH("/:id", srv.editFixture, isAdminMiddleware())
			fixturesRoute.DELETE("/:id", srv.removeFixture, isAdminMiddleware())
		}
	}
}
