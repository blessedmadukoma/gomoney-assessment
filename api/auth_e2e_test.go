package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blessedmadukoma/gomoney-assessment/db"
	models "github.com/blessedmadukoma/gomoney-assessment/db/models"
	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	envConfig   = utils.LoadEnvConfig("../.env")
	collections = make(map[string]*mongo.Collection)
	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisDBSource,
	})
)

func TestLoginE2E(t *testing.T) {

	ctx := context.Background()

	mongoclient, collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, collections)

	srv, err := NewServer(config, collections, redisclient)

	assert.NoError(t, err, "unable to create new server")

	user, password := randomUser(t)
	_, err = collections["users"].InsertOne(ctx, user)
	assert.NoError(t, err)

	testCases := []struct {
		name     string
		payload  map[string]interface{}
		expected int
	}{
		{
			name: "OK",
			payload: map[string]interface{}{
				"email":    user.Email,
				"password": password,
			},
			expected: http.StatusOK,
		},
		{
			name: "UserNotFound",
			payload: map[string]interface{}{
				"email":    "invalidemail@gmail.com",
				"password": password,
			},
			expected: http.StatusNotFound,
		},
		{
			name: "Invalid Email",
			payload: map[string]interface{}{
				"email":    "invalidemail",
				"password": password,
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Invalid Password",
			payload: map[string]interface{}{
				"email":    user.Email,
				"password": "sh",
			},
			expected: http.StatusUnauthorized,
		},
		{
			name:     "No Payload",
			payload:  map[string]interface{}{},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			srv.router.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expected, recorder.Code)
		})
	}

	filter := bson.M{"_id": user.ID}

	_, err = srv.collections["users"].DeleteOne(ctx, filter)
	assert.NoError(t, err)
}

func TestRegisterEnd2End(t *testing.T) {
	ctx := context.Background()

	mongoclient, collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, collections)

	type testCase struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectedBody   string
	}

	server, err := NewServer(envConfig, collections, redisclient)
	if err != nil {
		log.Fatal("failed to connect to server", err)
	}

	user, password := randomUser(t)

	testCases := []testCase{
		{
			name: "Register - Created",
			payload: map[string]interface{}{
				"firstname": user.FirstName,
				"lastname":  user.LastName,
				"email":     utils.RandomEmail(),
				// "role":      "admin",
				"password": password,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Register - InvalidEmail",
			payload:        map[string]interface{}{"email": "invalidemail", "password": "password123"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Register - InvalidPassword",
			payload:        map[string]interface{}{"email": "test@example.com", "password": "sh"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Register - NoRequestBody",
			payload:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		// Add more test cases for login if needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal payload to JSON
			payload, _ := json.Marshal(tc.payload)

			// Create a request
			req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewReader(payload))
			assert.NoError(t, err)

			// Set request headers
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder to record the response
			recorder := httptest.NewRecorder()

			// Serve the request
			server.router.ServeHTTP(recorder, req)

			// Check the response status code
			assert.Equal(t, tc.expectedStatus, recorder.Code)

			// Optionally, check the response body
			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, recorder.Body.String())
			}
		})
	}
}

func (srv *Server) obtainFanAuthToken(t *testing.T, ts *httptest.Server) string {
	user, password := randomFanUser(t)

	res, err := srv.collections["users"].InsertOne(context.Background(), &user)
	if err != nil {
		t.Fatal("error inserting into db:", err)
	}

	var newUser *models.UserParams
	query := bson.M{"_id": res.InsertedID}
	err = srv.collections["users"].FindOne(context.Background(), query).Decode(&newUser)
	if err != nil {
		t.Fatal("error retreiving record:", err)
	}

	loginRequest := map[string]interface{}{
		"email":    user.Email,
		"password": password,
	}

	loginRequestBody, err := json.Marshal(&loginRequest)
	if err != nil {
		t.Fatal("error marshalling user json:", err)
	}

	req, err := http.Post(ts.URL+"/api/auth/login", "application/json", bytes.NewReader(loginRequestBody))
	if err != nil {
		t.Fatal(err)
	}
	defer req.Body.Close()

	var responseBody map[string]interface{}
	err = json.NewDecoder(req.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}

	accessToken, ok := responseBody["data"].(map[string]interface{})["access_token"].(string)
	if !ok {
		log.Fatal("failed to retrieve access token")
		return ""
	}

	return accessToken

}
func (srv *Server) obtainAdminAuthToken(t *testing.T, ts *httptest.Server) string {
	user, password := randomAdminUser(t)

	res, err := srv.collections["users"].InsertOne(context.Background(), &user)
	if err != nil {
		t.Fatal("error inserting into db:", err)
	}

	var newUser *models.UserParams
	query := bson.M{"_id": res.InsertedID}
	err = srv.collections["users"].FindOne(context.Background(), query).Decode(&newUser)
	if err != nil {
		t.Fatal("error retreiving record:", err)
	}

	loginRequest := map[string]interface{}{
		"email":    user.Email,
		"password": password,
	}

	loginRequestBody, err := json.Marshal(&loginRequest)
	if err != nil {
		t.Fatal("error marshalling user json:", err)
	}

	req, err := http.Post(ts.URL+"/api/auth/login", "application/json", bytes.NewReader(loginRequestBody))
	if err != nil {
		t.Fatal(err)
	}
	defer req.Body.Close()

	var responseBody map[string]interface{}
	err = json.NewDecoder(req.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}

	accessToken, ok := responseBody["data"].(map[string]interface{})["access_token"].(string)
	if !ok {
		log.Fatal("failed to retrieve access token")
		return ""
	}

	return accessToken

}
