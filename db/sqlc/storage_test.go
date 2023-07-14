package db

import (
	"context"
	"testing"
	"time"

	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/stretchr/testify/require"
)

func TestRegisterTx(t *testing.T) {
	store := NewStorage(testDB)

	hashedPassword, err := util.HashPassword("secret")
	require.NoError(t, err)

	// run n concurrent user registrations to ensure transaction works well
	n := 5

	errs := make(chan error)
	results := make(chan RegisterTxResult)

	for i := 0; i < 5; i++ {
		go func() {
			result, err := store.RegisterTx(context.Background(), CreateUserParams{
				Username:       "Alex",
				HashedPassword: hashedPassword,
				Email:          util.RandomEmail(),
			})

			errs <- err
			results <- result
		}()
	}
	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check createdUser. because of unique constraint on emails, need to use
		//random emails and so cannot check them.
		user := result.User
		require.NotEmpty(t, user)

		require.NotZero(t, user.ID)
		require.Equal(t, "Alex", user.Username)
		require.Equal(t, hashedPassword, user.HashedPassword)
		require.NotEmpty(t, user.Email)

		require.WithinDuration(t, user.CreatedAt, time.Now(), time.Second)

		_, err = store.GetUser(context.Background(), user.ID)
		require.NoError(t, err)
	}
}
