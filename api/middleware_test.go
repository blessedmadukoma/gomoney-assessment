package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	// db "trackit/db/sqlc"

	"github.com/blessedmadukoma/gomoney-assessment/token"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// var tokenController *token.JWTToken

func addAuthorization(t *testing.T, request *http.Request, tokenController *token.JWTToken, authorizationType string, userId primitive.ObjectID, duration time.Duration) {
	tokenString, err := tokenController.CreateToken(userId, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, tokenString)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenController *token.JWTToken)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenController *token.JWTToken) {
				addAuthorization(t, request, tokenController, authorizationTypeBearer, primitive.NewObjectID(), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenController *token.JWTToken) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenController *token.JWTToken) {
				addAuthorization(t, request, tokenController, "unsupported", primitive.NewObjectID(), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenController *token.JWTToken) {
				addAuthorization(t, request, tokenController, "", primitive.NewObjectID(), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredAuthorizationToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenController *token.JWTToken) {
				addAuthorization(t, request, tokenController, authorizationTypeBearer, primitive.NewObjectID(), -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// server := newTestServer(t, db.Store{})
			server := newTestServer(t, map[string]*mongo.Collection{})

			authPath := "/api/auth"
			server.router.GET(authPath, AuthenticatedMiddleware(),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, tokenController)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
