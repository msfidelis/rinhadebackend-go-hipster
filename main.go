package main

import (
	"main/routes/clientes"
	"main/routines"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

func main() {
	// Routines
	routines.DatabaseMigration()
	routines.ClientesMemoryMapping()

	app := fiber.New(fiber.Config{
		JSONEncoder: jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		JSONDecoder: jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		Prefork:     true,
	})
	app.Get("/clientes/:id/extrato", clientes.Extrato)
	app.Post("/clientes/:id/transacoes", clientes.NovaTransacao)
	app.Listen(":8080")
}
