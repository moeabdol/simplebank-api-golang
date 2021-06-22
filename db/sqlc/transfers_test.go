package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/moeabdol/simplebank-api-golang/utils"
	"github.com/stretchr/testify/require"
)

func createTestTransfer(t *testing.T, fromAccountID, toAccountID int64) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testStore.CreateTransfer(context.Background(), arg)
	if err != nil {
		t.Error(err)
	}

	return transfer
}

func deleteTestTransfer(t *testing.T, id int64) {
	err := testStore.DeleteTransfer(context.Background(), id)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createTestAccount(t)
	toAccount := createTestAccount(t)

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testStore.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)

	require.NotEmpty(t, transfer.ID)
	require.NotEmpty(t, transfer.FromAccountID)
	require.NotEmpty(t, transfer.ToAccountID)
	require.NotEmpty(t, transfer.Amount)
	require.NotEmpty(t, transfer.CreatedAt)
	require.NotEmpty(t, transfer.UpdatedAt)

	deleteTestTransfer(t, transfer.ID)
	deleteTestAccount(t, fromAccount.ID)
	deleteTestAccount(t, toAccount.ID)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createTestAccount(t)
	toAccount := createTestAccount(t)
	transfer1 := createTestTransfer(t, fromAccount.ID, toAccount.ID)

	transfer2, err := testStore.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
	require.Equal(t, transfer1.UpdatedAt, transfer2.UpdatedAt)

	deleteTestTransfer(t, transfer2.ID)
	deleteTestAccount(t, fromAccount.ID)
	deleteTestAccount(t, toAccount.ID)
}

func TestListTransfers(t *testing.T) {
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)

	var testTransfers []Transfer
	for i := 0; i < 10; i++ {
		testTransfers = append(testTransfers, createTestTransfer(t, account1.ID, account2.ID))
		testTransfers = append(testTransfers, createTestTransfer(t, account2.ID, account1.ID))
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         100,
		Offset:        0,
	}

	transfers, err := testStore.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 10)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account2.ID)
	}

	for _, transfer := range testTransfers {
		deleteTestTransfer(t, transfer.ID)
	}
	deleteTestAccount(t, account1.ID)
	deleteTestAccount(t, account2.ID)
}

func TestDeleteTransfer(t *testing.T) {
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)
	transfer1 := createTestTransfer(t, account1.ID, account2.ID)

	err := testStore.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testStore.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)

	deleteTestAccount(t, account1.ID)
	deleteTestAccount(t, account2.ID)
}
