package database

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	repo "github.com/hemillainy/backend/repository"
	"github.com/hemillainy/backend/schemas"
	"github.com/jmoiron/sqlx"
	"gocloud.dev/postgres"
	"os"
	"strconv"
)

type Repository struct {
	db *sqlx.DB
}

// Open database and return DB instance
func Open(ctx context.Context) (r repo.IRepository, err error) {
	scheme := os.Getenv("BACKEND_DB_SCHEME")
	user := os.Getenv("BACKEND_USER")
	password := os.Getenv("BACKEND_PASSWORD")
	host := os.Getenv("BACKEND_DB_HOST")
	port := os.Getenv("BACKEND_DB_PORT")
	portInt, err := strconv.Atoi(port)
	DBName := os.Getenv("BACKEND_DB_NAME")

	psqlInfo := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", scheme, user, password, host, portInt, DBName)
	d, err := postgres.Open(ctx, psqlInfo)
	if err != nil {
		err = fmt.Errorf("error open database: %v", err)
		return &Repository{}, err
	}
	db := sqlx.NewDb(d, "postgresql")
	err = db.Ping()
	if err != nil {
		err = fmt.Errorf("error ping database: %v", err)
		return &Repository{}, err
	}
	r = &Repository{db}
	return
}

func (r *Repository) Close() (err error) {
	err = r.db.Close()
	return
}

func (r *Repository) CreateTransfer(ctx context.Context, transfer schemas.Transfer) (t schemas.Transfer, err error) {
	UUID := uuid.NewString()
	SQL := `INSERT INTO transfers(uuid, 
				value, 
				payer, 
				payee, 
				status) 
			VALUES($1, $2, $3, $4, $5);`

	_, err = r.db.ExecContext(ctx, SQL, UUID, transfer.Value, transfer.Payer, transfer.Payee, "done")
	SQL = `SELECT * FROM transfers WHERE uuid = $1`
	err = r.db.GetContext(ctx, &t, SQL, UUID)
	return
}

func (r *Repository) UpdateUserWallet(ctx context.Context, value float64, UUID string, commonUser bool) (u schemas.User, err error) {
	var SQL string
	if commonUser {
		SQL = `UPDATE users 
			SET wallet = wallet - $1
			WHERE uuid = $2`
	} else {
		SQL = `UPDATE users 
			SET wallet = wallet + $1
			WHERE uuid = $2`
	}
	_, err = r.db.ExecContext(ctx, SQL, value, UUID)
	if err != nil {
		return
	}

	SQL = `SELECT * FROM users WHERE uuid = $1`
	err = r.db.GetContext(ctx, &u, SQL, UUID)
	return
}

func (r *Repository) CreateUser(ctx context.Context, u schemas.User) (err error) {
	SQL := `INSERT INTO users(uuid, 
				name,
				document,
				email,
				password,
				wallet,
                common_user) 
		    VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = r.db.ExecContext(ctx, SQL, u.UUID, u.Name, u.Document, u.Email, u.Password, u.Wallet, u.CommonUser)
	return
}

func (r *Repository) SelectUser(ctx context.Context, payerUUID string) (u schemas.User, err error) {
	SQL := `SELECT * FROM users WHERE uuid = $1`
	err = r.db.GetContext(ctx, &u, SQL, payerUUID)
	return
}

func (r *Repository) DeleteUser(ctx context.Context, UUID string) (err error) {
	SQL := `DELETE FROM "users" WHERE "uuid"=$1`
	_, err = r.db.ExecContext(ctx, SQL, UUID)
	return
}
