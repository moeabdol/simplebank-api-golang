package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/moeabdol/simplebank-api-golang/utils"
	"github.com/stretchr/testify/require"
)

func createTestAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)
	if err != nil {
		t.Error(err)
	}

	return account
}

func deleteTestAccount(t *testing.T, id int64) {
	err := testStore.DeleteAccount(context.Background(), id)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.Owner, arg.Owner)
	require.Equal(t, account.Balance, arg.Balance)
	require.Equal(t, account.Currency, arg.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.UpdatedAt)

	deleteTestAccount(t, account.ID)
}

func TestGetAccount(t *testing.T) {
	account1 := createTestAccount(t)
	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.WithinDuration(t, account1.UpdatedAt, account2.UpdatedAt, time.Second)

	deleteTestAccount(t, account1.ID)
	deleteTestAccount(t, account2.ID)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createTestAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}

	account2, err := testStore.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.NotEqual(t, account1.Balance, account2.Balance)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
	require.NotEqual(t, account1.UpdatedAt, account2.UpdatedAt)

	deleteTestAccount(t, account1.ID)
	deleteTestAccount(t, account2.ID)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createTestAccount(t)

	err := testStore.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var randomAccounts []Account
	for i := 0; i < 10; i++ {
		randomAccounts = append(randomAccounts, createTestAccount(t))
	}

	arg := ListAccountsParams{
		Limit:  100,
		Offset: 0,
	}

	accounts, err := testStore.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 10)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

	for i := 0; i < 10; i++ {
		deleteTestAccount(t, randomAccounts[i].ID)
	}
}
