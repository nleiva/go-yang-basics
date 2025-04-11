package main

import (
	"fmt"
	model "github.com/nleiva/yang-gen/pkg"
	"github.com/openconfig/ygot/ygot"
)

func main() {
	// Parsed example
	input := `{ "base-container": { "base-container-leaf-3": 5 }}`

	t1 := model.Test{}
	if err := model.Unmarshal([]byte(input), &t1); err != nil {
		fmt.Printf("ERROR: Can't unmarshal JSON: %v\n", err)
	}

	err := t1.Validate()
	if err != nil {
		fmt.Printf("ERROR: Parsed input is not valid: %v\n", err)
	}

	// Built example
	t2 := model.Test{}
	base := t2.GetOrCreateBaseContainer()
	base.BaseContainerLeaf_3 = ygot.Int32(21)

	err = t2.Validate()
	if err != nil {
		fmt.Printf("ERROR: Built intance is not valid: %v\n", err)
	}
}
