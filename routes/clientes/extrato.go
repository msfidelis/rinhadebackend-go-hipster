package clientes

import (
	"main/dto"
	"main/pkg/memory"
	"main/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Extrato(c *fiber.Ctx) error {

	id := c.Params("id")

	// Checa no cache em memória da aplicação se o cliente existe
	cache := memory.GetCacheInstance()
	_, found := cache.Get("cliente:" + id)
	if !found {
		return dto.FiberError(c, fiber.StatusNotFound, "cliente não encontrado")
	}

	cliente, err := services.BuscaCliente(id)
	if err != nil {
		return dto.FiberError(c, fiber.StatusInternalServerError, "Erro ao recuperar o cliente")
	}

	transacoes, err := services.Extrato(id)
	if err != nil {
		return dto.FiberError(c, fiber.StatusInternalServerError, "Erro ao recuperar as transações")
	}

	response := dto.ExtratoResponse{
		UltimasTransacoes: transacoes,
	}

	response.Saldo.Total = cliente.Saldo
	response.Saldo.Limite = cliente.Limite
	response.Saldo.DataExtrato = time.Now().UTC().Format(time.RFC3339Nano)

	return c.JSON(response)

}
