package api

import (
	"github.com/gin-gonic/gin"
)

func (srv *Server) Routes(router *gin.Engine) {
	// func Routes(router *gin.Engine, srv *Server) {
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

		// // user routes: update user profile, delete user profile
		// userRoute := routes.Group("/user").Use(authMiddleware())
		// {
		// 	// userRoute.PATCH("/:id", isAdminMiddleware(srv.collections), srv.CreateUser)
		// 	userRoute.GET("/", isAdminMiddleware(srv.collections), srv.GetUsers)
		// 	userRoute.GET("/:id", isAdminMiddleware(srv.collections), srv.GetUserByID)
		// 	// userRoute.GET("/:email", isAdminMiddleware(srv.collections), srv.GetUserByEmail)
		// 	// userRoute.DELETE("/:id", isAdminMiddleware(srv.collections), srv.DeleteUser)
		// }

		// teams routes
		teamsRoute := routes.Group("/teams").Use(authMiddleware())
		{
			teamsRoute.POST("/", isAdminMiddleware(srv.collections), srv.createTeam)
			teamsRoute.GET("/", srv.getTeams)
			teamsRoute.GET("/search", srv.searchTeams)
			teamsRoute.GET("/:id", srv.getTeam)
			teamsRoute.PATCH("/:id", isAdminMiddleware(srv.collections), srv.editTeam)
			teamsRoute.DELETE("/:id", isAdminMiddleware(srv.collections), srv.removeTeam)
		}

		// fixtures routes
		fixturesRoute := routes.Group("/fixtures").Use(authMiddleware())
		{
			fixturesRoute.POST("/", isAdminMiddleware(srv.collections), srv.createFixture)
			fixturesRoute.GET("/", srv.getFixtures)
			fixturesRoute.GET("/:id", srv.getFixtureByID)
			fixturesRoute.GET("/link/:id", srv.getFixtureByLink)
			fixturesRoute.PATCH("/:id", isAdminMiddleware(srv.collections), srv.editFixture)
			fixturesRoute.DELETE("/:id", isAdminMiddleware(srv.collections), srv.removeFixture)
		}
	}
}
