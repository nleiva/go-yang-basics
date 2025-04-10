package main

import (
	"fmt"
	model "github.com/nleiva/yang-gen/pkg"
	"github.com/openconfig/ygot/ygot"
)

func main() {
	t := model.Test{}
	base := t.GetOrCreateBaseContainer()
	base.BaseContainerLeaf_1 = ygot.String("hello")
	base.BaseContainerLeaf_2 = ygot.Int32(1)

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
