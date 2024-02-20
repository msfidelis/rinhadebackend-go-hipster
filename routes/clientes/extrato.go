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

	var response dto.ExtratoResponse

	// Checa no cache em memória da aplicação se o cliente existe
	cache := memory.GetCacheInstance()
	_, found := cache.Get("cliente:" + id)
	if !found {
		return c.Status(fiber.StatusNotFound).
			JSON(&dto.HttpError{
				Message: "cliente não encontrado",
			})
	}

	cliente, err := services.BuscaCliente(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(&dto.HttpError{
				Message: err.Error(),
			})
	}

	transacoes, err := services.Extrato(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(&dto.HttpError{
				Message: err.Error(),
			})
	}

	response.UltimasTransacoes = transacoes
	response.Saldo.Total = cliente.Saldo
	response.Saldo.Limite = cliente.Limite
	response.Saldo.DataExtrato = time.Now().UTC().Format(time.RFC3339Nano)

	return c.JSON(response)

}
