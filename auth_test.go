package ledger

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockUserGetter struct {
	mock.Mock
}

func (m *mockUserGetter) GetUser(ctx context.Context, uuid string) (*User, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(*User), args.Error(1)
}

func Test_AuthenticateUser(t *testing.T) {
	expectedUser := &User{
		UserUUID:     "test",
		Login:        "test",
		PasswordHash: "passwordHash",
		PasswordSalt: "salt",
	}

	ctx := context.Background()
	repo := &mockUserGetter{}
	repo.On("GetUser", ctx, "test").Return(expectedUser, nil)
	hashFunc := func(password string, salt string, pepper string) string {
		assert.Equal(t, "password", password)
		assert.Equal(t, "salt", salt)
		assert.Equal(t, "pepper", pepper)
		return "passwordHash"
	}

	user, err := AuthenticateUser(ctx, repo, hashFunc, "test", "password")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	repo.AssertExpectations(t)
}
