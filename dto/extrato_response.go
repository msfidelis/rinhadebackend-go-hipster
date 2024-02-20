package dto

import "main/entities"

type ExtratoResponse struct {
	Saldo             Saldo                `json:"saldo"`
	UltimasTransacoes []entities.Transacao `json:"ultimas_transacoes"`
}

type Saldo struct {
	Total       float64 `json:"total"`
	DataExtrato string  `json:"data_extrato"`
	Limite      float64 `json:"limite"`
}
