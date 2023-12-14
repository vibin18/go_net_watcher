package handlers

import (
	"bufio"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go_net_watcher/internal/database"
	"go_net_watcher/internal/netwatcher"
	"log"
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

func Updater(ctx *fiber.Ctx) error {
	log.Println("Received SSE message request")
	//ctx.App().Use(cors.New(cors.Config{
	//	AllowOrigins:  "*",
	//	ExposeHeaders: "Content-Type",
	//}))
	ctx.Set("Access-Control-Allow-Origin", "*")
	ctx.Set("Access-Control-Expose-Headers", "Content-Type")
	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")
	ctx.Set("Transfer-Encoding", "chunked")

	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
	loop:
		for {
			select {
			case data := <-app.ComChan:
				log.Println("Message received")
				mydata := fmt.Sprintf("event: sse1\ndata: %s\n\n", data)
				fmt.Fprintf(w, "data: %s\n\n", mydata)
				err := w.Flush()
				if err != nil {
					// Refreshing page in web browser will establish a new
					// SSE connection, but only (the last) one is alive, so
					// dead connections must be closed here.
					log.Printf("Error while flushing: %v. Closing http connection.\n", err)

					break loop

				}

				//case <-ctx.Context().Done():
				//	log.Println("SSE breaking")
				//	break loop
			}
		}
	})
	return nil
}

func DbList(ctx *fiber.Ctx) error {
	ExistingDevices := []netwatcher.Device{}
	database.Database.Db.Find(&ExistingDevices)
	return ctx.Render("list", ExistingDevices)
}

func Home(ctx *fiber.Ctx) error {
	return ctx.Render("index", nil)
}
