package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"tech-school/util"
)

func createRandomEntry(t *testing.T) Entry {
	ctx := context.Background()

	account := createRandomAccount(t)

	var accID sql.NullInt64
	accID.Int64 = account.ID
	accID.Valid = true

	arg := CreateEntryParams{
		AccountID: accID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, account.ID, entry.AccountID.Int64)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	ctx := context.Background()

	entry1 := createRandomEntry(t)

	entry2, err := testQueries.GetEntry(ctx, entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	ctx := context.Background()

	entry1 := createRandomEntry(t)

	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandomMoney(),
	}

	entry2, err := testQueries.UpdateEntry(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, arg.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	ctx := context.Background()

	entry1 := createRandomEntry(t)

	err := testQueries.DeleteEntry(ctx, entry1.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(ctx, entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestListEntries(t *testing.T) {
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(ctx, arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, e := range entries {
		require.NotEmpty(t, e)
	}
}
