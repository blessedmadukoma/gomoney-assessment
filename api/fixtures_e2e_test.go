package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blessedmadukoma/gomoney-assessment/db"
	models "github.com/blessedmadukoma/gomoney-assessment/db/models"
	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateFixturesE2E(t *testing.T) {
	ctx := context.Background()

	mongoclient, collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, collections)

	srv, err := NewServer(config, collections, redisclient)
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.router)
	defer ts.Close()

	token := srv.obtainAdminAuthToken(t, ts)

	home := utils.RandomName()
	away := utils.RandomName()

	homeParams := models.CreateTeamsParams{
		TeamName:  home,
		ShortName: home[:3],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = srv.collections["teams"].InsertOne(ctx, homeParams)
	assert.NoError(t, err)

	awayParams := models.CreateTeamsParams{
		TeamName:  away,
		ShortName: away[:3],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = srv.collections["teams"].InsertOne(ctx, awayParams)
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

	mongoclient, collections := db.ConnectMongoDB(ctx, config)

	defer mongoclient.Disconnect(ctx)

	assert.NotNil(t, collections)

	srv, err := NewServer(config, collections, redisclient)
	assert.NoError(t, err)

	ts := httptest.NewServer(srv.router)
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
