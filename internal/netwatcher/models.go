package netwatcher

import (
	"sync"
)

type NetDevice struct {
	MAC  string `json:"mac"`
	IP   string `json:"ip"`
	Name string `json:"name"`
	//ID   uint   `json:"id" gorm:"primaryKey"`
}

type Device struct {
	MAC  string `json:"mac"`
	IP   string `json:"ip"`
	Name string `json:"name"`
	ID   uint   `json:"id" gorm:"primaryKey"`
}

type Mapping struct {
	Mac  string `yaml:"mac"`
	Name string `yaml:"name"`
}

type AppConfig struct {
	NetworkDeviceMap map[string]string
	Lock             *sync.Mutex
	MappedList       []Mapping
	FinalMap         []NetDevice
	ComChan          chan []byte
}

func NewAppConfig(a *AppConfig) *AppConfig {
	return a
}
