package main

import (
	"fmt"
	model "github.com/nleiva/go-yang-basics/pkg"
)

func main() {

	input := `{ "base-container": { "base-container-leaf-1": "hello", "base-container-leaf-2": 1 }}`

	t := model.Test{}
	if err := model.Unmarshal([]byte(input), &t); err != nil {
		fmt.Printf("Can't unmarshal JSON: %v", err)
	}

	base := t.GetBaseContainer()
	fmt.Printf("Leaf-1: %s\nLeaf-2: %d\n", *base.BaseContainerLeaf_1, *base.BaseContainerLeaf_2)
}
