package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go_net_watcher/internal/netwatcher"
)

var app *netwatcher.AppConfig

func NewRouteConfigs(a *netwatcher.AppConfig) {
	app = a
}

func Watcher(ctx *fiber.Ctx) error {
	app.Lock.Lock()
	defer app.Lock.Unlock()
	var amap []netwatcher.NetDevices
	amap = app.FinalMap

	return ctx.Render("list", amap)
}

func Home(ctx *fiber.Ctx) error {
	return ctx.Render("index", nil)
}
