package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go_net_watcher/internal/netwatcher"
)

var app *netwatcher.AppConfig

func NewRouteConfigs(a *netwatcher.AppConfig) {
	app = a
}

func Home(ctx *fiber.Ctx) error {
	app.Lock.Lock()
	defer app.Lock.Unlock()
	amap := make(map[string]netwatcher.NetDevices)
	amap = app.FinalMap

	return ctx.Render("index", amap)
}
