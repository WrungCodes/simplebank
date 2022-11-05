package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/WrungCodes/simplebank/util"
	"github.com/stretchr/testify/require"
)

func CreateAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.Balance, arg.Balance)
	require.Equal(t, account.Owner, arg.Owner)
	require.Equal(t, account.Currency, arg.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	CreateAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := CreateAccount(t)

	retrieved, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, retrieved)

	require.Equal(t, account.Balance, retrieved.Balance)
	require.Equal(t, account.Owner, retrieved.Owner)
	require.Equal(t, account.Currency, retrieved.Currency)
	require.Equal(t, account.ID, retrieved.ID)
	require.Equal(t, account.CreatedAt, retrieved.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	account := CreateAccount(t)

	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}

	updated, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updated)

	require.Equal(t, updated.Balance, arg.Balance)

	require.Equal(t, account.Owner, updated.Owner)
	require.Equal(t, account.Currency, updated.Currency)
	require.Equal(t, account.ID, updated.ID)
	require.Equal(t, account.CreatedAt, updated.CreatedAt)

	retrieved, err := testQueries.GetAccount(context.Background(), account.ID)

	require.Equal(t, updated.Balance, retrieved.Balance)
	require.Equal(t, updated.Owner, retrieved.Owner)
	require.Equal(t, updated.Currency, retrieved.Currency)
	require.Equal(t, updated.ID, retrieved.ID)
	require.Equal(t, updated.CreatedAt, retrieved.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	account := CreateAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	deleted, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Empty(t, deleted)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListAcconts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
