package main

import (
	"fmt"

	network "github.com/nleiva/go-yang-basics/pkg"
	"github.com/openconfig/ygot/ygot"
)

func main() {
	// Example 1: Parse JSON with invalid priority value (out of range)
	fmt.Println("=== Example 1: Invalid Parsed Input ===")
	input := `{ "interface": { "priority": 7 }}`

	device1 := network.Device{}
	if err := network.Unmarshal([]byte(input), &device1); err != nil {
		fmt.Printf("ERROR: Can't unmarshal JSON: %v\n", err)
	}

	err := device1.Validate()
	if err != nil {
		fmt.Printf("ERROR: Parsed input is not valid: %v\n", err)
	} else {
		fmt.Println("Parsed input is valid!")
	}

	// Example 2: Build instance with invalid priority value (out of range)
	fmt.Println("\n=== Example 2: Invalid Built Instance ===")
	device2 := network.Device{}
	iface := device2.GetOrCreateInterface()
	iface.Priority = ygot.Uint8(25) // Invalid: should be 1-5 or 10-15

	err = device2.Validate()
	if err != nil {
		fmt.Printf("ERROR: Built instance is not valid: %v\n", err)
	} else {
		fmt.Println("Built instance is valid!")
	}

	// Example 3: Valid configuration
	fmt.Println("\n=== Example 3: Valid Configuration ===")
	device3 := network.Device{}
	iface3 := device3.GetOrCreateInterface()
	iface3.Name = ygot.String("eth0")
	iface3.Mtu = ygot.Uint16(1500)
	iface3.Priority = ygot.Uint8(12) // Valid: within 10-15 range

	err = device3.Validate()
	if err != nil {
		fmt.Printf("ERROR: Configuration is not valid: %v\n", err)
	} else {
		fmt.Println("Valid configuration created successfully!")
	}
}
