package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"go_net_watcher/internal/netwatcher"
	"go_net_watcher/opts"
	"log"
	"os"
	"sync"
	"time"
)

var (
	argparser *flags.Parser
	arg       opts.Params
	wg        sync.WaitGroup
	lock      sync.Mutex
)

func main() {
	initArgparser()

	myIface, err := validateInterface(arg.Iface)
	if err != nil {
		panic(err)
	}

	app := netwatcher.AppConfig{
		NetworkDeviceMap: make(map[string]string),
		MappedList:       make([]netwatcher.Mapping, 0),
		FinalMap:         make(map[string]netwatcher.NetDevices),
		Lock:             &lock,
	}
	app.GetConf(arg.MapFile)

	wg.Add(2)
	// Start up a scan on each interface.
	go func() {
		defer wg.Done()
		start := time.Now()
		if err := app.ArpScan(&myIface); err != nil {
			log.Printf("interface %v: %v", myIface.Name, err)
		}
		et := time.Since(start)
		log.Printf("Scanning took %v seconds", et.Seconds())
	}()

	go func() {
		defer wg.Done()
		start := time.Now()
		for {
			start := time.Now()
			app.Lock.Lock()
			app.MapDevices()
			app.Lock.Unlock()
			et := time.Since(start)
			log.Printf("Mapping took %v seconds", et.Seconds())
		}

	}()

	go func() {
		defer wg.Done()
		for {
			app.Lock.Lock()
			err := PrettyPrint(app.FinalMap)
			if err != nil {
				println(err)
			}
			app.Lock.Unlock()
			time.Sleep(2 * time.Second)

		}

	}()

	wg.Wait()
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
