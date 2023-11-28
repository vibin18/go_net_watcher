package main

import (
	"github.com/gofiber/fiber/v2"
	"go_net_watcher/internal/netwatcher"
)

func NewWebTest() map[string]netwatcher.NetDevices {
	var test *netwatcher.AppConfig
	return test.FinalMap
}

func home(c *fiber.Ctx) error {
	gg := NewWebTest()
	FinalSerializer := make(map[string]netwatcher.NetDevices)
	FinalSerializer = gg
	return c.Render("index", FinalSerializer)
}
