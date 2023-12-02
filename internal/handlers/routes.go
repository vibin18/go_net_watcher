package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go_net_watcher/internal/database"
	"go_net_watcher/internal/netwatcher"
)

var app *netwatcher.AppConfig

func NewRouteConfigs(a *netwatcher.AppConfig) {
	app = a
}

func Watcher(ctx *fiber.Ctx) error {
	app.Lock.Lock()
	defer app.Lock.Unlock()
	var amap []netwatcher.NetDevice
	amap = app.FinalMap

	return ctx.Render("list", amap)
}

func DbList(ctx *fiber.Ctx) error {
	ExistingDevices := []netwatcher.Device{}
	database.Database.Db.Find(&ExistingDevices)
	return ctx.Render("list", ExistingDevices)
}

func Home(ctx *fiber.Ctx) error {
	return ctx.Render("index", nil)
}
