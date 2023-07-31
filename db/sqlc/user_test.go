package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomUsername(),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestGetUserByUsername(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)

}

func TestUpdateUserOnlyEmail(t *testing.T) {
	user := createRandomUser(t)

	newEmail := util.RandomEmail()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.Equal(t, user.ID, updatedUser.ID)
	require.Equal(t, user.Username, updatedUser.Username)
	require.NotEqual(t, user.Email, updatedUser.Email)
	require.Equal(t, user.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	user := createRandomUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username:       user.Username,
		HashedPassword: newHashedPassword,
	})

	require.NoError(t, err)
	require.Equal(t, user.ID, updatedUser.ID)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, user.Email, updatedUser.Email)
	require.NotEqual(t, user.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)

}

func TestUpdateUserAllFields(t *testing.T) {
	user := createRandomUser(t)

	newEmail := util.RandomEmail()

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username:       user.Username,
		HashedPassword: newHashedPassword,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)

	require.Equal(t, user.ID, updatedUser.ID)
	require.Equal(t, user.Username, updatedUser.Username)

	require.NotEqual(t, user.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)

	require.NotEqual(t, user.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)

}
