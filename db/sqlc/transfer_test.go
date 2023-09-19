package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"tech-school/util"
)

func createRandomTransfer(t *testing.T) Transfer {
	ctx := context.Background()

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	var fromAccountID sql.NullInt64
	fromAccountID.Int64 = account1.ID
	fromAccountID.Valid = true

	var toAccountID sql.NullInt64
	toAccountID.Int64 = account2.ID
	toAccountID.Valid = true

	arg := CreateTransferParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	require.Equal(t, arg.Amount, transfer.Amount)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID.Int64)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, account1.ID, transfer.FromAccountID.Int64)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	ctx := context.Background()

	transfer1 := createRandomTransfer(t)

	transfer2, err := testQueries.GetTransfer(ctx, transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.NotZero(t, transfer2.ID)
	require.NotZero(t, transfer2.CreatedAt)

	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
}

func TestUpdateTransfer(t *testing.T) {
	ctx := context.Background()

	transfer1 := createRandomTransfer(t)

	arg := UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: util.RandomMoney(),
	}

	transfer2, err := testQueries.UpdateTransfer(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, arg.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	ctx := context.Background()

	transfer1 := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(ctx, transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(ctx, transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}

func TestListTransfers(t *testing.T) {
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(ctx, arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, trn := range transfers {
		require.NotEmpty(t, trn)
	}

}
