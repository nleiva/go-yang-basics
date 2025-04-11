# Go YANG Basics

This tutorial demonstrates how to use Go with [YANG](https://datatracker.ietf.org/doc/html/rfc7950) data models using tools like [`goyang`](https://github.com/openconfig/goyang) and [`ygot`](https://github.com/openconfig/ygot).

## Table of Contents

- [Overview](#overview)
- [1. Define a YANG Model](#1-define-a-yang-model)
- [2. Inspect the Data Tree with goyang (Optional)](#2-inspect-the-data-tree-with-goyang-optional)
- [3. Generate Go Bindings with ygot](#3-generate-go-bindings-with-ygot)
- [4. Create and Populate a YANG Instance](#4-create-and-populate-a-yang-instance)
- [5. Parse a YANG Instance](#5-parse-a-yang-instance)
- [6. Validate Instance Values](#6-validate-instance-values)

---

## Overview

We’ll be working with a simple YANG model (`base.yang`) and show how to:

- Inspect it using `goyang`
- Generate Go code using `ygot`
- Create instances and marshal them to JSON
- Parse JSON back into Go structs

---

## 1. Define a YANG Model

Let’s use the following example: [`base.yang`](base.yang)

```c
module base {
  namespace "urn:mod";
  prefix "myprefix";

  typedef base-type {
    type int32;
  }

  container base-container {
    leaf base-container-leaf-1 { type string; }
    leaf base-container-leaf-2 { type base-type; }
  }
}
```

---

## 2. Inspect the Data Tree with `goyang` (Optional)

Install `goyang`:

```bash
go install github.com/openconfig/goyang@latest
```

Inspect the parsed data model:

```bash
goyang base.yang
```

Expected output:

```ruby
rw: myprefix:base {
  rw: myprefix:base-container {
    rw: string myprefix:base-container-leaf-1
    rw: int32 myprefix:base-container-leaf-2
  }
}
```

---

## 3. Generate Go Bindings with `ygot`

Install the generator:

```bash
go install github.com/openconfig/ygot/generator@latest
```

Generate Go code:

```bash
generator -path=. \
  -output_file=pkg/base.go \
  -enum_suffix_for_simple_union_enums \
  -package_name=test \
  -generate_fakeroot -fakeroot_name=test \
  -generate_getters \
  -generate_ordered_maps=false \
  -generate_simple_unions \
  base.yang
```

This creates Go structs in `pkg/base.go`.

---

## 4. Create and Populate a YANG Instance

Example: [`build/main.go`](build/main.go)

```go
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

	jsonOutput, _ := ygot.EmitJSON(&t, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
	})
	fmt.Println(jsonOutput)
}
```

Run it:

```bash
go run build/main.go
```

Output:

```json
{
  "base-container-leaf-1": "hello",
  "base-container-leaf-2": 1
}
```

---

## 5. Parse a YANG Instance

Example: [`parse/main.go`](parse/main.go)

```go
package main

import (
	"fmt"
	model "github.com/nleiva/yang-gen/pkg"
)

func main() {
	input := `{ "base-container": { "base-container-leaf-1": "hello", "base-container-leaf-2": 1 }}`
	t := model.Test{}

	if err := model.Unmarshal([]byte(input), &t); err != nil {
		fmt.Printf("Can't unmarshal JSON: %v", err)
		return
	}

	fmt.Println("Leaf1:", *t.BaseContainer.BaseContainerLeaf_1)
	fmt.Println("Leaf2:", *t.BaseContainer.BaseContainerLeaf_2)
}
```

Run it:

```bash
go run parse/main.go
```

Output:

```
Leaf1: hello
Leaf2: 1
```

---

## 6. Validate Instance Values

Let's add custom type to the YANG model:

```c
module base {
  // ...
  typedef my-base-int32-type {
    type int32 {
      range "1..4 | 10..20";
    }
  }
  container base-container {
    // ...
    leaf base-container-leaf-3 { type my-base-int32-type; } 
  }
}
```

Now, whether you built an instance with a value outside the specified range (`21`):

```go
func main() {
	// Built example
	t2 := model.Test{}
	base := t2.GetOrCreateBaseContainer()
	base.BaseContainerLeaf_3 = ygot.Int32(21)

	err = t2.Validate()
	if err != nil {
		fmt.Printf("ERROR: Built intance is not valid: %v\n", err)
	}
}
```

Or parsed one with an invalid value (`5`):

```go
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
}
```

The `Validate()` method will warn you about it.


Run it:

```bash
go run validate/main.go 
```

Output:

```
ERROR: Parsed input is not valid: /test/base-container: /test/base-container/base-container-leaf-3: schema "base-container-leaf-3": signed integer value 5 is outside specified ranges
ERROR: Built intance is not valid: /test/base-container: /test/base-container/base-container-leaf-3: schema "base-container-leaf-3": signed integer value 21 is outside specified ranges
```