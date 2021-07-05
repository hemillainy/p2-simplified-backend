package transfer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	error2 "github.com/hemillainy/backend/error"
	"github.com/hemillainy/backend/rabbit"
	repo "github.com/hemillainy/backend/repository"
	database "github.com/hemillainy/backend/repository/database"
	"github.com/hemillainy/backend/schemas"
	"io/ioutil"
	"net/http"
)

type responseMessage struct {
	Message string `json:"message"`
}

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	t := schemas.Transfer{}
	err := json.NewDecoder(r.Body).Decode(&t)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	type errorMsg struct {
		error	string	`json:"error"`
	}
	repository, err := database.Open(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMsg{error: err.Error()})
		return
	}

	valid, err := validTransaction(context.Background(), t.Payer, t.Value, repository)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMsg{error: err.Error()})
		return
	}

	msg, err := authorizerRequest()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMsg{error: err.Error()})
		return
	}

	if msg != "Autorizado" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMsg{error: errors.New("transação não autorizada").Error()})
		return
	}

	tCreated, err := repository.CreateTransfer(context.Background(), t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMsg{error: err.Error()})
		return
	}

	repository.UpdateUserWallet(context.Background(), t.Value, t.Payer, true)
	repository.UpdateUserWallet(context.Background(), t.Value, t.Payee, false)

	notify(t)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tCreated)
}

func validTransaction(ctx context.Context, payerUUID string, transactionValue float64, r repo.IRepository) (bool, error) {
	payer, err := r.SelectUser(ctx, payerUUID)
	if err != nil {
		return false, err
	}
	if payer.Wallet < transactionValue {
		return false, error2.ErrInvalidTransaction
	}
	return true, nil
}

func authorizerRequest() (string, error) {
	resp, err := http.Get("https://run.mocky.io/v3/8fafdd68-a090-496f-8c9a-3442cf30dae6")
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var message responseMessage
	err = json.Unmarshal(body, &message)
	if err != nil {
		return "", err
	}
	return message.Message, nil
}

func notify(t schemas.Transfer) error {
	resp, err := http.Get("http://o4d9z.mocklab.io/notify")
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var message responseMessage
	err = json.Unmarshal(body, &message)
	if err != nil {
		return err
	}

	b := []byte(fmt.Sprintf("%v", t))
	m := rabbit.Message{
		Queue:       "notification",
		Body:        &b,
		ContentType: "application/json",
	}

	fmt.Printf("msg > %v\n", message.Message)
	if message.Message != "Sucess" {
		rabbit.Publish(context.Background(), m)
	}
	return nil
}