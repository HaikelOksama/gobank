package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/haikeloksama/gobank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	fetchRandName, err := util.RandomOwner()

	require.NoError(t, err)

	arg := CreateAccountParams{
		Owner:    fetchRandName,
		Balance:  util.RandomInt(1, 1000),
		Currency: util.RandomCurrency(),
	}
	defer fmt.Printf("arg: %v\n", arg)
	account, err := testQueries.CreateAccount(context.Background(), arg)
	defer fmt.Printf("account: %v\n", account)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Balance, account.Balance)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	actualAccount := createRandomAccount(t)
	getAccount, err := testQueries.GetAccount(context.Background(), actualAccount.ID)

	require.NoError(t, err)

	require.NotEmpty(t, getAccount)
	require.Equal(t, actualAccount.ID, getAccount.ID)
	require.Equal(t, actualAccount.Owner, getAccount.Owner)
	require.Equal(t, actualAccount.Balance, getAccount.Balance)
	require.Equal(t, actualAccount.Currency, getAccount.Currency)
	require.WithinDuration(t, actualAccount.CreatedAt, getAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	actualAccount := createRandomAccount(t) 

	arg := UpdateAccountParams{
		ID:      actualAccount.ID,
		Balance: util.RandomInt(1, 1000),
	}

	getAccount, err := testQueries.GetAccount(context.Background(), actualAccount.ID)

	require.NoError(t, err)

	err = testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)

	updatedAccount, err := testQueries.GetAccount(context.Background(), actualAccount.ID)
	require.NoError(t, err)

	require.NotEmpty(t, getAccount)
	require.Equal(t, actualAccount.ID, getAccount.ID)
	require.Equal(t, actualAccount.Owner, getAccount.Owner)
	require.NotEqual(t, actualAccount.Balance, updatedAccount.Balance)
	require.Equal(t, actualAccount.Currency, getAccount.Currency)
	require.WithinDuration(t, actualAccount.CreatedAt, getAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	actualAccount := createRandomAccount(t) 

	err := testQueries.DeleteAccount(context.Background(), actualAccount.ID)
	require.NoError(t, err)

	deletedAccount, err := testQueries.GetAccount(context.Background(), actualAccount.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedAccount)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.NotEmpty(t, account.ID)
		require.NotEmpty(t, account.Balance)
		require.NotEmpty(t, account.Owner)
	}
}
