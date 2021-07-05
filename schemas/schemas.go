package schemas

type User struct {
	UUID     string  `json:"uuid" database:"uuid"`
	Name     string  `json:"name" database:"name"`
	CPF      string  `json:"cpf" database:"cpf"`
	Email    string  `json:"email" database:"email"`
	Password string  `json:"password" database:"password"`
	Wallet   float64 `json:"wallet" database:"wallet"`
}

type Shopkeeper struct {
	UUID     string `json:"uuid" database:"uuid"`
	Name     string `json:"name" database:"name"`
	CNPJ     string `json:"cnpj" database:"cnpj"`
	Email    string `json:"email" database:"email"`
	Password string `json:"password" database:"password"`
	Wallet   float64 `json:"wallet" database:"wallet"`
}

type Transfer struct {
	UUID   string  `json:"uuid" database:"uuid"`
	Value  float64 `json:"value" database:"value"`
	Payer  string  `json:"payer" database:"payer"`
	Payee  string  `json:"payee" database:"payee"`
	Status string  `json:"status" database:"status"`
}
