package routines

import (
	"context"
	"fmt"
	"main/entities"
	"main/pkg/database"
	"main/pkg/memory"
)

func ClientesMemoryMapping() {

	consumerName := "ClientesMemoryMapping"

	var clientes []entities.Cliente
	ctx := context.Background()

	cache := memory.GetCacheInstance()

	db := database.GetDB()
	err := db.NewSelect().Model(&clientes).OrderExpr("id_cliente ASC").Limit(10).Scan(ctx)

	if err != nil {
		fmt.Printf("[%s] Erro ao recuperar os clientes do database principal %v:\n", consumerName, err)
		return
	}

	// Criando um cache em memória dos valores que não mudam
	// Será utilizado para verificar o limite e verificar se o cliente existe
	for _, u := range clientes {
		cache.Set("cliente:"+u.ID, u.ID)
		cache.Set("limite:"+u.ID, u.Limite)
	}

	fmt.Printf("[%s] Mapping de memória finalizado:\n", consumerName)
}
