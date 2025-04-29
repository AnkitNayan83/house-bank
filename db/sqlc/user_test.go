package db

import (
	"context"
	"testing"
	"time"

	"github.com/AnkitNayan83/houseBank/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomString(7),
		FullName:       util.RandomString(10),
		Email:          util.RandomString(10) + "@gmail.com",
		HashedPassword: hashedPassword,
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
	require.WithinDuration(t, user.CreatedAt, userInDb.CreatedAt, time.Second)
}

func TestGetUserByEmail(t *testing.T) {
	user := createRandomUser(t)

	userInDb, err := testQueries.GetUserByEmail(context.Background(), user.Email)

	require.NoError(t, err)
	require.NotEmpty(t, userInDb)

	require.Equal(t, user.FullName, userInDb.FullName)
	require.Equal(t, user.Username, userInDb.Username)
	require.Equal(t, user.Email, userInDb.Email)
	require.Equal(t, user.HashedPassword, userInDb.HashedPassword)
	require.WithinDuration(t, user.CreatedAt, userInDb.CreatedAt, time.Second)
}

func TestChangeUserPassword(t *testing.T) {
	user := createRandomUser(t)

	arg := ChangePasswordParams{
		Username:          user.Username,
		HashedPassword:    util.RandomString(10),
		PasswordChangedAt: time.Now(),
	}

	changedUser, err := testQueries.ChangePassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, changedUser)

	require.Equal(t, user.FullName, changedUser.FullName)
	require.Equal(t, user.Username, changedUser.Username)
	require.Equal(t, user.Email, changedUser.Email)
	require.Equal(t, arg.HashedPassword, changedUser.HashedPassword)
	require.WithinDuration(t, user.CreatedAt, changedUser.CreatedAt, time.Second)
}

func TestUpdateUserEmailVerification(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUserEmailVerificationParams{
		Username:        user.Username,
		EmailVerifiedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	changedUser, err := testQueries.UpdateUserEmailVerification(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, changedUser)
	require.NotEmpty(t, changedUser.EmailVerifiedAt)

	require.Equal(t, user.FullName, changedUser.FullName)
	require.Equal(t, user.Username, changedUser.Username)
	require.Equal(t, user.Email, changedUser.Email)
	require.Equal(t, user.HashedPassword, changedUser.HashedPassword)
	require.WithinDuration(t, user.CreatedAt, changedUser.CreatedAt, time.Second)
}
