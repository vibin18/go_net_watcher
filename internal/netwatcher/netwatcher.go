package netwatcher

import (
	"fmt"
	"go_net_watcher/internal/database"
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

					dmap := NetDevice{
						mac,
						ip,
						item.Name,
					}

					a.FinalMap = append(a.FinalMap, dmap)
					break
				}
				a.FinalMap = append(a.FinalMap, NetDevice{
					mac,
					ip,
					mac,
				})

			}
		}
	}
}

func (a *AppConfig) AddDeviceToDb(ip net.IP, mac net.HardwareAddr) {
	if ip == nil {
		log.Printf("Missing IP!")
	}
	if mac == nil {
		log.Printf("Missing MAC!")
	}
	myipString := fmt.Sprint(ip)
	mymacString := fmt.Sprint(mac)
	device := NetDevice{
		mymacString,
		myipString,
		mymacString,
	}
	log.Printf("Trying to add device with MAC : %v", mymacString)

	//  Check device already exist with an ID
	// loop through existing MAC(devices) and if a device with an ID exist
	if CheckDeviceExist(device) {
		log.Printf("Ignoring MAC : %v", device.MAC)
		return
	}
	// TODO
	// Device has a different MAC
	// Update
	// database.Database.Db.Update("ID", &device)

	log.Printf("Device NOT found in DB with MAC : %v", device.MAC)
	log.Printf("Adding device in DB with MAC : %v", device.MAC)
	CreateDeviceToDb(device, a.MappedList)
	//  Device not found with above match
	// CreateDeviceToDb(device, a.MappedList)

}

func CheckDeviceExist(device NetDevice) bool {
	log.Printf("Fetching existing device list from db")
	ExistingDevices := []Device{}
	database.Database.Db.Find(&ExistingDevices)
	for _, dev := range ExistingDevices {
		// Check device has same MAC

		//log.Printf("Checking for MAC: %v with db mac: %v", device.MAC, dev.MAC)
		if device.MAC == dev.MAC {
			//log.Printf("Device already in DB with MAC : %v", dev.MAC)
			// Device has same MAC
			// Continue to next loop

			// TODO Check if Update works when new mapping is given in the mapping yaml file
			// Device has a different MAC
			// Update

			return true
		}
	}
	return false
}

func CreateDeviceToDb(device NetDevice, mappedList []Mapping) {

	// Check device mac has a mapping available

	fmac := IFMapExist(device.MAC, mappedList)

	myDevice := Device{
		MAC:  device.MAC,
		IP:   device.IP,
		Name: fmac,
	}
	database.Database.Db.Create(&myDevice)

	// Device with no mapping
	// Add device to DB with MAC mapping
}

func IFExist(device string, devices []NetDevice) bool {
	c1 := make(chan bool, len(devices))
	for _, dev := range devices {
		go func(dev NetDevice) {
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

func IFMapExist(device string, devices []Mapping) string {
	c1 := make(chan bool, len(devices))
	c2 := make(chan string, len(devices))
	for _, dev := range devices {
		go func(dev Mapping) {
			if dev.Mac == device {
				c1 <- true
				c2 <- dev.Name
			} else {
				c1 <- false
				c2 <- dev.Mac
			}
		}(dev)
	}

	for i := 0; i < len(devices); i++ {
		select {
		case status := <-c1:

			if status {
				mapped := <-c2
				return mapped
			}
			break
		}
	}
	return device

}

func (a *AppConfig) AddDevicesToNetworkMap(ip net.IP, mac net.HardwareAddr) {
	if ip == nil {
		log.Printf("Missing IP!")
	}
	if mac == nil {
		log.Printf("Missing MAC!")
	}
	myipString := fmt.Sprint(ip)
	mymacString := fmt.Sprint(mac)
	a.Lock.Lock()
	a.NetworkDeviceMap[mymacString] = myipString
	a.Lock.Unlock()
}
