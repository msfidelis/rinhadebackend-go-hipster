package services

import (
	"context"
	"log"
	"main/entities"
	"main/pkg/database"
	"time"
)

func Extrato(id_cliente string) ([]entities.Transacao, error) {
	var transacoes []entities.Transacao

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := database.GetDB()
	err := db.NewSelect().
		Model(&transacoes).
		Where("id_cliente = ?", id_cliente).
		Limit(10).
		Order("realizada_em DESC").
		Scan(ctx)
	if err != nil {
		log.Printf("[%s] Erro ao retornor o extrato do cliente %s: %v", "Extrato", id_cliente, err)
		return nil, err
	}

	return transacoes, nil

}
