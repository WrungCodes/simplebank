package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/WrungCodes/simplebank/util"
	"github.com/stretchr/testify/require"
)

func CreateTransfer(t *testing.T, account1 Account, account2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.Amount, transfer.Amount)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	CreateTransfer(t, CreateAccount(t), CreateAccount(t))
}

func TestGetTransfer(t *testing.T) {
	transfer := CreateTransfer(t, CreateAccount(t), CreateAccount(t))

	retrieved, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, retrieved)

	require.Equal(t, transfer.ID, retrieved.ID)
	require.Equal(t, transfer.FromAccountID, retrieved.FromAccountID)
	require.Equal(t, transfer.ToAccountID, retrieved.ToAccountID)
	require.Equal(t, transfer.CreatedAt, retrieved.CreatedAt)
	require.Equal(t, transfer.Amount, retrieved.Amount)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := CreateTransfer(t, CreateAccount(t), CreateAccount(t))

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	retrieved, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Error(t, err)
	require.Empty(t, retrieved)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListAllTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateTransfer(t, CreateAccount(t), CreateAccount(t))
	}

	arg := ListAllTransfersParams{
		Offset: 5,
		Limit:  5,
	}

	transfers, err := testQueries.ListAllTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestListAccountTransfers(t *testing.T) {
	account := CreateAccount(t)

	for i := 0; i < 10; i++ {
		CreateTransfer(t, CreateAccount(t), account)
		CreateTransfer(t, account, CreateAccount(t))
	}

	arg := ListAccountTransfersParams{
		Offset:    5,
		Limit:     5,
		AccountID: account.ID,
	}

	transfers, err := testQueries.ListAccountTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Contains(t, []int64{transfer.FromAccountID, transfer.ToAccountID}, account.ID)
	}
}
