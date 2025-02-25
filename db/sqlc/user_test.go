package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/mrafid01/simplebank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashedPassword(util.RandomString(6))
	require.NoError(t, err)

	args := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), args)
	assert.Nil(t, err)
	assert.NotEmpty(t, user)

	assert.Equal(t, args.Username, user.Username)
	assert.Equal(t, args.HashedPassword, user.HashedPassword)
	assert.Equal(t, args.FullName, user.FullName)
	assert.Equal(t, args.Email, user.Email)

	assert.True(t, user.PasswordChangedAt.IsZero())
	assert.NotEmpty(t, user.CreatedAt)

	return user
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	assert.Nil(t, err)
	assert.NotEmpty(t, user2)

	assert.Equal(t, user1.Username, user2.Username)
	assert.Equal(t, user1.HashedPassword, user2.HashedPassword)
	assert.Equal(t, user1.FullName, user2.FullName)
	assert.Equal(t, user1.Email, user2.Email)
	assert.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	assert.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, updatedUser.FullName, newFullName)
	require.Equal(t, updatedUser.Email, oldUser.Email)
	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, updatedUser.Email, oldUser.Email)
	require.Equal(t, updatedUser.Email, newEmail)
	require.Equal(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashedPassword(newPassword)
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.HashedPassword, newHashedPassword)
	require.Equal(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, updatedUser.Email, oldUser.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashedPassword(newPassword)
	require.NoError(t, err)
	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, updatedUser.FullName, newFullName)
	require.NotEqual(t, updatedUser.Email, oldUser.Email)
	require.Equal(t, updatedUser.Email, newEmail)
	require.NotEqual(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.HashedPassword, newHashedPassword)
}
