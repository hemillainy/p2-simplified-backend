package transfer

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	error2 "github.com/hemillainy/backend/error"
	repo "github.com/hemillainy/backend/repository"
	"github.com/hemillainy/backend/repository/database"
	"github.com/hemillainy/backend/schemas"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func insertDefault(r repo.IRepository, commonUser schemas.User, shopkeeper schemas.User) func() {
	r.CreateUser(context.Background(), commonUser)
	r.CreateUser(context.Background(), shopkeeper)
	return func() {
		r.DeleteUser(context.Background(), commonUser.UUID)
		r.DeleteUser(context.Background(), shopkeeper.UUID)
	}
}

func TestTransactionHandler(t *testing.T) {
	setEnvs()
	r, err := database.Open(context.Background())
	if err != nil {
		t.Errorf("expected no errors but was %v", err)
		return
	}

	commonUser := schemas.User{
		UUID:       uuid.NewString(),
		Name:       "User Test",
		Document:   "11111111111",
		Email:      "user@backend.com",
		Password:   "senha123",
		Wallet:     100.00,
		CommonUser: true,
	}
	shopkeeper := schemas.User{
		UUID:       uuid.NewString(),
		Name:       "Shopkeeper",
		Document:   "11111111111111",
		Email:      "shopkeeper@backend.com",
		Password:   "senha123",
		Wallet:     0.00,
		CommonUser: false,
	}

	f := insertDefault(r, commonUser, shopkeeper)
	defer f()

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	payload := schemas.Transfer{
		Value: 200.00,
		Payer: commonUser.UUID,
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
		expected    int
	}{
		{
			description: "user send money to other user with success",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/transaction", bytes.NewReader(b)),
			},
			wantErr:  false,
			expected: http.StatusCreated,
		},
		{
			description: "user send money with error",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/transaction", bytes.NewReader(b)),
			},
			wantErr:  true,
			expected: http.StatusBadRequest,
		},
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			TransactionHandler(c.args.w, c.args.r)
		})
		resp := c.args.w.Result()
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}

		resp.Body.Close()

		if c.wantErr && resp.StatusCode != c.expected {
			t.Errorf("expected: %v\n has: %v\n", c.expected, err)
			return
		}
	}
}

func Test_validTransaction(t *testing.T) {
	setEnvs()
	r, err := database.Open(context.Background())
	if err != nil {
		t.Errorf("expected no errors but was %v", err)
		return
	}

	commonUser := schemas.User{
		UUID:       uuid.NewString(),
		Name:       "User Test",
		Document:   "11111111111",
		Email:      "user@backend.com",
		Password:   "senha123",
		Wallet:     100.00,
		CommonUser: true,
	}
	shopkeeper := schemas.User{
		UUID:       uuid.NewString(),
		Name:       "Shopkeeper",
		Document:   "11111111111111",
		Email:      "shopkeeper@backend.com",
		Password:   "senha123",
		Wallet:     0.00,
		CommonUser: false,
	}

	f := insertDefault(r, commonUser, shopkeeper)
	defer f()

	testCases := []struct {
		description string
		value       float64
		expected    bool
	}{
		{
			description: "invalid transaction",
			value:       200,
			expected:    false,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			valid, err := validTransaction(context.Background(), commonUser.UUID, 200, r)
			if err != error2.ErrInvalidTransaction {
				t.Errorf("expected %v but has %v", error2.ErrInvalidTransaction, err)
				return
			}
			if valid != test.expected {
				t.Errorf("expected %v but has %v", valid, test.expected)
				return
			}

		})
	}

}

