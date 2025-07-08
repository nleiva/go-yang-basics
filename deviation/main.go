package main

import (
	"fmt"

	network "github.com/nleiva/go-yang-basics/pkg"
	"github.com/openconfig/ygot/ygot"
)

func main() {
	device := network.Device{}
	iface := device.GetOrCreateInterface()

	// Example 1: Valid interface name (matches pattern ethX or wlanX)
	fmt.Println("=== Example 1: Valid Interface Name ===")
	iface.Name = ygot.String("eth0")

	err := device.Validate()
	if err != nil {
		fmt.Printf("ERROR: Built instance is not valid: %v\n", err)
	} else {
		fmt.Println("Valid interface name: eth0")
	}

	// Example 2: Another valid interface name
	fmt.Println("\n=== Example 2: Another Valid Interface Name ===")
	iface.Name = ygot.String("wlan1")

	err = device.Validate()
	if err != nil {
		fmt.Printf("ERROR: Built instance is not valid: %v\n", err)
	} else {
		fmt.Println("Valid interface name: wlan1")
	}

	// Example 3: Invalid interface name (doesn't match pattern)
	fmt.Println("\n=== Example 3: Invalid Interface Name ===")
	iface.Name = ygot.String("lo0") // loopback interface - doesn't match ethX or wlanX pattern

	err = device.Validate()
	if err != nil {
		fmt.Printf("ERROR: Built instance is not valid: %v\n", err)
	} else {
		fmt.Println("Interface name is valid (unexpected)")
	}
}
