package netwatcher

import (
	"github.com/gofiber/fiber/v2"
	"sync"
)

type NetDevices struct {
	IP   string
	Name string
}

type Mapping struct {
	Mac  string `yaml:"mac"`
	Name string `yaml:"name"`
}

type AppConfig struct {
	NetworkDeviceMap map[string]string
	Lock             *sync.Mutex
	MappedList       []Mapping
	FinalMap         map[string]NetDevices
	Fiber            *fiber.Ctx
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}
