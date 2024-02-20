package entities

import "github.com/uptrace/bun"

type Cliente struct {
	bun.BaseModel `bun:"table:clientes,alias:u"`
	ID            string  `json:"id" bun:"id_cliente,pk,autoincrement"`
	Saldo         float64 `json:"saldo"`
	Limite        float64 `json:"limite"`
}
