package main

import (
	"fmt"

	network "github.com/nleiva/go-yang-basics/pkg"
	"github.com/openconfig/ygot/ygot"
)

func main() {
	device := network.Device{}
	iface := device.GetOrCreateInterface()

	// Configure basic interface properties
	iface.Name = ygot.String("eth0")
	iface.Mtu = ygot.Uint16(1500)
	iface.Priority = ygot.Uint8(12)

	// Configure extended properties added via augmentation
	iface.Status = network.NetworkDevice_Interface_Status_up
	iface.Bandwidth = ygot.Uint32(1000) // 1000 Mbps

	// Validate the configuration
	err := device.Validate()
	if err != nil {
		fmt.Printf("ERROR: Built instance is not valid: %v\n", err)
		return
	}

	// Generate JSON output showing all configuration including augmented fields
	jsonOutput, err := ygot.EmitJSON(iface, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: false,
		},
	})

	if err != nil {
		fmt.Printf("Error generating JSON: %v\n", err)
		return
	}

	fmt.Println("Network Interface Configuration (with Augmented Fields):")
	fmt.Printf("%s\n", jsonOutput)

	// Example with custom maintenance status
	fmt.Println("\n=== Example with Custom Status ===")
	device2 := network.Device{}
	iface2 := device2.GetOrCreateInterface()
	iface2.Name = ygot.String("wlan0")
	iface2.Status = network.UnionString("maintenance-scheduled")

	jsonOutput2, _ := ygot.EmitJSON(iface2, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: false,
		},
	})

	fmt.Printf("%s\n", jsonOutput2)
}
