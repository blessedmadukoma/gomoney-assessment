package api

import (

	// "database/sql"
	// db "fintrax/db/sqlc"

	"net/http"
	"strings"
	"time"

	db "github.com/blessedmadukoma/gomoney-assessment/db/models"
	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Role string

const (
	Admin Role = "admin"
	Fan        = "fan"
)

func (srv *Server) register(ctx *gin.Context) {
	var user db.RegisterParams

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("cannot bind user data", err))
		return
	}

	if user.Role != string(Admin) {
		user.Role = string(Fan)
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("could not hash password", err))
		return
	}

	user.Password = hashedPassword

	res, err := srv.collections["users"].InsertOne(ctx, &user)

	if err != nil {

		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			ctx.JSON(http.StatusBadRequest, errorResponse("user already exists", er))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse("could not create user", err))
		return
	}

	// Create a unique index for the email field
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

	if _, err := srv.collections["users"].Indexes().CreateOne(ctx, index); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("could not create index for email", err))
	}

	var newUser *db.UserParams
	query := bson.M{"_id": res.InsertedID}

	err = srv.collections["users"].FindOne(ctx, query).Decode(&newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("unable to find newly created user", err))
		return
	}

	data := gin.H{
		"user": db.ToUserResponse(newUser),
	}

	ctx.JSON(http.StatusCreated, successResponse("user created successfully", data))
}

type LoginParams struct {
	Email    string `json:"email" bson:"email" binding:"required,email"`
	Password string `json:"password" bson:"password" binding:"required"`
}

func (srv *Server) login(ctx *gin.Context) {
	var user LoginParams

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("cannot bind user data", err))
		return
	}
	dbUser, err := srv.FindUserByEmail(ctx, user.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusBadRequest, errorResponse("no user record found", err))
			return
		}

		ctx.JSON(http.StatusBadRequest, errorResponse("", err))
		return
	}

	if err := utils.VerifyPassword(user.Password, dbUser.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid password", err))
		return
	}

	// Generate Tokens
	access_token, err := tokenController.CreateToken(dbUser.ID, srv.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refresh_token, err := tokenController.CreateToken(dbUser.ID, srv.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("could not create token", err))
		return
	}

	data := gin.H{
		"access_token":  access_token,
		"refresh_token": refresh_token,
		"user":          db.ToUserResponse(dbUser),
	}

	ctx.JSON(http.StatusOK, successResponse("login successful", data))
}

func (srv *Server) logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
	return
}

func (srv *Server) refresh(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
	return
}
