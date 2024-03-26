package api

import (
	"github.com/gin-gonic/gin"
)

func (srv *Server) Routes(router *gin.Engine) {
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

		// teams routes
		teamsRoute := routes.Group("/teams").Use(authMiddleware())
		{
			teamsRoute.POST("/", isAdminMiddleware(srv.Collections), srv.createTeam)
			teamsRoute.GET("/", srv.getTeams)
			teamsRoute.GET("/search", srv.searchTeams)
			teamsRoute.GET("/:id", srv.getTeam)
			teamsRoute.PATCH("/:id", isAdminMiddleware(srv.Collections), srv.editTeam)
			teamsRoute.DELETE("/:id", isAdminMiddleware(srv.Collections), srv.removeTeam)
		}

		// fixtures routes
		fixturesRoute := routes.Group("/fixtures").Use(authMiddleware())
		{
			fixturesRoute.POST("/", isAdminMiddleware(srv.Collections), srv.createFixture)
			fixturesRoute.GET("/", srv.getFixtures)
			fixturesRoute.GET("/:id", srv.getFixtureByID)
			fixturesRoute.GET("/link/:id", srv.getFixtureByLink)
			fixturesRoute.PATCH("/:id", isAdminMiddleware(srv.Collections), srv.editFixture)
			fixturesRoute.DELETE("/:id", isAdminMiddleware(srv.Collections), srv.removeFixture)
		}
	}
}
