package main

import (
	"go_net_watcher/internal/netwatcher"
)

func home(c netwatcher.AppConfig) error {
	return c.Fiber.Render("index", c.FinalMap)
}
