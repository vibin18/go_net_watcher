package main

import (
	"github.com/gofiber/fiber/v2"
	"go_net_watcher/internal/netwatcher"
)

var test *netwatcher.AppConfig

func NewWebTest() *netwatcher.AppConfig {
	return test
}

func home(c *fiber.Ctx) error {
	gg := NewWebTest()
	return c.Render("index", gg.FinalMap)
}
