package api

import (

	// "database/sql"
	// db "fintrax/db/sqlc"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Role string

const (
	Admin Role = "admin"
	Fan        = "fan"
)

type RegisterParams struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `json:"firstname" binding:"required"`
	LastName  string             `json:"lastname" binding:"required"`
	Email     string             `json:"email" binding:"required,email"`
	Role      string             `json:"role"`
	Password  string             `json:"password" binding:"required"`
}

// func (u UserResponse) toNewUserResponse(user *db.User) *UserResponse {
// 	return &UserResponse{
// 		ID:        user.ID,
// 		Email:     user.Email,
// 		CreatedAt: user.CreatedAt,
// 		UpdatedAt: user.UpdatedAt,
// 	}
// }

func (srv *Server) register(ctx *gin.Context) {
	var user RegisterParams

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("cannot bind user data", err))
		return
	}

	if user.Role != string(Admin) {
		user.Role = string(Fan)
	}

	// hashedPassword, err := utils.HashPassword(user.Password)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// arg := db.CreateUserParams{
	// 	Email:          user.Email,
	// 	HashedPassword: hashedPassword,
	// }

	// newUser, err := a.server.queries.CreateUser(context.Background(), arg)
	// if err != nil {
	// 	if pgErr, ok := err.(*pq.Error); ok {
	// 		// violated unique constraint i.e. user already exists
	// 		if pgErr.Code == "23505" {
	// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 			return
	// 		}
	// 	}

	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// c.JSON(http.StatusCreated, UserResponse{}.toNewUserResponse(&newUser))
	data := gin.H{
		"user": user,
	}

	ctx.JSON(http.StatusCreated, successResponse("user created successfully", data))
}

func (srv *Server) logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
	return
}

type LoginParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (srv *Server) login(ctx *gin.Context) {
	var user LoginParams

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("cannot bind user data", err))
		return
	}

	// dbUser, err := a.server.queries.GetUserByEmail(ctx, user.Email)

	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("user not found: %v", err.Error())})
	// 		return
	// 	}

	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// if err := utils.VerifyPassword(user.Password, dbUser.HashedPassword); err != nil {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("invalid password: %v", err.Error())})
	// 	return
	// }

	dbUserID := 1

	token, err := tokenController.CreateToken(int64(dbUserID), srv.config.AccessTokenDuration)
	// token, err := tokenController.CreateToken(dbUser.ID, srv.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("could not create token: %v", err))
	}

	data := gin.H{
		"token": token,
	}

	ctx.JSON(http.StatusOK, successResponse("login successful", data))
}
