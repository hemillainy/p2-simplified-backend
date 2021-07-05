package database

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hemillainy/backend/config"
	repo "github.com/hemillainy/backend/repository"
	"github.com/hemillainy/backend/schemas"
	"github.com/jmoiron/sqlx"
	"gocloud.dev/postgres"
)

type Repository struct {
	db *sqlx.DB
}

//cfg.DB.Scheme,
//cfg.DB.User,
//cfg.DB.Password,
//cfg.DB.Host,
//cfg.DB.Port,
//cfg.DB.Name,

// Open database and return DB instance
func Open(ctx context.Context, cfg *config.Config) (r repo.IRepository, err error) {
	//psqlInfo := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
	//	"postgresql",
	//	"postgres",
	//	"postgres",
	//	"postgres",
	//	5432,
	//	"backend",
	//)

	psqlInfo := "postgres://postgres:postgres@localhost:5432/backend?sslmode=disable"
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

	r.UpdateUserWallet(ctx, t.Value, transfer.Payer)
	r.UpdateShopkeeperWallet(ctx, t.Value, transfer.Payee)

	err = r.db.GetContext(ctx, &t, SQL, UUID)
	return
}

func (r *Repository) CreateUser(ctx context.Context, u schemas.User) (err error) {
	SQL := `INSERT INTO users(uuid, 
				name,
				cpf,
				email,
				password,
				wallet) 
		    VALUES ($1, $2, $3, $4, $5, $6)`

	fmt.Printf("INSERT\n")
	_, err = r.db.ExecContext(ctx, SQL, u.UUID, u.Name, u.CPF, u.Email, u.Password, u.Wallet)
	return
}

func (r *Repository) CreateShopkeeper(ctx context.Context, s schemas.Shopkeeper) (err error) {
	SQL := `INSERT INTO shopkeepers(uuid,
                       email,
                       password,
                       name,
                       cnpj,
                       wallet)
			   	VALUES($1, $2, $3, $4, $5, $6)`

	_, err = r.db.ExecContext(ctx, SQL, s.UUID, s.Email, s.Password, s.Name, s.CNPJ, s.Wallet)
	return
}

func (r *Repository) SelectUser(ctx context.Context, payerUUID string) (u schemas.User, err error) {
	SQL := `SELECT * FROM users WHERE uuid = $1`
	err = r.db.GetContext(ctx, &u, SQL, payerUUID)
	return
}


func (r *Repository) DeleteUser(ctx context.Context, UUID string) (err error){
	SQL := `DELETE FROM "users" WHERE "uuid"=$1`
	_, err = r.db.ExecContext(ctx, SQL, UUID)
	return
}

func (r *Repository) DeleteShopkeeper(ctx context.Context, UUID string) (err error){
	SQL := `DELETE FROM "shopkeepers" WHERE "uuid"=$1`
	_, err = r.db.ExecContext(ctx, SQL, UUID)
	return
}

func (r *Repository) UpdateUserWallet(ctx context.Context, value float64, UUID string) (err error) {
	SQL := `UPDATE users 
			SET wallet = u.wallet - $1
			FROM (
				SELECT wallet FROM users WHERE uuid = $2
			) AS u
			WHERE uuid = $2`
	_, err = r.db.ExecContext(ctx, SQL, value, UUID)
	fmt.Printf("err > %v\n", err)
	return
}

func (r *Repository) UpdateShopkeeperWallet(ctx context.Context, value float64, UUID string) (err error) {
	SQL := `UPDATE shopkeepers 
			SET wallet = s.wallet + $1
			FROM (
				SELECT wallet FROM shopkeepers WHERE uuid = $2
			) AS s
			WHERE uuid = $2`
	_, err = r.db.ExecContext(ctx, SQL, value, UUID)
	fmt.Printf("err > %v\n", err)
	return
}