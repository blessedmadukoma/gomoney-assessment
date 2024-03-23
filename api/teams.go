package api

import (
	"log"
	"net/http"
	"time"

	db "github.com/blessedmadukoma/gomoney-assessment/db/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (srv *Server) createTeam(ctx *gin.Context) {
	var teamsParams db.CreateTeamsParams

	if err := ctx.ShouldBindJSON(&teamsParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("cannot bind teams data", err))
		return
	}

	teamsParams.CreatedAt = time.Now()
	teamsParams.UpdatedAt = teamsParams.CreatedAt

	filter := bson.D{
		{Key: "teamname", Value: teamsParams.TeamName},
		{Key: "shortname", Value: teamsParams.ShortName},
	}
	count, err := srv.collections["teams"].CountDocuments(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to count records", err))
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusConflict, errorResponse("team already exists", nil))
		return
	}

	res, err := srv.collections["teams"].InsertOne(ctx, teamsParams)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			ctx.JSON(http.StatusBadRequest, errorResponse("team already exists", er))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to create team", err))
		return
	}

	var newTeam *db.TeamsParams
	query := bson.M{"_id": res.InsertedID}

	err = srv.collections["teams"].FindOne(ctx, query).Decode(&newTeam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("unable to find newly created user", err))
		return
	}

	data := gin.H{
		"team": db.ToTeamsResponse(newTeam),
	}

	ctx.JSON(http.StatusCreated, successResponse("Team created successfully", data))
	return
}

func (srv *Server) getTeams(ctx *gin.Context) {
	// Define a filter to match all documents
	filter := bson.D{}

	// Find all teams
	cursor, err := srv.collections["teams"].Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and decode each document into a Team struct
	var teams []db.TeamsParams

	for cursor.Next(ctx) {
		var team db.TeamsParams
		if err := cursor.Decode(&team); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse("failed to decode team values", err))
			return
		}
		teams = append(teams, team)
	}

	// Check for errors during cursor iteration
	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("cursor iteration error", err))
		return
	}
	ctx.JSON(http.StatusOK, successResponse("teams retrieved successfully", teams))
}

func (srv *Server) searchTeams(ctx *gin.Context) {
	// Define a filter to match all documents
	filter := bson.D{}
	ctx.JSON(http.StatusOK, filter)
	return
}

func (srv *Server) getTeam(ctx *gin.Context) {
	objectID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to verify team ID", err))
		return
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
	}

	var team db.TeamsParams
	err = srv.collections["teams"].FindOne(ctx, filter).Decode(&team)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, errorResponse("Team not found", err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse("error retrieving team", err))
			return
		}
	}

	ctx.JSON(http.StatusOK, successResponse("team retrieved successfully", team))
}

func (srv *Server) editTeam(ctx *gin.Context) {
	return
}

func (srv *Server) removeTeam(ctx *gin.Context) {
	return
}
