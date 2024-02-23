package clientes

import (
	"fmt"
	"main/dto"
	"main/entities"
	"main/pkg/memory"
	"main/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

var functionName = "NovaTransacao"

var request dto.TransacaoRequest

var saldo float64
var limite float64
var semLimite bool
var err error

func NovaTransacao(c *fiber.Ctx) error {

	id := c.Params("id")

	// Checa no cache em memória da aplicação se o cliente existe
	cache := memory.GetCacheInstance()
	_, found := cache.Get("cliente:" + id)
	if !found {
		return dto.FiberError(c, fiber.StatusNotFound, "cliente não encontrado")
	}

	// Parser do Request
	if err := c.BodyParser(&request); err != nil {
		return dto.FiberError(c, fiber.StatusBadRequest, err.Error())
	}

	transacao := &entities.Transacao{
		IDCliente:   id,
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
		return dto.FiberError(c, fiber.StatusBadRequest, "tipo de operação inválida")
	}

	if semLimite {
		return dto.FiberError(c, fiber.StatusUnprocessableEntity, "cliente sem limite disponível")
	}

	// Verifica Erros durante as operações de crédito ou débito
	if err != nil {
		return dto.FiberError(c, fiber.StatusInternalServerError, err.Error())
	}

	fmt.Printf("[%v] Cliente: %v\n", c.Context().ID(), id)
	fmt.Printf("[%v] Tipo: %v\n", c.Context().ID(), request.Tipo)
	fmt.Printf("[%v] Valor: %v\n", c.Context().ID(), request.Valor)
	fmt.Printf("[%v] Descricao: %v\n", c.Context().ID(), request.Descricao)
	fmt.Printf("[%v] Saldo: %v\n", c.Context().ID(), saldo)
	fmt.Printf("[%v] Limite: %v\n", c.Context().ID(), limite)

	response := dto.TransacaoResponse{
		Limite: limite,
		Saldo:  saldo,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
