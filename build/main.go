package main

import (
	"fmt"
	network "github.com/nleiva/go-yang-basics/pkg"
	"github.com/openconfig/ygot/ygot"
)

func main() {
	// Create a new network device instance
	device := network.Device{}

	// Configure the network interface
	iface := device.GetOrCreateInterface()
	iface.Name = ygot.String("eth0")
	iface.Mtu = ygot.Uint16(1500)
	iface.Priority = ygot.Uint8(3)

	// Generate JSON output for the interface configuration
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

	fmt.Println(">> Network Interface Configuration:")
	fmt.Printf("%s\n", jsonOutput)
}
