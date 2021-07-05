package repository

import (
	"context"
	"github.com/hemillainy/backend/schemas"
)

type IRepository interface {
	CreateTransfer(ctx context.Context, transfer schemas.Transfer) (schemas.Transfer, error)
	SelectUser(ctx context.Context, payerUUID string) (u schemas.User, err error)
	CreateUser(ctx context.Context, u schemas.User) (err error)
	DeleteUser(ctx context.Context, UUID string) (err error)
	UpdateUserWallet(ctx context.Context, value float64, UUID string, commonUser bool) (schemas.User, error)
	Close() (err error)
}
