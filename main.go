package main

import (
	"main/routes/clientes"
	"main/routines"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	jsoniter "github.com/json-iterator/go"
)

func main() {
	// Routines
	routines.DatabaseMigration()
	routines.ClientesMemoryMapping()

	app := fiber.New(fiber.Config{
		JSONEncoder: jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		JSONDecoder: jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		Immutable:   true,
		Prefork:     true,
	})

	// Initialize default config
	app.Use(compress.New())

	// Or extend your config for customization
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	// app.Use(pprof.New())

	app.Get("/clientes/:id/extrato", clientes.Extrato)
	app.Post("/clientes/:id/transacoes", clientes.NovaTransacao)
	app.Listen(":8080")
}
