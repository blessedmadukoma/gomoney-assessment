package api

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine, srv *Server) {
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
		// }

		// teams routes
		teamsRoute := routes.Group("/teams").Use(authMiddleware(*tokenController))
		{
			teamsRoute.POST("/", isAdminMiddleware(srv.collections), srv.createTeam)
			teamsRoute.GET("/", srv.getTeams)
			teamsRoute.GET("/:id", srv.getTeam)
			teamsRoute.PATCH("/:id", isAdminMiddleware(srv.collections), srv.editTeam)
			teamsRoute.DELETE("/:id", isAdminMiddleware(srv.collections), srv.removeTeam)
		}

		// fixtures routes
		fixturesRoute := routes.Group("/fixtures").Use(authMiddleware(*tokenController))
		{
			fixturesRoute.POST("/", isAdminMiddleware(srv.collections), srv.createFixture)
			fixturesRoute.GET("/", srv.getFixtures)
			fixturesRoute.GET("/:id", srv.getFixture)
			fixturesRoute.PATCH("/:id", isAdminMiddleware(srv.collections), srv.editFixture)
			fixturesRoute.DELETE("/:id", isAdminMiddleware(srv.collections), srv.removeFixture)
		}
	}
}
