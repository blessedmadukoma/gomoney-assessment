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
	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/stretchr/testify/assert"
)

func (srv *Server) createTestTeamsData(t *testing.T, ts *httptest.Server) *http.Response {
	token := srv.obtainAdminAuthToken(t, ts)

	teamName := utils.RandomName()

	teamData := map[string]interface{}{
		"teamName":  teamName,
		"shortName": teamName[:3],
	}

	payload, _ := json.Marshal(teamData)

	req, err := http.NewRequest("POST", ts.URL+"/api/teams", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	return res
}

func TestCreateTeamsE2E(t *testing.T) {
	ctx := context.Background()

	mongoclient, collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, collections)

	srv, err := NewServer(config, collections, redisclient)
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.router)
	defer ts.Close()

	res := srv.createTestTeamsData(t, ts)

	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}
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

	req, err := http.NewRequest("GET", ts.URL+"/api/teams", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTeamByID(t *testing.T) {
	ctx := context.Background()

	mongoclient, collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, collections)

	srv, err := NewServer(config, collections, redisclient)
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.router)
	defer ts.Close()

	token := srv.obtainFanAuthToken(t, ts)

	response := srv.createTestTeamsData(t, ts)
	defer response.Body.Close()

	assert.Equal(t, http.StatusCreated, response.StatusCode)

	var teamResponseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&teamResponseBody)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("teamResponseBody:", teamResponseBody)

	teamID := teamResponseBody["data"].(map[string]interface{})["team"].(map[string]interface{})["id"].(string)

	log.Println("teamID:", teamID)

	req, err := http.NewRequest("GET", ts.URL+"/api/teams/"+teamID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveTeam(t *testing.T) {
	ctx := context.Background()

	mongoclient, collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, collections)

	srv, err := NewServer(config, collections, redisclient)
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.router)
	defer ts.Close()

	token := srv.obtainAdminAuthToken(t, ts)

	response := srv.createTestTeamsData(t, ts)
	defer response.Body.Close()

	assert.Equal(t, http.StatusCreated, response.StatusCode)

	var teamResponseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&teamResponseBody)
	if err != nil {
		t.Fatal(err)
	}

	teamID := teamResponseBody["data"].(map[string]interface{})["team"].(map[string]interface{})["id"].(string)

	req, err := http.NewRequest("DELETE", ts.URL+"/api/teams/"+teamID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}
}
