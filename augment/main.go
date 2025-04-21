package main

import (
	model "github.com/nleiva/go-yang-basics/pkg"
	"fmt"
	"github.com/openconfig/ygot/ygot"
)

func main() {
	t := model.Test{}
	base := t.GetOrCreateBaseContainer()
	base.BaseContainerLeaf_1 = ygot.String("hello")
	base.BaseContainerLeaf_2 = ygot.Int32(1)
	base.BaseContainerLeaf_4 = model.UnionString("$goodbye")

	err := t.Validate()
	if err != nil {
		fmt.Printf("ERROR: Built instance is not valid: %v\n", err)
	}

	jsonOutput, err := ygot.EmitJSON(base, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: false,
		},
	})

	if err != nil {
		fmt.Printf("can't emit JSON config: %v", err)
	}
	fmt.Printf("%s\n", jsonOutput)
}
