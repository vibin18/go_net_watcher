package netwatcher

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
)

func (c *AppConfig) GetConf(file string) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Mapping file %v not found! #%v ", file, err)
	}
	err = yaml.Unmarshal(yamlFile, &c.MappedList)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

}

func (c *AppConfig) MapDevices() {
	for mac, ip := range c.NetworkDeviceMap {
		for _, item := range c.MappedList {
			if mac == item.Mac {
				c.FinalMap[mac] = NetDevices{
					ip,
					item.Name,
				}
				break
			}
			c.FinalMap[mac] = NetDevices{
				ip,
				mac,
			}

		}
	}
}

func (c *AppConfig) AddDevicesToNetworkMap(ip net.IP, mac net.HardwareAddr) {
	if ip == nil {
		log.Printf("Missing IP provide while adding to the nework map")
	}
	if mac == nil {
		log.Printf("Missing MAC while adding to the nework map")
	}
	myipString := fmt.Sprint(ip)
	mymacString := fmt.Sprint(mac)
	c.Lock.Lock()
	c.NetworkDeviceMap[mymacString] = myipString
	c.Lock.Unlock()
}
