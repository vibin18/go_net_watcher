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

	ExistingDevices := []Device{}
	database.Database.Db.Find(&ExistingDevices)
	//  Check device already exist with an ID
	// loop through existing MAC(devices) and if a device with an ID exist
	for _, dev := range ExistingDevices {
		// Check device has same MAC

		if dev.MAC == device.MAC {
			log.Printf("Device already in DB with MAC : %v", dev.MAC)
			// Device has same MAC
			// Continue to next loop
			// TODO
			// TODO Check if Update works when new mapping is given in the mapping yaml file
			// Device has a different MAC
			// Update
			log.Printf("Ignoring MAC : %v", dev.MAC)
			break
		}
		// Device has a different MAC
		// Update

		log.Printf("Device NOT found in DB with MAC : %v", dev.MAC)
		CreateDeviceToDb(device, a.MappedList)
		//database.Database.Db.Update("ID", &device)
	}

	//  Device not found with above match
	// CreateDeviceToDb(device, a.MappedList)

}

func CreateDeviceToDb(device NetDevice, mappedList []Mapping) {

	c1 := make(chan bool, len(mappedList))
	c2 := make(chan NetDevice, len(mappedList))
	// Check device mac has a mapping available

	for _, md := range mappedList {
		// Device with mapping
		go func(md Mapping) {
			if md.Mac == device.MAC {
				c1 <- true
				c2 <- device
			} else {
				c1 <- false
				c2 <- device
			}
		}(md)
		// Add device to DB with Name mapping
	}

	for i := 0; i < len(mappedList); i++ {
		select {
		case status := <-c1:
			if status {

				myname := <-c2
				log.Printf("Mac to Name found in MAP for mac: %v with name: %v", myname.MAC, myname.Name)
				// Device with mapping
				myDevice := Device{
					MAC:  device.MAC,
					IP:   device.IP,
					Name: myname.Name,
				}
				database.Database.Db.Create(&myDevice)
				break
			}

			myname := <-c2
			log.Printf("Name NOT found in the Map for MAC: %v", myname.MAC, myname.MAC)
			// Device without mapping
			myDevice := Device{
				MAC:  device.MAC,
				IP:   device.IP,
				Name: myname.MAC,
			}
			database.Database.Db.Create(&myDevice)
		}
	}

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
