package services

import (
	"context"
	"fmt"
	"main/entities"
	"main/pkg/database"
)

func Extrato(id_cliente string) ([]entities.Transacao, error) {
	functionName := "Extrato"

	var transacoes []entities.Transacao

	ctx := context.Background()
	db := database.GetDB()

	err := db.NewSelect().Model(&transacoes).Where("id_cliente = ?", id_cliente).Limit(10).Order("realizada_em DESC").Scan(ctx)
	if err != nil {
		fmt.Printf("[%s] Erro ao retornar o extrato do cliente %v:\n", functionName, err)
		return transacoes, err
	}

	return transacoes, err

}
