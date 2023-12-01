package netwatcher

import (
	"sync"
)

type NetDevices struct {
	IP   string `json:"ip"`
	Name string `json:"name"`
	//ID   uint   `json:"id" gorm:"primaryKey"`
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
}

func NewAppConfig(a *AppConfig) *AppConfig {
	return a
}
