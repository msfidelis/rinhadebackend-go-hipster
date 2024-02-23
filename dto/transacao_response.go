package dto

type TransacaoResponse struct {
	Limite int64 `json:"limite"`
	Saldo  int64 `json:"saldo"`
}
