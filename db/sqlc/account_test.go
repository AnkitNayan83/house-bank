package db

import (
	"context"
	"testing"
	"time"

	"github.com/AnkitNayan83/houseBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
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
	require.WithinDuration(t, account.CreatedAt, accountInDb.CreatedAt, time.Second)
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
	require.WithinDuration(t, account.CreatedAt, updatedAccount.CreatedAt, time.Second)

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
	var owner string
	for range 10 {
		account := createRandomAccount(t)
		owner = account.Owner
	}

	arg := GetAccountsParams{
		Owner:  owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.GetAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, owner, account.Owner)
	}
}
