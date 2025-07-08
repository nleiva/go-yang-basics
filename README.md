# Practical YANG data modeling
[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/nleiva/go-yang-basics?quickstart=1)

This tutorial demonstrates how to use Go with [YANG](https://datatracker.ietf.org/doc/html/rfc7950) data models using tools such as [`goyang`](https://github.com/openconfig/goyang) and [`ygot`](https://github.com/openconfig/ygot).

[YANG](https://datatracker.ietf.org/doc/html/rfc7950) is a modeling language to define configuration and operational state data models in a hierarchical data tree. It expresses the structure of data, not the data itself. Each node in the model has a name and a value or a set of child nodes. Models can describe constraints to enforce on the data. Instances of data can be expressed in XML, JSON, Protobuf, etc. and are considered valid if they adhere to the YANG data model (schema).

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

We'll be working with a practical network device configuration YANG model (`base.yang`) and show how to:

- Inspect it using `goyang`
- Generate Go code using `ygot`
- Create instances and marshal them to JSON
- Parse JSON back into Go structs
- Validate configurations against YANG constraints

---

## 1. Define a YANG Model

Create a practical YANG model that represents a network device interface configuration with realistic constraints and data types.

We'll define a YANG module that models network interface properties ([`base.yang`](base.yang)) including name, MTU (Maximum Transmission Unit), and priority levels. This demonstrates how to use YANG's built-in types, custom typedefs, ranges, and documentation.

```c
module network-device {
  namespace "urn:example:network";
  prefix "net";

  typedef priority-level {
    type uint8 {
      range "1..5 | 10..15";
    }
    description "Network priority levels: 1-5 (low priority) or 10-15 (high priority)";
  }

  container interface {
    description "Network interface configuration";
    
    leaf name {
      type string;
      description "Interface name (e.g., eth0, wlan0)";
    }
    
    leaf mtu {
    type uint16 {
      range "68..9216";
    }
      description "Maximum Transmission Unit in bytes";
    }
    
    leaf priority {
      type priority-level;
      description "Interface priority level";
    }
  }
}
```

This model defines:
- **`priority-level`**: Network priorities with two valid ranges (1-5 for low, 10-15 for high)
- **`interface`**: A container with network interface configuration properties

---

## 2. Inspect the Data Tree with `goyang` (Optional)

Use the goyang tool to inspect and visualize the structure of our YANG model to understand the data hierarchy before generating Go code.

The `goyang` tool parses YANG files and displays them in a tree format, showing the namespace, data types, and structure. This helps verify our model is correct and understand how it will be represented.

Install `goyang`:

```bash
go install github.com/openconfig/goyang@latest
```

Inspect the parsed data model:

```bash
goyang base.yang
```

This parses our YANG file and displays its structure in a human-readable tree format.

Expected output:

```ruby
$ goyang base.yang
rw: net:network-device {
  
  // Network interface configuration
  rw: net:interface {
    
    // Maximum Transmission Unit in bytes
    rw: uint16 net:mtu
    
    // Interface name (e.g., eth0, wlan0)
    rw: string net:name
    
    // Interface priority level
    rw: priority-level net:priority
  }
}
```

---

## 3. Generate Go Bindings with `ygot`

Transform our YANG model into Go structs and types that can be used in Go applications for configuration management and validation.

The `ygot` generator reads YANG files and creates corresponding Go code with proper type safety, validation methods, and JSON marshaling/unmarshaling capabilities. We'll generate a Go package that represents our network device model.

Install the generator:

```bash
go install github.com/openconfig/ygot/generator@latest
```

Generate Go code:

```bash
generator -path=. \
  -output_file=pkg/network.go \
  -enum_suffix_for_simple_union_enums \
  -package_name=network -generate_fakeroot -fakeroot_name=device \
  -generate_getters \
  -generate_ordered_maps=false \
  -generate_simple_unions \
  base.yang
```

This runs the generator with options to create Go structs in the `network` package, generate getter methods, create a root device container, and process our [`base.yang`](base.yang) file.

This creates Go structs in `pkg/network.go` with:
- `Device` struct as the root container
- `NetworkDevice_Interface` struct for interface configuration
- Proper type constraints and validation methods

---

## 4. Create and Populate a YANG Instance

Create a Go program that instantiates our generated structs, populates them with network configuration data, and outputs the result as JSON.

We'll create a `Device` instance, configure its interface properties (name, MTU, priority).

Example: [`build/main.go`](build/main.go)

```go
func main() {
  // Create a new network device instance
  device := network.Device{}
  
  // Configure the network interface
  iface := device.GetOrCreateInterface()
  iface.Name = ygot.String("eth0")
  iface.Mtu = ygot.Uint16(1500)
  iface.Priority = ygot.Uint8(3)
  // ...
}
```

Use ygot's JSON emission capabilities to output RFC 7951 compliant JSON configuration data.

```go
func main() {
  // ...
  // Generate JSON output for the interface configuration
  jsonOutput, err := ygot.EmitJSON(iface, &ygot.EmitJSONConfig{Format: ygot.RFC7951})
  if err != nil {
    fmt.Printf("Error generating JSON: %v\n", err)
    return
  }
  
  fmt.Println(">> Network Interface Configuration:")
  fmt.Printf("%s\n", jsonOutput)
}
```

Compile and run our build example, which creates a network device configuration and outputs it as JSON.

```bash
go run build/main.go
```



Output:

```json
{
  "mtu": 1500,
  "name": "eth0",
  "priority": 3
}
```

---

## 5. Parse a YANG Instance

Demonstrate how to parse JSON configuration data back into Go structs using the generated YANG bindings.

We'll take a JSON string representing network interface configuration, unmarshal it into our `Device` struct.

Example: [`parse/main.go`](parse/main.go)

```go
func main() {
  // Sample JSON input representing a network interface configuration
  input := `{ "interface": { "name": "eth0", "mtu": 1500 }}`

  // Create a new device instance
  device := network.Device{}
  
  // Parse the JSON input into the Go struct
  if err := network.Unmarshal([]byte(input), &device); err != nil {
    fmt.Printf("Error parsing JSON: %v\n", err)
    return
  }
  // ...
}
```

Then, access the parsed values through the generated getter methods. This shows the round-trip capability from JSON to Go structs.

```go
func main() {
  // ...
  // Access the parsed values
  iface := device.GetInterface()
  if iface != nil {
    fmt.Println(">> Parsed Network Interface Configuration:")
    if iface.Name != nil {
      fmt.Printf("Interface Name: %s\n", *iface.Name)
    }
    if iface.Mtu != nil {
      fmt.Printf("MTU: %d\n", *iface.Mtu)
    }
  }
}
```

Run our parsing example, which takes JSON input and converts it back into Go structs, then displays the parsed values:

```bash
go run parse/main.go
```

Output:

```bash
Interface Name: eth0
MTU: 1500
```

---

## 6. Validate Instance Values

Show how YANG constraints are enforced in Go code through ygot's validation mechanisms, catching invalid data before it's used.

We'll create examples with both valid and invalid priority values to demonstrate how the `Validate()` method enforces the range constraints defined in our YANG model (1-5 or 10-15). This shows how YANG's declarative constraints become runtime validation in Go.

### Priority Level Validation

Our model defines valid priority levels as ranges `1..5` (low priority) or `10..15` (high priority).

Example: [`validate/main.go`](validate/main.go)

```go
func main() {
  // Example 1: Invalid priority value (out of range)
  input := `{ "interface": { "priority": 7 }}`

  device1 := network.Device{}
  network.Unmarshal([]byte(input), &device1)

  err := device1.Validate()
  if err != nil {
    fmt.Printf("ERROR: Parsed input is not valid: %v\n", err)
  }
  // ...
}
```

```go
func main() {
  // ...
  // Example 2: Build instance with invalid priority
  device2 := network.Device{}
  iface := device2.GetOrCreateInterface()
  iface.Priority = ygot.Uint8(25) // Invalid: outside valid ranges

  err = device2.Validate()
  if err != nil {
    fmt.Printf("ERROR: Built instance is not valid: %v\n", err)
  }
}
```

Run the validation example, which demonstrates how YANG constraints are enforced by creating invalid configurations and showing the resulting error messages.

```bash
go run validate/main.go 
```

Output:

```bash
ERROR: ...: schema "priority": value 7 is outside specified ranges
ERROR: ...: schema "priority": value 25 is outside specified ranges
```

## 7. Change a YANG Model

Demonstrate how to modify existing YANG models indirectly with `deviation` statements without directly editing the original model files.

YANG deviations allow you to override or modify parts of an imported model. We'll add a pattern constraint to interface names, restricting them to common network interface naming conventions (ethX, wlanX). This shows how to adapt existing models to specific organizational requirements.


### Add Interface Name Pattern Restriction

The deviation statement targets a specific leaf in the imported model and replaces its type definition with a more restrictive pattern. This creates a specialized version of the model without modifying the original.

Let's add a pattern restriction to interface names using the [`deviation` statement](https://datatracker.ietf.org/doc/html/rfc7950#page-39):

```c
module network-device-restrictions {
  namespace "urn:example:network:restrictions";
  prefix "net-restrict";

  import network-device { prefix net; }

  deviation /net:interface/net:name {
    deviate replace {
      type string {
        pattern 'eth[0-9]+|wlan[0-9]+';
      }
    }
    description "Restrict interface names to ethernet (ethX) or wireless (wlanX) patterns";
  }
}
```

This restricts valid interface names to `ethX` or `wlanX` patterns. Re-run the code generation including the [`deviation.yang`](deviation.yang) file:

```bash
generator -path=. \
  # ...
  base.yang \
  deviation.yang
```

This regenerates the Go bindings including both the base model and the deviation, applying the interface name restrictions to the generated code.

### Validate Pattern Restrictions

After regenerating the Go bindings with the deviation, we test various interface names to see which ones pass validation. The pattern `eth[0-9]+|wlan[0-9]+` only allows ethernet and wireless interface names with numbers.

Example: [`deviation/main.go`](deviation/main.go)

```go
func main() {
  device := network.Device{}
  iface := device.GetOrCreateInterface()

  // Valid interface names
  iface.Name = ygot.String("eth0")    // Valid
  iface.Name = ygot.String("wlan1")   // Valid
  
  // Invalid interface name
  iface.Name = ygot.String("lo0")     // Invalid: doesn't match pattern
  // ...
}
```

Call ygot's `Validate()` method. 

```go
func main() {
  // ...
  err := device.Validate()
  if err != nil {
    fmt.Printf("ERROR: Built instance is not valid: %v\n", err)
  }
}
```

Run the deviation example, which tests both valid interface names (eth0, wlan1) and invalid ones (lo0) to demonstrate the pattern constraint enforcement.

```bash
go run deviation/main.go 
```

Output:

```bash
ERROR: ...: schema "name": "lo0" does not match regular expression pattern "^(eth[0-9]+|wlan[0-9]+)$"
```

## 8. Extend a YANG Model

Show how to extend existing YANG models by adding new data elements using `augment` statements, demonstrating YANG's modular design.

YANG augments allow you to add new leaves, containers, or other elements to existing models without modifying the original. We'll add operational status and bandwidth information to our interface model, showing how to extend configuration models with operational data.


### Add Operational Status and Bandwidth

The augment statement targets the existing `/net:interface` container and adds two new leaves: a `status` field using union types (enum + pattern) and a `bandwidth` field with range constraints. This demonstrates YANG's advanced type system including unions and enumerations.

Let's add operational fields to the interface using the [`augment` statement](https://datatracker.ietf.org/doc/html/rfc7950#page-28):

```C
module network-device-extensions {
  namespace "urn:example:network:extensions";
  prefix "net-ext";

  import network-device { prefix net; }

  augment "/net:interface" {
    leaf status {
      description "Interface operational status";
      type union {
        type enumeration {
          enum up { description "Interface is operational"; }
          enum down { description "Interface is not operational"; }
          enum testing { description "Interface is in testing mode"; }
        }
        type string {
          pattern "maintenance-.*";
          description "Custom maintenance status";
        }
      }
    }
    
    leaf bandwidth {
      type uint32 {
        range "1..10000";
      }
      units "Mbps";
      description "Interface bandwidth in Megabits per second";
    }
  }
}
```

This adds:
- **`status`**: Operational status (enum or custom maintenance string)
- **`bandwidth`**: Interface bandwidth with range validation

Re-run the code generation including the [`augment.yang`](augment.yang) file:

```bash
generator -path=. \
  # ...
  base.yang \
  deviation.yang \
  augment.yang
```

This regenerates the Go bindings including the base model, deviation restrictions, and augmented fields, creating a complete unified Go interface.

### Use Extended Configuration

After regenerating the Go bindings with the augmented model, we can now set both the original fields (name, MTU, priority) and the new augmented fields (status, bandwidth). The generated Go code seamlessly combines all fields into a single struct, demonstrating how YANG's modular extensions become unified Go types.

Example: [`augment/main.go`](augment/main.go)

```go
func main() {
  device := network.Device{}
  iface := device.GetOrCreateInterface()
  
  // Configure basic and extended properties
  iface.Name = ygot.String("eth0")
  iface.Mtu = ygot.Uint16(1500)
  iface.Priority = ygot.Uint8(12)
  iface.Status = network.NetworkDevice_Interface_Status_up
  iface.Bandwidth = ygot.Uint32(1000) // 1000 Mbps
  // ...
}
```

Use ygot's JSON emission capabilities to output RFC 7951 compliant JSON configuration data.

```go
func main() {
  // ...
  jsonOutput, _ := ygot.EmitJSON(iface, &ygot.EmitJSONConfig{Format: ygot.RFC7951})
  fmt.Printf("%s\n", jsonOutput)
}
```

Run the example to demonstrate the extended interface configuration with both original fields and augmented operational data (status, bandwidth).

```bash
go run augment/main.go 
```

Output:

```json
{
  "bandwidth": 1000,
  "mtu": 1500,
  "name": "eth0", 
  "priority": 12,
  "status": "up"
}
```

## Appendix
- [RFC 7950](https://datatracker.ietf.org/doc/html/rfc7950): The YANG 1.1 Data Modeling Language
- [RFC 7951](https://datatracker.ietf.org/doc/html/rfc7951): JSON Encoding of Data Modeled with YANG
- [OpenConfig](https://openconfig.net/): Vendor-neutral configuration and telemetry standards for network devices

### YANG Tools & Validators
- [YANG Explorer (Nokia)](https://yang.srlinux.dev): Interactive YANG model explorer
- [YANG Data Model Explorer (Juniper)](https://apps.juniper.net/ydm-explorer/): Juniper's YANG browser
- [YANG Catalog (IETF/Cisco)](https://www.yangcatalog.org/yang-search): Comprehensive YANG model repository
- [Clixon](https://github.com/clicon/clixon): YANG-based configuration manager, with interactive CLI, NETCONF and RESTCONF interfaces, an embedded database and transaction mechanism.
- [pyang](https://github.com/mbj4668/pyang): YANG validator and code generator

### OpenConfig Resources  
- [OpenConfig Paths](https://openconfig.net/projects/models/paths/): Standard configuration paths
- [OpenConfig Tree View](https://openconfig.net/projects/models/schemadocs/):Schema documentation
- [ygot](https://github.com/openconfig/ygot): YANG Go Tools library