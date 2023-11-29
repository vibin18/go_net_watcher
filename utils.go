package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"net"
)

func getAllInterfaces() []net.Interface {
	allIfaces, err := net.Interfaces()
	if err != nil {
		log.Panicf("Error validating network interface %v. \n%v", allIfaces, err)
	}
	return allIfaces
}

func validateInterface(iface string) (net.Interface, error) {
	allIfaces := getAllInterfaces()
	for _, i := range allIfaces {
		log.Warn("Checking interface : ", i)
		if iface == i.Name {
			return i, nil
		}
	}
	return net.Interface{}, errors.New("Interface not found ")
}

// PrettyPrint is only used for debugging
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
