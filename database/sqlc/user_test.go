package database

import (
	"context"
	"testing"

	"github.com/akshay237/backend-with-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	// 1. create the args passed to the query
	args := CreateUSerParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	// 2. calls the create Account function generated through sqlc
	user, err := testQueries.CreateUSer(context.Background(), args)

	// 3. Check for no error and other properties using testify library
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {

	// 1. create a random user
	user1 := createRandomUser(t)

	// 2. fetch the user details based on the  username
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
}
