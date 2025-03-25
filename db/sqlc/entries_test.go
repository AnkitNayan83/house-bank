package db

import (
	"context"
	"testing"
	"time"

	"github.com/AnkitNayan83/houseBank/util"
	"github.com/stretchr/testify/require"
)

const (
	account1Id = 1
	account2Id = 2
)

func createRandomEntry(t *testing.T) Entry {
	arg := CreateEntryParams{
		AccountID: account1Id,
		Amount:    util.RandomInt(100, 10000),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntryById(t *testing.T) {
	entry := createRandomEntry(t)
	entryInDb, err := testQueries.GetEntryById(context.Background(), entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entryInDb)

	require.Equal(t, entry.ID, entryInDb.ID)
	require.Equal(t, entry.AccountID, entryInDb.AccountID)
	require.Equal(t, entry.Amount, entryInDb.Amount)
	require.WithinDuration(t, entry.CreatedAt.Time, entryInDb.CreatedAt.Time, time.Second)
}

func TestGetEntriesByAccountId(t *testing.T) {
	arg := GetEntriesByAccountIdParams{
		AccountID: account1Id,
		Limit:     5,
		Offset:    0,
	}
	entries, err := testQueries.GetEntriesByAccountId(context.Background(), arg)

	require.NoError(t, err)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}

}
