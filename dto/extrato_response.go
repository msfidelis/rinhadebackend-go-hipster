package dto

import "main/entities"

type ExtratoResponse struct {
	Saldo             Saldo                `json:"saldo"`
	UltimasTransacoes []entities.Transacao `json:"ultimas_transacoes"`
}

type Saldo struct {
	Total       int64  `json:"total"`
	DataExtrato string `json:"data_extrato"`
	Limite      int64  `json:"limite"`
}
