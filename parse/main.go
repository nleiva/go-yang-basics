package main

import (
	"fmt"

	network "github.com/nleiva/go-yang-basics/pkg"
)

func main() {
	// Sample JSON input representing a network interface configuration
	input := `{ "interface": { "name": "eth0", "mtu": 1500 }}`

	// Create a new device instance
	device := network.Device{}

	// Parse the JSON input into the Go struct
	if err := network.Unmarshal([]byte(input), &device); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Access the parsed values
	iface := device.GetInterface()
	if iface == nil {
		fmt.Println("No interface configuration found")
		return
	}

	fmt.Println(">> Parsed Network Interface Configuration:")
	if iface.Name != nil {
		fmt.Printf("Interface Name: %s\n", *iface.Name)
	}
	if iface.Mtu != nil {
		fmt.Printf("MTU: %d\n", *iface.Mtu)
	}
}
