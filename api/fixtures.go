package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	db "github.com/blessedmadukoma/gomoney-assessment/db/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
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
	filter := bson.D{}

	cursor, err := srv.collections["fixtures"].Find(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to find all fixtures", err))
		return
	}
	defer cursor.Close(ctx)

	var fixtures []db.FixturesParams

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
	ctx.JSON(http.StatusOK, successResponse("fixtures retrieved successfully", fixtures))
}

// func (srv *Server) searchFixtures(ctx *gin.Context) {
// 	// Example of API endpoint: <<BASE_URL>>/api/fixtures/search?q=che -> to get `che`lsea or man`che`ster united

// 	// Note: you can search by fixture name, shortname or object ID i.e. ID

// 	searchQuery := ctx.Query("q")

// 	var filter primitive.M

// 	fixtureID, err := primitive.ObjectIDFromHex(searchQuery)
// 	if err != nil {
// 		// seach by fixture or short name
// 		filter = bson.M{
// 			"$or": []bson.M{
// 				{"fixturename": primitive.Regex{Pattern: searchQuery, Options: "i"}},
// 				{"shortname": primitive.Regex{Pattern: searchQuery, Options: "i"}},
// 			},
// 		}
// 	} else {
// 		// search by Object ID
// 		filter = bson.M{"_id": fixtureID}
// 	}

// 	cursor, err := srv.collections["fixtures"].Find(ctx, filter)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse("Error searching for fixtures", err))
// 		return
// 	}
// 	defer cursor.Close(ctx)

// 	var fixtures []db.FixturesParams
// 	if err := cursor.All(ctx, &fixtures); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse("Error retrieving fixtures", err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, successResponse("Fixtures retrieved successfully", fixtures))
// }

// func (srv *Server) getFixture(ctx *gin.Context) {
// 	objectID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to verify fixture ID", err))
// 		return
// 	}

// 	filter := bson.D{
// 		{Key: "_id", Value: objectID},
// 	}

// 	var fixture db.FixturesParams
// 	err = srv.collections["fixtures"].FindOne(ctx, filter).Decode(&fixture)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			ctx.JSON(http.StatusNotFound, errorResponse("Fixture not found", err))
// 			return
// 		} else {
// 			ctx.JSON(http.StatusInternalServerError, errorResponse("error retrieving fixture", err))
// 			return
// 		}
// 	}

// 	ctx.JSON(http.StatusOK, successResponse("fixture retrieved successfully", fixture))
// }

// func (srv *Server) getFixtureByLink(ctx *gin.Context) {
// 	return
// }

// func (srv *Server) editFixture(ctx *gin.Context) {
// 	objectID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to verify fixture ID", err))
// 		return
// 	}

// 	var payload bson.M
// 	if err := ctx.ShouldBindJSON(&payload); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse("failed to bind payload data", err))
// 		return
// 	}

// 	payload["updated_at"] = time.Now()

// 	filter := bson.M{"_id": objectID}

// 	update := bson.M{"$set": payload}

// 	result, err := srv.collections["fixtures"].UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse("error updating fixture", err))
// 		return
// 	}

// 	if result.ModifiedCount == 0 {
// 		ctx.JSON(http.StatusNotFound, errorResponse("fixture not modified", nil))
// 		return
// 	}

// 	var updatedFixture *db.FixturesParams
// 	query := bson.M{"_id": objectID}

// 	err = srv.collections["fixtures"].FindOne(ctx, query).Decode(&updatedFixture)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse("failed to find updated fixture", err))
// 		return
// 	}

// 	data := gin.H{
// 		"fixture": db.ToFixturesResponse(updatedFixture),
// 	}

// 	// Return success response
// 	ctx.JSON(http.StatusOK, successResponse("fixture updated successfully", data))
// }

// func (srv *Server) removeFixture(ctx *gin.Context) {
// 	fixtureID := ctx.Param("id")

// 	objectID, err := primitive.ObjectIDFromHex(fixtureID)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse("invalid fixture ID", err))
// 		return
// 	}

// 	filter := bson.M{"_id": objectID}

// 	result, err := srv.collections["fixtures"].DeleteOne(ctx, filter)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to delete fixture", err))
// 		return
// 	}

// 	if result.DeletedCount == 0 {
// 		ctx.JSON(http.StatusNotFound, errorResponse("fixture not found", nil))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, successResponse("fixture deleted successfully", nil))
// }
