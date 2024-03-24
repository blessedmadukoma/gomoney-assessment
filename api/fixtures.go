package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	db "github.com/blessedmadukoma/gomoney-assessment/db/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (srv *Server) generateFixture(ctx *gin.Context, home, away string) string {
	return fmt.Sprintf("%s-vs-%s", strings.ToLower(srv.getShortName(ctx, home)), strings.ToLower(srv.getShortName(ctx, away)))
}

func (srv *Server) generateFixtureLink(ctx *gin.Context, home, away string) string {
	shortUUID := xid.New().String()
	return fmt.Sprintf("%s-vs-%s-%s", strings.ToLower(srv.getShortName(ctx, home)), strings.ToLower(srv.getShortName(ctx, away)), shortUUID)
}

func (srv *Server) createFixture(ctx *gin.Context) {
	var createParams db.CreateFixturesParams
	if err := ctx.ShouldBindJSON(&createParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("failed to bind fixture data", err))
		return
	}

	// Check if both teams exist
	missingTeams := make([]string, 0)
	if !srv.teamExists(ctx, createParams.Home) {
		missingTeams = append(missingTeams, createParams.Home)
	}
	if !srv.teamExists(ctx, createParams.Away) {
		missingTeams = append(missingTeams, createParams.Away)
	}

	if len(missingTeams) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("the following team(s) do not exist: %v", missingTeams),
		})
		return
	}

	// Generate fixture and link
	createParams.Link = srv.generateFixtureLink(ctx, createParams.Home, createParams.Away)

	createParams.Fixture = srv.generateFixture(ctx, createParams.Home, createParams.Away)

	if createParams.Status == "" {
		createParams.Status = "pending"
	}

	now := time.Now()
	createParams.CreatedAt = now
	createParams.UpdatedAt = now

	// Check if the fixture already exists
	filter := bson.M{
		"home": createParams.Home,
		"away": createParams.Away,
		// "link":    createParams.Link,  // link is unique, therefore, the fixture will always be unique
		"status":  createParams.Status,
		"fixture": createParams.Fixture,
	}
	count, err := srv.collections["fixtures"].CountDocuments(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to count fixture records", err))
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusConflict, errorResponse("fixture already exists", nil))
		return
	}

	// Insert the fixture into the database
	res, err := srv.collections["fixtures"].InsertOne(ctx, createParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to create fixture", err))
		return
	}

	// Fetch the newly created fixture
	var newFixture db.FixturesParams
	err = srv.collections["fixtures"].FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&newFixture)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to retrieve newly created fixture", err))
		return
	}

	// Prepare response data
	data := gin.H{
		"fixture": db.ToFixturesResponse(&newFixture),
	}

	// Return success response
	ctx.JSON(http.StatusCreated, successResponse("fixture created successfully", data))
}

func (srv *Server) getFixtures(ctx *gin.Context) {
	var fixtures []db.FixturesParams

	// get data from redis if not expired after 5 minutes: cache hit
	fixturesData, err := srv.GetDataFromRedis(ctx, "fixtures", &fixtures)
	if err == nil {
		ctx.JSON(http.StatusOK, successResponse("fixtures retrieved successfully from redis", fixturesData))
		return
	}

	log.Println("Cache Miss - failed to get fixtures data from redis:", err)

	filter := bson.D{}

	cursor, err := srv.collections["fixtures"].Find(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to find all fixtures", err))
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var fixture db.FixturesParams
		if err := cursor.Decode(&fixture); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse("failed to decode fixture values", err))
			return
		}
		fixtures = append(fixtures, fixture)
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("cursor iteration error", err))
		return
	}

	// store the data in redis
	err = srv.SetDataIntoRedis(ctx, "fixtures", fixtures)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("failed to set fixtures data to redis", err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse("fixtures retrieved successfully", fixtures))
}

func (srv *Server) getFixtureByLink(ctx *gin.Context) {
	var fixture db.FixturesParams

	fixtureLink := ctx.Param("id")

	filter := bson.D{
		{Key: "link", Value: fixtureLink},
	}

	// get data from redis if not expired after 5 minutes: cache hit
	fixtureData, err := srv.GetDataFromRedis(ctx, fmt.Sprintf("fixture-link-%s", fixtureLink), &fixture)
	if err == nil {
		ctx.JSON(http.StatusOK, successResponse("fixture by link retrieved successfully from redis", fixtureData))
		return
	}

	log.Println("Cache Miss - failed to get fixture-by-link data from redis:", err)

	err = srv.collections["fixtures"].FindOne(ctx, filter).Decode(&fixture)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, errorResponse("fixture not found", err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse("error retrieving fixture", err))
			return
		}
	}

	// store the data in redis
	err = srv.SetDataIntoRedis(ctx, fmt.Sprintf("fixture-link-%s", fixtureLink), fixture)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("failed to set team data to redis", err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse("fixture retrieved successfully", fixture))
}

func (srv *Server) getFixtureByID(ctx *gin.Context) {
	var fixture db.FixturesParams

	objectID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to verify fixture ID", err))
		return
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
	}

	// get data from redis if not expired after 5 minutes: cache hit
	fixtureData, err := srv.GetDataFromRedis(ctx, fmt.Sprintf("team-%s", objectID), &fixture)
	if err == nil {
		ctx.JSON(http.StatusOK, successResponse("fixture retrieved successfully from redis", fixtureData))
		return
	}

	log.Println("Cache Miss - failed to get fixture data from redis:", err)

	err = srv.collections["fixtures"].FindOne(ctx, filter).Decode(&fixture)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, errorResponse("fixture not found", err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse("error retrieving fixture", err))
			return
		}
	}

	// store the data in redis
	err = srv.SetDataIntoRedis(ctx, fmt.Sprintf("fixture-%s", objectID), fixture)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("failed to set fixture data to redis", err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse("fixture retrieved successfully", fixture))
}

func (srv *Server) editFixture(ctx *gin.Context) {
	objectID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to verify fixture ID", err))
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

	result, err := srv.collections["fixtures"].UpdateOne(ctx, filter, update)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("error updating fixture", err))
		return
	}

	if result.ModifiedCount == 0 {
		ctx.JSON(http.StatusNotFound, errorResponse("fixture not modified", nil))
		return
	}

	var updatedFixture *db.FixturesParams
	query := bson.M{"_id": objectID}

	err = srv.collections["fixtures"].FindOne(ctx, query).Decode(&updatedFixture)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("failed to find updated fixture", err))
		return
	}

	data := gin.H{
		"fixture": db.ToFixturesResponse(updatedFixture),
	}

	// Return success response
	ctx.JSON(http.StatusOK, successResponse("fixture updated successfully", data))
}

func (srv *Server) removeFixture(ctx *gin.Context) {
	fixtureID := ctx.Param("id")

	objectID, err := primitive.ObjectIDFromHex(fixtureID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid fixture ID", err))
		return
	}

	filter := bson.M{"_id": objectID}

	result, err := srv.collections["fixtures"].DeleteOne(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to delete fixture", err))
		return
	}

	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusNotFound, errorResponse("fixture not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, successResponse("fixture deleted successfully", nil))
}
