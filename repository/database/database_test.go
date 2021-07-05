package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/hemillainy/backend/config"
	repo "github.com/hemillainy/backend/repository"
	"github.com/hemillainy/backend/schemas"
	"testing"
)

func insertDefault(r repo.IRepository, u schemas.User, s schemas.Shopkeeper) func() {
	r.CreateUser(context.Background(), u)
	r.CreateShopkeeper(context.Background(), s)
	return func() {
		r.DeleteUser(context.Background(), u.UUID)
		r.DeleteShopkeeper(context.Background(), s.UUID)
	}
}

func TestCreateTransfer(t *testing.T) {
	repo, err := Open(context.Background(), &config.Config{})
	if err != nil {
		t.Errorf("no error expected but has: %v", err)
		return
	}

	user := schemas.User{
		UUID:     uuid.NewString(),
		Name:     "User Test",
		CPF:      "11111111111",
		Email:    "user@backend.com",
		Password: "senha123",
		Wallet:   200,
	}
	shopkeeper := schemas.Shopkeeper{
		UUID:     uuid.NewString(),
		Name:     "Shopkeeper",
		CNPJ:     "11111111111111",
		Email:    "shopkeeper@backend.com",
		Password: "senha123",
		Wallet:   0,
	}

	_ = insertDefault(repo, user, shopkeeper)
	//defer f()

	value := 100.00
	status := "done"

	_, err = repo.CreateTransfer(context.Background(), schemas.Transfer{
		Value:  value,
		Payer:  user.UUID,
		Payee:  shopkeeper.UUID,
		Status: status,
	})

	if err != nil {
		t.Errorf("error no expected but has %v", err.Error())
		return
	}
}

func TestRepository_UpdateUserWallet(t *testing.T) {

}

