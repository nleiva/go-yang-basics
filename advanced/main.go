package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

var (
	SchemaTree map[string]*yang.Entry = make(map[string]*yang.Entry)
	ΛEnumTypes map[string][]reflect.Type
)

func init() {
	SchemaTree["emptyBranchTestOne"] = &yang.Entry{
		Name: "empty-branch-test-one",
		Kind: yang.DirectoryEntry,
		Dir: map[string]*yang.Entry{
			"string": {
				Name: "string",
				Kind: yang.LeafEntry,
				Type: &yang.YangType{Kind: yang.Ystring},
			},
		},
	}
}

func String(s string) *string { return &s }

type emptyBranchTestOne struct {
	String *string `path:"string"`
}

func (*emptyBranchTestOne) IsYANGGoStruct() {}

// Validate validates s against the YANG schema corresponding to its type.
func (e *emptyBranchTestOne) ΛValidate(opts ...ygot.ValidationOption) error {
	if err := ytypes.Validate(SchemaTree["emptyBranchTestOne"], e, opts...); err != nil {
		return err
	}
	return nil
}

// Validate validates s against the YANG schema corresponding to its type.
func (e *emptyBranchTestOne) Validate(opts ...ygot.ValidationOption) error {
	return e.ΛValidate(opts...)
}

// ΛEnumTypeMap returns a map, keyed by YANG schema path, of the enumerated types
// that are included in the generated code.
func (e *emptyBranchTestOne) ΛEnumTypeMap() map[string][]reflect.Type { return ΛEnumTypes }

// ΛBelongingModule returns the name of the module that defines the namespace
// of Base_BaseContainer.
func (*emptyBranchTestOne) ΛBelongingModule() string {
	return "base"
}

// Unmarshal unmarshals data, which must be RFC7951 JSON format, into
// destStruct, which must be non-nil and the correct GoStruct type. It returns
// an error if the destStruct is not found in the schema or the data cannot be
// unmarshaled. The supplied options (opts) are used to control the behaviour
// of the unmarshal function - for example, determining whether errors are
// thrown for unknown fields in the input JSON.
func Unmarshal(data []byte, destStruct ygot.GoStruct, opts ...ytypes.UnmarshalOpt) error {
	tn := reflect.TypeOf(destStruct).Elem().Name()
	schema, ok := SchemaTree[tn]
	if !ok {
		return fmt.Errorf("could not find schema for type %s", tn)
	}
	var jsonTree interface{}
	if err := json.Unmarshal([]byte(data), &jsonTree); err != nil {
		return err
	}
	return ytypes.Unmarshal(schema, destStruct, jsonTree, opts...)
}

func main() {
	b := []byte(`{"Number": 1}`)

	var model map[string]interface{}

	/////////
	// TEST 1
	/////////
	err := json.Unmarshal(b, &model)
	if err != nil {
		fmt.Println("error:", err)
	}

	switch value := model["Number"].(type) {
	case int:
		fmt.Printf("\nTEST 1: int: %+v\n", value)
	case float64:
		fmt.Printf("\nTEST 1: float: %+v\n", value)
	default:
		fmt.Printf("\nTEST 1: something else: %+v\n", value)
	}

	/////////
	// TEST 2
	/////////
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	err = d.Decode(&model)
	if err != nil {
		fmt.Println("error:", err)
	}
	switch value := model["Number"].(type) {
	case int:
		fmt.Printf("\nTEST 2: int: %+v\n", value)
	case float64:
		fmt.Printf("\nTEST 2: float: %+v\n", value)
	case json.Number:
		n, err := value.Int64()
		if err != nil {
			break
		}
		fmt.Printf("\nTEST 2: int: %+v\n", n)
	default:
		fmt.Printf("\nTEST 2: something else: %+v\n", value)
	}

	/////////
	// TEST 3
	/////////
	data := map[string]interface{}{"Number": 1}

	b, err = json.Marshal(data)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = json.Unmarshal(b, &model)
	if err != nil {
		fmt.Println("error:", err)
	}

	switch value := model["Number"].(type) {
	case int:
		fmt.Printf("\nTEST 3: int: %+v\n", value)
	case float64:
		fmt.Printf("\nTEST 3: float: %+v\n", value)
	default:
		fmt.Printf("\nTEST 3: something else: %+v\n", value)
	}

	/////////
	// TEST 4
	/////////
	input := `{"string": "hello"}`

	load := &emptyBranchTestOne{}
	if err := Unmarshal([]byte(input), load); err != nil {
		fmt.Printf("Can't unmarshal JSON: %v\n", err)
	}
	fmt.Printf("\nTEST 4: %s\n", *load.String)

	/////////
	// TEST 5
	/////////
	inStruct := &emptyBranchTestOne{
		String: String("goodbye"),
	}

	err = inStruct.Validate()
	if err != nil {
		fmt.Printf("Input is not valid: %v", err)
	}

	json, err := ygot.EmitJSON(inStruct, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: false,
		},
	})

	if err != nil {
		fmt.Printf("can't emit JSON config: %v", err)
	}

	fmt.Printf("\nTEST 5: %s\n", json)

}

// ValidatedGoStruct is an interface implemented by all Go structs (YANG
// container or lists), *except* when the default validate_fn_name generation
// flag is overridden.
// type ValidatedGoStruct interface {
// 	GoStruct
// 	Validate(...ValidationOption) error
// 	ΛEnumTypeMap() map[string][]reflect.Type
// 	ΛBelongingModule() string
// }

// GoStruct is an interface which can be implemented by Go structs that are
// generated to represent a YANG container or list member. It simply allows
// handling code to ensure that it is interacting with a struct that will meet
// the expectations of the interface - such as the fields being tagged with
// appropriate metadata (tags) that allow mapping of the struct into a YANG
// schematree.
// type GoStruct interface {
// 	IsYANGGoStruct()
// }

// validatedGoStruct is an interface used for validating GoStructs.
// This interface is implemented by all Go structs (YANG container or lists),
// regardless of generation flag.
// type validatedGoStruct interface {
// 	GoStruct
// 	ΛValidate(...ValidationOption) error
// }
