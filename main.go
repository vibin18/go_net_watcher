package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/jessevdk/go-flags"
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

	myIface, err := validateInterface(arg.Iface)
	if err != nil {
		panic(err)
	}

	engine := html.New("html", ".html")
	engine.Reload(true)
	engine.Debug(true)

	web := fiber.New(fiber.Config{Views: engine})

	app := netwatcher.AppConfig{
		NetworkDeviceMap: make(map[string]string),
		MappedList:       make([]netwatcher.Mapping, 0),
		FinalMap:         make(map[string]netwatcher.NetDevices),
		Lock:             &lock,
	}
	app.GetConf(arg.MapFile)

	// Start up a scan on each interface.
	go func() {
		if err := app.ArpScan(&myIface); err != nil {
			log.Printf("interface %v: %v", myIface.Name, err)
		}
	}()

	go func() {
		for {
			app.Lock.Lock()
			app.MapDevices()
			app.Lock.Unlock()
		}

	}()

	web.Get("/", home)
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
