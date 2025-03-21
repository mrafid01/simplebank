package db

import (
	"context"
	"testing"
	"time"

	"github.com/mrafid01/simplebank/util"
	"github.com/stretchr/testify/assert"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), args)
	assert.Nil(t, err)
	assert.NotEmpty(t, account)

	assert.Equal(t, args.Owner, account.Owner)
	assert.Equal(t, args.Balance, account.Balance)
	assert.Equal(t, args.Currency, account.Currency)

	assert.NotEmpty(t, account.ID)
	assert.NotEmpty(t, account.CreatedAt)

	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, account2)

	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	assert.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	args := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testStore.UpdateAccount(context.Background(), args)
	assert.Nil(t, err)
	assert.NotEmpty(t, account2)

	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, args.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	assert.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testStore.DeleteAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrRecordNotFound.Error())
	assert.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i <= 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	args := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testStore.ListAccounts(context.Background(), args)
	assert.NoError(t, err)
	assert.NotEmpty(t, accounts)

	for _, account := range accounts {
		assert.NotEmpty(t, account)
		assert.Equal(t, account.Owner, lastAccount.Owner)
	}
}
