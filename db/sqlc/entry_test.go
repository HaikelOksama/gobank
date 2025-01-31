package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/haikeloksama/gobank/util"
	"github.com/stretchr/testify/require"
)

func createTestEntry(t *testing.T) (Account, Entry) {
	account := createRandomAccount(t)

	args := CreateEntryParams {
		AccountID: account.ID,
		Amount: util.RandomInt(1, 1000),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)

	require.Equal(t, account.ID, entry.AccountID)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return account, entry
}

func TestCreateEntry(t *testing.T) {
	createTestEntry(t)
}

func TestGetEntry(t *testing.T) {
	_, entry := createTestEntry(t)

	getEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getEntry)

	require.Equal(t, entry.ID, getEntry.ID)

}

func TestUpdateEntry(t *testing.T) {
	_, entry := createTestEntry(t)

	args := UpdateEntryParams {
		ID: entry.ID,
		Amount: util.RandomInt(1, 1000),
	}

	err := testQueries.UpdateEntry(context.Background(), args)
	require.NoError(t, err)
	newEntry,err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEqual(t, entry.Amount, newEntry.Amount)
}

func TestDeleteEntry(t *testing.T) {
	_ , entry := createTestEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)

	require.NoError(t, err)

	getEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, getEntry)

}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestEntry(t)
	}

	args := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
	
}
