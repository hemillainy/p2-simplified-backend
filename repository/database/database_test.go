package database

import (
	"context"
	"github.com/google/uuid"
	repo "github.com/hemillainy/backend/repository"
	"github.com/hemillainy/backend/schemas"
	"os"
	"testing"
)

func setEnvs() {
	os.Setenv("BACKEND_DB_HOST", "localhost")
	os.Setenv("BACKEND_DB_SCHEME", "postgres")
	os.Setenv("BACKEND_USER", "postgres")
	os.Setenv("BACKEND_PASSWORD", "postgres")
	os.Setenv("BACKEND_DB_PORT", "5432")
	os.Setenv("BACKEND_DB_SCHEME", "postgres")
	os.Setenv("BACKEND_DB_NAME", "backend")
}

func insertDefault(r repo.IRepository, user schemas.User, shopkeeper schemas.User) func() {
	r.CreateUser(context.Background(), user)
	r.CreateUser(context.Background(), shopkeeper)
	return func() {
		r.DeleteUser(context.Background(), user.UUID)
		r.DeleteUser(context.Background(), shopkeeper.UUID)
	}
}

func TestCreateTransfer(t *testing.T) {
	setEnvs()
	repo, err := Open(context.Background())
	if err != nil {
		t.Errorf("no error expected but has: %v", err)
		return
	}

	commonUser := schemas.User{
		UUID:       uuid.NewString(),
		Name:       "User Test",
		Document:   "11111111111",
		Email:      "user@backend.com",
		Password:   "senha123",
		Wallet:     200,
		CommonUser: true,
	}
	shopkeeper := schemas.User{
		UUID:       uuid.NewString(),
		Name:       "Shopkeeper Test",
		Document:   "11111111111111",
		Email:      "shopkeeper@backend.com",
		Password:   "senha123",
		Wallet:     0,
		CommonUser: false,
	}

	f := insertDefault(repo, commonUser, shopkeeper)
	defer f()

	value := 100.00
	status := "done"

	transferCreated, err := repo.CreateTransfer(context.Background(), schemas.Transfer{
		Value:  value,
		Payer:  commonUser.UUID,
		Payee:  shopkeeper.UUID,
		Status: status,
	})

	if err != nil {
		t.Errorf("error no expected but has %v", err.Error())
		return
	}

	if value != transferCreated.Value {
		t.Errorf("expected %v but has %v", value, transferCreated.Value)
		return
	}
}

func TestRepository_UpdateUserWallet(t *testing.T) {
	setEnvs()
	repo, err := Open(context.Background())
	if err != nil {
		t.Errorf("no error expected but has: %v", err)
		return
	}

	defer repo.Close()

	commonUser := schemas.User{
		UUID:       uuid.NewString(),
		Name:       "User Test",
		Document:   "11111111111",
		Email:      "user@backend.com",
		Password:   "senha123",
		Wallet:     200,
		CommonUser: true,
	}
	shopkeeper := schemas.User{
		UUID:       uuid.NewString(),
		Name:       "Shopkeeper Test",
		Document:   "11111111111111",
		Email:      "shopkeeper@backend.com",
		Password:   "senha123",
		Wallet:     0,
		CommonUser: false,
	}

	f := insertDefault(repo, commonUser, shopkeeper)
	defer f()

	testCases := []struct {
		description string
		value       float64
		expected    float64
		commonUser  bool
		user        schemas.User
	}{
		{
			description: "decreases user wallet",
			value:       100,
			expected:    100,
			commonUser:  true,
			user:        commonUser,
		},
		{
			description: "increase the user's wallet",
			value:       100,
			expected:    100,
			commonUser:  false,
			user:        shopkeeper,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			u, err := repo.UpdateUserWallet(context.Background(), test.value, test.user.UUID, test.commonUser)
			if err != nil {
				t.Errorf("no error expected but has: %v\n", err)
				return
			}

			if test.expected != u.Wallet {
				t.Errorf("expected %v but has %v", test.expected, u.Wallet)
				return
			}
		})
	}
}
