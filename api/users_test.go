package api

import (
	"testing"

	db "github.com/blessedmadukoma/gomoney-assessment/db/models"
	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (user db.UserParams, password string) {
	password = utils.RandomString(8)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = db.UserParams{
		FirstName: utils.RandomName(),
		LastName:  utils.RandomName(),
		Email:     utils.RandomEmail(),
		Password:  hashedPassword,
	}
	return
}

func randomFanUser(t *testing.T) (user db.RegisterParams, password string) {
	password = utils.RandomString(8)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = db.RegisterParams{
		FirstName: utils.RandomName(),
		LastName:  utils.RandomName(),
		Email:     utils.RandomEmail(),
		Password:  hashedPassword,
	}
	return
}

func randomAdminUser(t *testing.T) (user db.RegisterParams, password string) {
	password = utils.RandomString(8)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = db.RegisterParams{
		FirstName: utils.RandomName(),
		LastName:  utils.RandomName(),
		Email:     utils.RandomEmail(),
		Role:      "admin",
		Password:  hashedPassword,
	}
	return
}
