package netwatcher

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
)

func (a *AppConfig) GetConf(file string) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Mapping file %v not found! #%v ", file, err)
	}
	err = yaml.Unmarshal(yamlFile, &a.MappedList)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

}

func (a *AppConfig) MapDevices() {
	for mac, ip := range a.NetworkDeviceMap {
		for _, item := range a.MappedList {

			if !IFExist(mac, a.FinalMap) {
				if mac == item.Mac {

					dmap := NetDevices{
						mac,
						ip,
						item.Name,
					}

					a.FinalMap = append(a.FinalMap, dmap)
					break
				}
				a.FinalMap = append(a.FinalMap, NetDevices{
					mac,
					ip,
					mac,
				})

			}
		}
	}
}

func IFExist(device string, devices []NetDevices) bool {
	c1 := make(chan bool, len(devices))
	for _, dev := range devices {
		go func(dev NetDevices) {
			if dev.MAC == device {
				c1 <- true
			} else {
				c1 <- false
			}
		}(dev)
	}

	for i := 0; i < len(devices); i++ {
		select {
		case status := <-c1:
			if status {
				return true
			}
			break
		}
	}
	return false

}

func (a *AppConfig) AddDevicesToNetworkMap(ip net.IP, mac net.HardwareAddr) {
	if ip == nil {
		log.Printf("Missing IP provide while adding to the nework map")
	}
	if mac == nil {
		log.Printf("Missing MAC while adding to the nework map")
	}
	myipString := fmt.Sprint(ip)
	mymacString := fmt.Sprint(mac)
	a.Lock.Lock()
	a.NetworkDeviceMap[mymacString] = myipString
	a.Lock.Unlock()
}
