package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/haikeloksama/gobank/util"
	"github.com/stretchr/testify/require"
)

func createTestTransfer(t *testing.T) (from Account, to Account, tf Transfer) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	args := CreateTransferParams {
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: util.RandomInt(1, 1000),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.Amount)
	require.NotZero(t, transfer.CreatedAt)

	return account1, account2, transfer
}

func TestCreateTransfer(t *testing.T) {
	createTestTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	_, _, actualTransfer := createTestTransfer(t)
	transfer, err := testQueries.GetTransfer(context.Background(), actualTransfer.ID)
	
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, actualTransfer.ID, transfer.ID)
	require.Equal(t, actualTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, actualTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, actualTransfer.Amount, transfer.Amount)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i <= 10 ;i++ {
		createTestTransfer(t)
	} 

	args := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.NotEmpty(t, transfer.ID)
		require.NotEmpty(t, transfer.FromAccountID)
		require.NotEmpty(t, transfer.ToAccountID)
		require.NotEmpty(t, transfer.Amount)
	}
}

func TestUpdateTransfer(t *testing.T) {
	_, _, actualTransfer := createTestTransfer(t)
	
	args := UpdateTransferParams {
		ID: actualTransfer.ID,
		Amount: util.RandomInt(1, 1000),
	}

	err := testQueries.UpdateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotZero(t,  actualTransfer.ID)
	require.NotZero(t,  actualTransfer.Amount)
	transfer, err := testQueries.GetTransfer(context.Background(), actualTransfer.ID)

	require.NoError(t, err)
	require.NotEqual(t, actualTransfer.Amount, transfer.Amount)
}

func TestDeleteTransfer(t *testing.T) {
	_, _, actualTransfer := createTestTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), actualTransfer.ID)
	require.NoError(t, err)

	transfer, err := testQueries.GetTransfer(context.Background(), actualTransfer.ID)	

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer)
}