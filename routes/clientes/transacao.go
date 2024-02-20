package clientes

import (
	"main/dto"
	"main/entities"
	"main/pkg/memory"
	"main/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

var functionName = "NovaTransacao"

func NovaTransacao(c *fiber.Ctx) error {

	var request dto.TransacaoRequest
	var response dto.TransacaoResponse

	var saldo float64
	var limite float64
	var semLimite bool
	var err error

	clienteID := c.Params("id")

	// Checa no cache em memória da aplicação se o cliente existe
	cache := memory.GetCacheInstance()
	_, found := cache.Get("cliente:" + clienteID)
	if !found {
		return c.Status(fiber.StatusNotFound).
			JSON(&dto.HttpError{
				Message: "cliente não encontrado",
			})
	}

	// Parser do Request
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(&dto.HttpError{
				Message: err.Error(),
			})
	}

	transacao := &entities.Transacao{
		IDCliente:   clienteID,
		Tipo:        request.Tipo,
		Valor:       request.Valor,
		Descricao:   request.Descricao,
		RealizadaEm: time.Now().UTC().Format(time.RFC3339Nano),
	}

	// Operação de Crédito ou Débito
	switch request.Tipo {
	case "c":
		saldo, limite, semLimite, err = services.Credito(*transacao)
	case "d":
		saldo, limite, semLimite, err = services.Debito(*transacao)
	default:
		return c.Status(fiber.StatusBadRequest).
			JSON(&dto.HttpError{
				Message: "Operação inválida",
			})
	}

	if semLimite {
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(&dto.HttpError{
				Message: "cliente sem limite disponível",
			})
	}

	// Verifica Erros durante as operações de crédito ou débito
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(&dto.HttpError{
				Message: err.Error(),
			})

	}

	response.Limite = limite
	response.Saldo = saldo

	return c.JSON(response)

}
