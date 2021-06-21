package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/moeabdol/simplebank-api-golang/utils"
	"github.com/stretchr/testify/require"
)

func createTestEntry(t *testing.T, id int64) Entry {
	arg := CreateEntryParams{
		AccountID: id,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	if err != nil {
		t.Error(err)
	}

	return entry
}

func deleteTestEntry(t *testing.T, id int64) {
	err := testQueries.DeleteEntry(context.Background(), id)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateEntry(t *testing.T) {
	account := createTestAccount(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.UpdatedAt)

	deleteTestEntry(t, entry.ID)
	deleteTestAccount(t, account.ID)
}

func TestGetEntry(t *testing.T) {
	account := createTestAccount(t)
	entry1 := createTestEntry(t, account.ID)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	require.Equal(t, entry1.UpdatedAt, entry2.UpdatedAt)

	deleteTestEntry(t, entry2.ID)
	deleteTestAccount(t, account.ID)
}

func TestDeleteEntry(t *testing.T) {
	account := createTestAccount(t)
	entry1 := createTestEntry(t, account.ID)

	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)

	deleteTestAccount(t, account.ID)
}

func TestListEntries(t *testing.T) {
	account := createTestAccount(t)
	var randomEntries []Entry
	for i := 0; i < 10; i++ {
		randomEntries = append(randomEntries, createTestEntry(t, account.ID))
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     100,
		Offset:    0,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, 10)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, account.ID)
	}

	for i := 0; i < 10; i++ {
		deleteTestEntry(t, entries[i].ID)
	}

	deleteTestAccount(t, account.ID)
}
