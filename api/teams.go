package api

import (
	"net/http"

	db "github.com/blessedmadukoma/gomoney-assessment/db/models"
	"github.com/gin-gonic/gin"
)

func (srv *Server) createTeam(ctx *gin.Context) {
	var teamsParams db.TeamsParams

	if err := ctx.ShouldBindJSON(&teamsParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("cannot bind teams data", err))
		return
	}

	ctx.JSON(http.StatusCreated, teamsParams)
}

func (srv *Server) getTeams(ctx *gin.Context) {
	return
}

func (srv *Server) searchTeams(ctx *gin.Context) {
	return
}

func (srv *Server) getTeam(ctx *gin.Context) {
	return
}

func (srv *Server) editTeam(ctx *gin.Context) {
	return
}

func (srv *Server) removeTeam(ctx *gin.Context) {
	return
}
