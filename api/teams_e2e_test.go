package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blessedmadukoma/gomoney-assessment/db"
	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateTeamsE2E(t *testing.T) {
	ctx := context.Background()

	mongoclient, collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, collections)

	srv, err := NewServer(config, collections, redisclient)
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.router)
	defer ts.Close()

	token := srv.obtainAdminAuthToken(t, ts)

	teamName := utils.RandomName()

	// Sample team data for testing
	teamData := map[string]interface{}{
		"teamName":  teamName,
		"shortName": teamName[:3],
	}

	// Convert team data to JSON
	payload, _ := json.Marshal(teamData)

	// Send POST request to the create-team endpoint
	req, err := http.NewRequest("POST", ts.URL+"/api/teams", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Assert the response status code
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// Decode the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}

	// Assert the response data as needed
	// Example: assert.Equal(t, "Team created successfully", responseBody["message"])
}

func TestGetTeams(t *testing.T) {
	ctx := context.Background()

	mongoclient, collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, collections)

	srv, err := NewServer(config, collections, redisclient)
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.router)
	defer ts.Close()

	token := srv.obtainFanAuthToken(t, ts)

	// Sample team data for testing
	// teamData := map[string]interface{}{
	// 	"teamName":  "Test Team",
	// 	"shortName": "TT",
	// }

	// Convert team data to JSON
	// payload, _ := json.Marshal(teamData)

	// Send POST request to the create-team endpoint
	req, err := http.NewRequest("GET", ts.URL+"/api/teams", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Assert the response status code
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Decode the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}

	// Assert the response data as needed
	// Example: assert.Equal(t, "Team created successfully", responseBody["message"])
}
