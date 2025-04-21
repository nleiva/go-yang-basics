package main

import (
	model "github.com/nleiva/go-yang-basics/pkg"
	"fmt"
	"github.com/openconfig/ygot/ygot"
)

func main() {
	t := model.Test{}
	base := t.GetOrCreateBaseContainer()

	// String "hello"
	base.BaseContainerLeaf_1 = ygot.String("hello")

	err := t.Validate()
	if err != nil {
		fmt.Printf("ERROR: Built instance is not valid: %v\n", err)
	}

	// String "hell"
	base.BaseContainerLeaf_1 = ygot.String("hell")

	err = t.Validate()
	if err != nil {
		fmt.Printf("ERROR: Built instance is not valid: %v\n", err)
	}

}
