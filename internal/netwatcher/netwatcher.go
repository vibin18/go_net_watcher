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
			if mac == item.Mac {
				a.FinalMap = append(a.FinalMap, NetDevices{
					mac,
					ip,
					item.Name,
				})
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
