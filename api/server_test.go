package api

import (
	"os"
	"testing"
	"time"

	// db "trackit/db/sqlc"
	"github.com/blessedmadukoma/gomoney-assessment/utils"

	"github.com/gin-gonic/gin"
	// _ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// func newTestServer(t *testing.T, store db.Store) *Server {
func newTestServer(t *testing.T) *Server {
	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	// server, err := NewServer(config, &store)
	server, err := NewServer(config)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
