package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/blessedmadukoma/gomoney-assessment/db/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (srv *Server) teamExists(ctx *gin.Context, teamName string) bool {
	filter := bson.M{"teamname": teamName}
	count, err := srv.collections["teams"].CountDocuments(ctx, filter)
	if err != nil {
		return false
	}
	return count > 0
}

func (srv *Server) getShortName(ctx *gin.Context, teamName string) string {
	var teamInfo db.TeamsParams
	filter := bson.M{"teamname": teamName}
	err := srv.collections["teams"].FindOne(ctx, filter).Decode(&teamInfo)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse("failed to find team", err))
	}

	return teamInfo.ShortName
}

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
	var teams []db.TeamsParams
	filter := bson.D{}

	teamsRedisData, err := srv.GetDataFromRedis(ctx, "teams", &teams)
	if err == nil {
		ctx.JSON(http.StatusOK, successResponse("teams retrieved successfully from redis", teamsRedisData))
		return
	}

	log.Println("Cache Miss - failed to get data from redis:", err)

	// Find all teams from mongodb: cache miss
	cursor, err := srv.collections["teams"].Find(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to find teams", err))
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var team db.TeamsParams
		if err := cursor.Decode(&team); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse("failed to decode team values", err))
			return
		}
		teams = append(teams, team)
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("cursor iteration error", err))
		return
	}

	// store the data in redis
	err = srv.SetDataIntoRedis(ctx, "teams", teams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("failed to set teams data to redis", err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse("teams retrieved successfully", teams))
}

func (srv *Server) searchTeams(ctx *gin.Context) {
	// Example of API endpoint: <<BASE_URL>>/api/teams/search?q=che -> to get `che`lsea or man`che`ster united

	// Note: you can search by team name, shortname or object ID i.e. ID

	searchQuery := ctx.Query("q")

	var filter primitive.M

	teamID, err := primitive.ObjectIDFromHex(searchQuery)
	if err != nil {
		// seach by team or short name
		filter = bson.M{
			"$or": []bson.M{
				{"teamname": primitive.Regex{Pattern: searchQuery, Options: "i"}},
				{"shortname": primitive.Regex{Pattern: searchQuery, Options: "i"}},
			},
		}
	} else {
		// search by Object ID
		filter = bson.M{"_id": teamID}
	}

	cursor, err := srv.collections["teams"].Find(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Error searching for teams", err))
		return
	}
	defer cursor.Close(ctx)

	var teams []db.TeamsParams
	if err := cursor.All(ctx, &teams); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Error retrieving teams", err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse("Teams retrieved successfully", teams))
}

func (srv *Server) getTeam(ctx *gin.Context) {
	var team db.TeamsParams

	objectID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to verify team ID", err))
		return
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
	}

	teamData, err := srv.GetDataFromRedis(ctx, fmt.Sprintf("team-%s", objectID), &team)
	if err == nil {
		ctx.JSON(http.StatusOK, successResponse("team retrieved successfully from redis", teamData))
		return
	}

	log.Println("Cache Miss - failed to get team data from redis:", err)

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

	// store the data in redis
	err = srv.SetDataIntoRedis(ctx, fmt.Sprintf("team-%s", objectID), team)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("failed to set team data to redis", err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse("team retrieved successfully", team))
}

func (srv *Server) editTeam(ctx *gin.Context) {
	objectID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to verify team ID", err))
		return
	}

	var payload bson.M
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("failed to bind payload data", err))
		return
	}

	payload["updated_at"] = time.Now()

	filter := bson.M{"_id": objectID}

	update := bson.M{"$set": payload}

	result, err := srv.collections["teams"].UpdateOne(ctx, filter, update)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("error updating team", err))
		return
	}

	if result.ModifiedCount == 0 {
		ctx.JSON(http.StatusNotFound, errorResponse("team not modified", nil))
		return
	}

	var updatedTeam *db.TeamsParams
	query := bson.M{"_id": objectID}

	err = srv.collections["teams"].FindOne(ctx, query).Decode(&updatedTeam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("failed to find updated team", err))
		return
	}

	data := gin.H{
		"team": db.ToTeamsResponse(updatedTeam),
	}

	ctx.JSON(http.StatusOK, successResponse("team updated successfully", data))
}

func (srv *Server) removeTeam(ctx *gin.Context) {
	teamID := ctx.Param("id")

	objectID, err := primitive.ObjectIDFromHex(teamID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid team ID", err))
		return
	}

	filter := bson.M{"_id": objectID}

	result, err := srv.collections["teams"].DeleteOne(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to delete team", err))
		return
	}

	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusNotFound, errorResponse("team not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, successResponse("team deleted successfully", nil))
}
