package token

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTToken struct {
	config *utils.Config
}

var (
	ErrInvalidToken = errors.New("invalid authentication token")
	ErrExpiredToken = errors.New("token has expired")
)

func NewJWTToken(cfg *utils.Config) *JWTToken {
	return &JWTToken{
		config: cfg,
	}
}

type jwtClaim struct {
	jwt.RegisteredClaims
	UserID    primitive.ObjectID `json:"user_id"`
	ExpiredAt int64              `json:"expired_at"`
}

// CreateToken creates a new JWT token
func (j *JWTToken) CreateToken(user_id primitive.ObjectID, duration time.Duration) (string, error) {
	claims := jwtClaim{
		UserID: user_id,
		// ExpiredAt: time.Now().Add(time.Minute * 15).Unix(),
		ExpiredAt: int64(duration),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.config.TokenSymmetricKey))

	if err != nil {
		log.Println("Error signing string:", tokenString)
		return "", err
	}

	return string(tokenString), nil
}

// VerifyToken verifies a JWT token
func (j *JWTToken) VerifyToken(tokenString string) (primitive.ObjectID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaim{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid authentication token")
		}

		return []byte(j.config.TokenSymmetricKey), nil
	})

	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("invalid authentication token")
	}

	claims, ok := token.Claims.(*jwtClaim)

	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("invalid authentication token")
	}

	if claims.ExpiredAt < time.Now().Unix() {
		return primitive.ObjectID{}, fmt.Errorf("token has expired")
	}

	return claims.UserID, nil
}
