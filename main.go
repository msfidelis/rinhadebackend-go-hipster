package main

import (
	"encoding/json"
	"main/routes/clientes"
	"main/routines"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Routines
	routines.DatabaseMigration()
	routines.ClientesMemoryMapping()

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		Prefork:     true,
	})
	app.Get("/clientes/:id/extrato", clientes.Extrato)
	app.Post("/clientes/:id/transacoes", clientes.NovaTransacao)
	app.Listen(":8080")
}
