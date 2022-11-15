package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/WrungCodes/simplebank/util"
	"github.com/stretchr/testify/require"
)

func CreateEntry(t *testing.T, account Account) Entry {

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	CreateEntry(t, CreateAccount(t))
}

func TestGetEntry(t *testing.T) {
	entry := CreateEntry(t, CreateAccount(t))

	retrieved, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, retrieved)

	require.Equal(t, entry.AccountID, retrieved.AccountID)
	require.Equal(t, entry.Amount, retrieved.Amount)

	require.Equal(t, entry.ID, retrieved.ID)
	require.Equal(t, entry.CreatedAt, retrieved.CreatedAt)
}

func TestDeleteEntry(t *testing.T) {
	entry := CreateEntry(t, CreateAccount(t))

	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	deleted, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.Empty(t, deleted)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateEntry(t, CreateAccount(t))
	}

	arg := ListAllEntriesParams{
		Limit:  5,
		Offset: 2,
	}

	entries, err := testQueries.ListAllEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestListAccountEntry(t *testing.T) {
	account := CreateAccount(t)

	for i := 0; i < 10; i++ {
		CreateEntry(t, account)
	}

	arg := ListAccountEntriesParams{
		Limit:     5,
		Offset:    2,
		AccountID: account.ID,
	}

	entries, err := testQueries.ListAccountEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
	}
}
