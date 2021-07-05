package schemas

type User struct {
	UUID       string  `json:"uuid" db:"uuid"`
	Name       string  `json:"name" db:"name"`
	Document   string  `json:"document" db:"document"`
	Email      string  `json:"email" db:"email"`
	Password   string  `json:"password" db:"password"`
	Wallet     float64 `json:"wallet" db:"wallet"`
	CommonUser bool    `json:"common_user" db:"common_user"`
}

type Shopkeeper struct {
	UUID     string  `json:"uuid" db:"uuid"`
	Name     string  `json:"name" db:"name"`
	CNPJ     string  `json:"cnpj" db:"cnpj"`
	Email    string  `json:"email" db:"email"`
	Password string  `json:"password" db:"password"`
	Wallet   float64 `json:"wallet" db:"wallet"`
}

type Transfer struct {
	UUID   string  `json:"uuid" database:"uuid"`
	Value  float64 `json:"value" database:"value"`
	Payer  string  `json:"payer" database:"payer"`
	Payee  string  `json:"payee" database:"payee"`
	Status string  `json:"status" database:"status"`
}
