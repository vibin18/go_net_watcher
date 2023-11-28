package main

import (
	"github.com/gofiber/fiber/v2"
	"go_net_watcher/internal/netwatcher"
)

func home(c *fiber.Ctx) error {
	lock.Lock()
	defer lock.Unlock()
	var test *netwatcher.AppConfig
	gg := make(map[string]netwatcher.NetDevices)
	gg = test.FinalMap
	return c.Render("index", gg)
}
