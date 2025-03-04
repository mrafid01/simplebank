package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5              // number of transactions that occur
	amount := int64(10) // money to transfer

	errChnl := make(chan error)
	resultChnl := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})
			errChnl <- err
			resultChnl <- result
		}()
	}

	// Checks results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errChnl
		require.NoError(t, err)

		result := <-resultChnl
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check accounts balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balances
	updatesAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatesAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatesAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatesAccount2.Balance)
}

// 2 accounts transfer each other
func TestTransferDeadlock(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10            // number of transactions that occur
	amount := int64(5) // money to transfer

	errChnl := make(chan error)
	resultChnl := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		accountID1 := account1.ID
		accountID2 := account2.ID
		if i%2 == 1 {
			accountID1, accountID2 = accountID2, accountID1
		}
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: accountID1,
				ToAccountId:   accountID2,
				Amount:        amount,
			})
			errChnl <- err
			resultChnl <- result
		}()
	}

	// To execute to DB
	// without this DB won't update the transfer data
	for i := 0; i < n; i++ {
		err := <-errChnl
		require.NoError(t, err)
		result := <-resultChnl
		require.NotEmpty(t, result)
	}

	// check the final updated balances
	updatesAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatesAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatesAccount1.Balance)
	require.Equal(t, account2.Balance, updatesAccount2.Balance)
}
