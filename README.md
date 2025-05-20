# Practical YANG data modeling
[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/nleiva/go-yang-basics?quickstart=1)

This tutorial demonstrates how to use with [YANG](https://datatracker.ietf.org/doc/html/rfc7950) data models with tools like [`goyang`](https://github.com/openconfig/goyang) and [`ygot`](https://github.com/openconfig/ygot).

## Table of Contents

- [Overview](#overview)
- [1. Define a YANG Model](#1-define-a-yang-model)
- [2. Inspect the Data Tree with goyang (Optional)](#2-inspect-the-data-tree-with-goyang-optional)
- [3. Generate Go Bindings with ygot](#3-generate-go-bindings-with-ygot)
- [4. Create and Populate a YANG Instance](#4-create-and-populate-a-yang-instance)
- [5. Parse a YANG Instance](#5-parse-a-yang-instance)
- [6. Validate Instance Values](#6-validate-instance-values)
- [7. Change a YANG Model](#7-change-a-yang-model)
- [8. Extend a YANG Model](#8-extend-a-yang-model)

---

## Overview

We’ll be working with a simple YANG model (`base.yang`) and show how to:

- Inspect it using `goyang`
- Generate Go code using `ygot`
- Create instances and marshal them to JSON
- Parse JSON back into Go structs

---

## 0. Where to run

Alternatives to run the examples:

1. [GitHub Codespaces](https://codespaces.new/nleiva/go-yang-basics?quickstart=1) by clicking the badge "Open in GutHub Codespaces", or 
2. Inside a container: `docker pull ghcr.io/nleiva/practical-yang:latest`, or
3. Any environment with [Go installed](https://go.dev/doc/install). 

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
  model "bbgithub.dev.bloomberg.com/nleiva2/go-yang-basics/pkg"
  "github.com/openconfig/ygot/ygot"
)

func main() {
  t := model.Test{}
  base := t.GetOrCreateBaseContainer()
  base.BaseContainerLeaf_1 = ygot.String("hello")
  base.BaseContainerLeaf_2 = ygot.Int32(1)

  jsonOutput, _ := ygot.EmitJSON(&t, &ygot.EmitJSONConfig{Format: ygot.RFC7951})
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
  model "bbgithub.dev.bloomberg.com/nleiva2/go-yang-basics/pkg"
)

func main() {
  input := `{ "base-container": { "base-container-leaf-1": "hello", "base-container-leaf-2": 1 }}`
 
  t := model.Test{}
  model.Unmarshal([]byte(input), &t)

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

You can validate values against constraints defined in your YANG model using the `Validate()` method provided by `ygot`.

### Add a Range Restriction

Let’s define a new typedef in `base.yang` using the [`range` statement](https://datatracker.ietf.org/doc/html/rfc7950#section-9.2.4):

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

This means valid values for `base-container-leaf-3` must be between `1..4` or `10..20`.

### Example 1: Invalid Parsed Input

Parse a model instance with an invalid value (`5`):

```go
func main() {
  // Parsed example
  input := `{ "base-container": { "base-container-leaf-3": 5 }}`

  t1 := model.Test{}
  model.Unmarshal([]byte(input), &t1)

  err := t1.Validate()
  if err != nil {
    fmt.Printf("ERROR: Parsed input is not valid: %v\n", err)
  }
}
```

### Example 2: Invalid Built Instance

Build a model instance with a value outside the specified range (`21`):

```go
func main() {
  // Built example
  t2 := model.Test{}
  base := t2.GetOrCreateBaseContainer()
  base.BaseContainerLeaf_3 = ygot.Int32(21)

  err = t2.Validate()
  if err != nil {
    fmt.Printf("ERROR: Built instance is not valid: %v\n", err)
  }
}
```

Run it:

```bash
go run validate/main.go 
```

Output:

```bash
ERROR: Parsed input is not valid: /test/base-container: /test/base-container/base-container-leaf-3: schema "base-container-leaf-3": signed integer value 5 is outside specified ranges
ERROR: Built intance is not valid: /test/base-container: /test/base-container/base-container-leaf-3: schema "base-container-leaf-3": signed integer value 21 is outside specified ranges
```

## 7. Change a YANG Model

You can indirectly change a YANG model with a `deviation` statements.

### Change a datatype in a deviation YANG model.

Let’s add a pattern restriction to the base model using the [`deviation` statement](https://datatracker.ietf.org/doc/html/rfc7950#page-39):

```c
module base-dev {
  namespace "urn:dev";
  prefix "my-dev";

  import base { prefix myprefix; }

  deviation /myprefix:base-container/myprefix:base-container-leaf-1 {
    deviate replace {
      type string {
        pattern 'h.*o';
      }
    }
  }
}
```

This means valid values for `base-container-leaf-1` must start with `h` and finish with an `o`. Re-run the code generation CLI command including the [`deviation.yang`](deviation.yang) file:

```bash
generator -path=. \
  # ...
  base.yang \
  deviation.yang
```

### Validate changes are enforced

Create a model instance with an invalid value (`hell`):

```go
func main() {
  t := model.Test{}
  base := t.GetOrCreateBaseContainer()

  // String "hell"
  base.BaseContainerLeaf_1 = ygot.String("hell")

  err = t.Validate()
  if err != nil {
    fmt.Printf("ERROR: Built instance is not valid: %v\n", err)
  }
}
```

Run it:

```bash
go run deviation/main.go 
```

Output:

```bash
ERROR: Built instance is not valid: /test/base-container: /test/base-container/base-container-leaf-1: schema "base-container-leaf-1": "hell" does not match regular expression pattern "^(h.*o)$"
```

## 8. Extend a YANG Model

You can also add more items to a YANG data model with an `augment` statement.

### Add a new item in an augment YANG model.

Let’s add a new leaf to the base model using the [`augment` statement](https://datatracker.ietf.org/doc/html/rfc7950#page-28):

```c
module base-aug {
  namespace "urn:aug";
  prefix "my-aug";

  import base { prefix myprefix; }

  augment "/myprefix:base-container" {
    leaf base-container-leaf-4 {
      description "Another leaf";
      type union {
        type string {
          pattern "<.*>|$.*";
        }
        type uint32 {
          range "1 .. 1000";
        }
      }
    }
  }
}
```

This means we now have a `base-container-leaf-4` that can be either a string between `<` and `>`, or a string that starts with `$`, or an uint32 between `0` and `1000`. Re-run the code generation CLI command including the [`augment.yang`](augment.yang) file:

```bash
generator -path=. \
  # ...
  base.yang \
  deviation.yang \
  augment.yang
```

### Validate changes are enforced

Create a model instance with a value for the new leaf (`$goodbye`):

```go
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

  jsonOutput, _ := ygot.EmitJSON(&t, &ygot.EmitJSONConfig{Format: ygot.RFC7951})
  fmt.Println(jsonOutput)
}
```

Run it:

```bash
go run augment/main.go 
```

Output:

```bash
{
  "base-container-leaf-1": "hello",
  "base-container-leaf-2": 1,
  "base-container-leaf-4": "$goodbye"
}
```