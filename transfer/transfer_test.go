package transfer

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hemillainy/backend/config"
	error2 "github.com/hemillainy/backend/error"
	repo "github.com/hemillainy/backend/repository"
	"github.com/hemillainy/backend/repository/database"
	"github.com/hemillainy/backend/schemas"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func TestMakeTransfer(t *testing.T) {
	r, err := database.Open(context.Background(), &config.Config{})
	if err != nil {
		t.Errorf("expected no errors but was %v", err)
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

	_ = insertDefault(r, user, shopkeeper)
	//defer f()

	type args struct {
		w        *httptest.ResponseRecorder
		r        *http.Request
	}

	payload := schemas.Transfer{
		Value: 100,
		Payer: user.UUID,
		Payee: shopkeeper.UUID,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	testCases := []struct {
		description string
		args        args
		wantErr     bool
		expected	error
	}{
		{
			description: "user send money to other user with success",
			args: args{
				w:        httptest.NewRecorder(),
				r:        httptest.NewRequest(http.MethodGet, "/transaction", bytes.NewReader(b)),
			},
			wantErr: false,
			expected: nil,
		},
		{
			description: "user send money with error",
			args: args{
				w:        httptest.NewRecorder(),
				r:        httptest.NewRequest(http.MethodGet, "/transaction", bytes.NewReader(b)),
			},
			wantErr: true,
			expected: error2.ErrInvalidBalance,
		},
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			TransactionHandler(c.args.w, c.args.r)
		})
		resp := c.args.w.Result()
		b, err := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf(string(b))
			return
		}

		if err != nil {
			t.Error(err)
			return
		}
		resp.Body.Close()

		tCreated := schemas.Transfer{}
		err = json.Unmarshal(b, &tCreated)
		if !c.wantErr && err != nil {
			t.Errorf("no error expected but has: %v", err)
			return
		}

		if c.wantErr && err != c.expected {
			t.Errorf("expected: %v\n has: %v\n", c.expected, err)
			return
		}
	}
}

func TestTransactionHandlerErr(t *testing.T) {
	//testCases :=
}

func Test_notify(t *testing.T) {
	 _ = notify(schemas.Transfer{
		UUID:   uuid.NewString(),
		Value:  100,
		Payer:  uuid.NewString(),
		Payee:  uuid.NewString(),
		Status: "done",
	})

}