package db

import (
	"context"
	"testing"
	"time"

	"github.com/AnkitNayan83/houseBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, arg.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccountById(t *testing.T) {
	account := createRandomAccount(t)

	accountInDb, err := testQueries.GetAccountById(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountInDb)

	require.Equal(t, account.ID, accountInDb.ID)
	require.Equal(t, account.Owner, accountInDb.Owner)
	require.Equal(t, account.Balance, accountInDb.Balance)
	require.Equal(t, account.Currency, accountInDb.Currency)
	require.WithinDuration(t, account.CreatedAt.Time, accountInDb.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	arg := UpdateAccountBalanceParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}

	updatedAccount, err := testQueries.UpdateAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, account.ID, updatedAccount.ID)
	require.Equal(t, account.Owner, updatedAccount.Owner)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
	require.Equal(t, account.Currency, updatedAccount.Currency)
	require.WithinDuration(t, account.CreatedAt.Time, updatedAccount.CreatedAt.Time, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	_, err = testQueries.GetAccountById(context.Background(), account.ID)

	require.Error(t, err)
	require.EqualError(t, err, "no rows in result set")
}

func TestListAccount(t *testing.T) {
	for range 10 {
		createRandomAccount(t)
	}

	arg := GetAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.GetAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
