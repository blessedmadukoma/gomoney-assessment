package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blessedmadukoma/gomoney-assessment/db"
	models "github.com/blessedmadukoma/gomoney-assessment/db/models"
	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/stretchr/testify/assert"
)

func (srv *Server) createTestFixtureData(t *testing.T, ts *httptest.Server) *http.Response {
	ctx := context.Background()

	token := srv.obtainAdminAuthToken(t, ts)

	home := utils.RandomName()
	away := utils.RandomName()

	homeParams := models.CreateTeamsParams{
		TeamName:  home,
		ShortName: home[:3],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := srv.Collections["teams"].InsertOne(ctx, homeParams)
	assert.NoError(t, err)

	awayParams := models.CreateTeamsParams{
		TeamName:  away,
		ShortName: away[:3],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = srv.Collections["teams"].InsertOne(ctx, awayParams)
	assert.NoError(t, err)

	fixtureData := map[string]interface{}{
		"home":   home,
		"away":   away,
		"status": "pending",
	}

	payload, _ := json.Marshal(fixtureData)

	// Send POST request to the create-fixture endpoint
	req, err := http.NewRequest("POST", ts.URL+"/api/fixtures", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	return res
}

func TestCreateFixturesE2E(t *testing.T) {
	ctx := context.Background()

	mongoclient, Collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, Collections)

	srv, err := NewServer(config, Collections, redisclient)
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.Router)
	defer ts.Close()

	res := srv.createTestFixtureData(t, ts)
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFixtures(t *testing.T) {
	ctx := context.Background()

	mongoclient, Collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, Collections)

	srv, err := NewServer(config, Collections, redisclient)
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.Router)
	defer ts.Close()

	token := srv.obtainFanAuthToken(t, ts)

	req, err := http.NewRequest("GET", ts.URL+"/api/fixtures", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Decode the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFixtureByLink(t *testing.T) {
	ctx := context.Background()

	mongoclient, Collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, Collections)

	srv, err := NewServer(config, Collections, redisclient)
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.Router)
	defer ts.Close()

	token := srv.obtainFanAuthToken(t, ts)

	response := srv.createTestFixtureData(t, ts)
	defer response.Body.Close()

	assert.Equal(t, http.StatusCreated, response.StatusCode)

	var fixtureResponseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&fixtureResponseBody)
	if err != nil {
		t.Fatal(err)
	}

	fixtureLink := fixtureResponseBody["data"].(map[string]interface{})["fixture"].(map[string]interface{})["link"].(string)

	req, err := http.NewRequest("GET", ts.URL+"/api/fixtures/link/"+fixtureLink, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Decode the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFixtureByID(t *testing.T) {
	ctx := context.Background()

	mongoclient, Collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, Collections)

	srv, err := NewServer(config, Collections, redisclient)
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.Router)
	defer ts.Close()

	token := srv.obtainFanAuthToken(t, ts)

	response := srv.createTestFixtureData(t, ts)
	defer response.Body.Close()

	assert.Equal(t, http.StatusCreated, response.StatusCode)

	var fixtureResponseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&fixtureResponseBody)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("fixtureResponseBody:", fixtureResponseBody)

	fixtureID := fixtureResponseBody["data"].(map[string]interface{})["fixture"].(map[string]interface{})["id"].(string)

	log.Println("fixtureID:", fixtureID)

	req, err := http.NewRequest("GET", ts.URL+"/api/fixtures/"+fixtureID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Decode the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveFixture(t *testing.T) {
	ctx := context.Background()

	mongoclient, Collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, Collections)

	srv, err := NewServer(config, Collections, redisclient)
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.Router)
	defer ts.Close()

	token := srv.obtainAdminAuthToken(t, ts)

	response := srv.createTestFixtureData(t, ts)
	defer response.Body.Close()

	assert.Equal(t, http.StatusCreated, response.StatusCode)

	var fixtureResponseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&fixtureResponseBody)
	if err != nil {
		t.Fatal(err)
	}

	fixtureID := fixtureResponseBody["data"].(map[string]interface{})["fixture"].(map[string]interface{})["id"].(string)

	req, err := http.NewRequest("DELETE", ts.URL+"/api/fixtures/"+fixtureID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Decode the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}
}
