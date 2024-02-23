package entities

import "github.com/uptrace/bun"

type Transacao struct {
	bun.BaseModel `bun:"table:transacoes,alias:t"`
	IDCliente     string `json:"id_cliente"`
	Valor         int64  `json:"valor"`
	Tipo          string `json:"tipo"`
	Descricao     string `json:"descricao"`
	RealizadaEm   string `json:"realizada_em"`
}
