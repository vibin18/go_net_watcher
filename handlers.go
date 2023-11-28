package main

import (
	"github.com/gofiber/fiber/v2"
	"go_net_watcher/internal/netwatcher"
)

func home(c *fiber.Ctx) error {
	var test *netwatcher.AppConfig
	gg := test.FinalMap
	FinalSerializer := make(map[string]netwatcher.NetDevices)
	FinalSerializer = gg
	return c.Render("index", FinalSerializer)
}
