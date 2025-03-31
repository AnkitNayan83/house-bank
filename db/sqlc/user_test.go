package db

import (
	"context"
	"testing"
	"time"

	"github.com/AnkitNayan83/houseBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomString(7),
		FullName:       util.RandomString(10),
		Email:          util.RandomString(10) + "@gmail.com",
		HashedPassword: util.RandomString(10),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, arg.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByUsername(t *testing.T) {
	user := createRandomUser(t)

	userInDb, err := testQueries.GetUserByUsername(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, userInDb)

	require.Equal(t, user.FullName, userInDb.FullName)
	require.Equal(t, user.Username, userInDb.Username)
	require.Equal(t, user.Email, userInDb.Email)
	require.Equal(t, user.HashedPassword, userInDb.HashedPassword)
	require.WithinDuration(t, user.CreatedAt.Time, userInDb.CreatedAt.Time, time.Second)
}
