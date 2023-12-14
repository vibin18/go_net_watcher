package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/jessevdk/go-flags"
	"go_net_watcher/internal/database"
	"go_net_watcher/internal/handlers"
	"go_net_watcher/internal/netwatcher"
	"go_net_watcher/opts"
	"log"
	"os"
	"sync"
)

var (
	argparser *flags.Parser
	arg       opts.Params
	lock      sync.Mutex
)

func main() {
	initArgparser()
	database.ConnectDB()
	// Validating Network interface
	myIface, err := validateInterface(arg.Iface)
	if err != nil {
		panic(err)
	}

	// enable template rendering on fiber
	engine := html.New("html", ".html")
	engine.Reload(true)
	engine.Debug(true)

	web := fiber.New(fiber.Config{Views: engine})

	// load global configs
	app := &netwatcher.AppConfig{
		NetworkDeviceMap: make(map[string]string),
		MappedList:       make([]netwatcher.Mapping, 0),
		FinalMap:         []netwatcher.NetDevice{},
		Lock:             &lock,
		ComChan:          make(chan []byte),
	}
	myapp := netwatcher.NewAppConfig(app)
	app.GetConf(arg.MapFile)
	handlers.NewRouteConfigs(app)

	// Start up a scan on each interface.
	go func() {
		if err := myapp.ArpScan(&myIface); err != nil {
			log.Printf("interface %v: %v", myIface.Name, err)
		}
	}()

	web.Get("/", handlers.Home)
	web.Get("/watcher", handlers.Watcher)
	web.Get("/dblist", handlers.DbList)
	web.Get("/update", handlers.Updater)
	log.Fatal(web.Listen(":3000"))

}

func initArgparser() {
	argparser = flags.NewParser(&arg, flags.Default)
	_, err := argparser.Parse()

	// check if there is a parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}
