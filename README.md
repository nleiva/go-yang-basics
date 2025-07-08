# Practical YANG data modeling
[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/nleiva/go-yang-basics?quickstart=1)

In this tutorial, we'll turn a [YANG](https://datatracker.ietf.org/doc/html/rfc7950) data model into Go structs and types that we can use in Go programs for configuration management and validation.

We'll build a Go package that models our network device using YANG. To visualize and verify the structure of our YANG model, we can use [`goyang`](https://github.com/openconfig/goyang), which displays the model as a tree. Once we're confident in the model's structure, we use the [`ygot`](https://github.com/openconfig/ygot) generator to convert the YANG files into Go code, complete with type-safe structs, validation functions, and JSON serialization support. 

[YANG](https://datatracker.ietf.org/doc/html/rfc7950) is a language for modeling configuration and operational data in a tree structure. Think of it as a schema that describes what your data should look like, not the actual data itself. Each part of the model has a name and either contains a value or has child elements. You can set rules to make sure the data is valid. The actual data can be in XML, JSON, Protobuf, or other formats, as long as it follows the YANG model's rules.

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

We'll work with a network device configuration YANG model (`base.yang`) and learn how to:

- Look at it using `goyang`
- Turn it into Go code using `ygot`
- Create instances and convert them to JSON
- Parse JSON back into Go structs
- Check that configurations follow YANG rules

---

## 1. Define a YANG Model

Let's create a YANG model that represents a network device interface with constraints and data types.

In [`base.yang`](base.yang), we define a YANG module that models network interface properties such as name, MTU (Maximum Transmission Unit), and priority levels. This demonstrates how to use YANG's built-in types, create custom types, specify value ranges, and include descriptive documentation.

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
- **`priority-level`**: Network priority with two valid ranges (1-5 for low, 10-15 for high)
- **`interface`**: A container with network interface configuration properties

---

## 2. Inspect the Data Tree with `goyang` (Optional)

Use `goyang` to inspect and visualize the structure of our YANG model to understand the data hierarchy before generating Go code.

`goyang` parses YANG files and displays them in a tree format, showing the namespace, data types, and structure. This helps verify our model is correct and understand how it will be represented.

First make sure `goyang` is present in your environment. If you don't have it installed, you can get it with:

```bash
go install github.com/openconfig/goyang@latest
```

Look at the parsed data model:

```bash
goyang base.yang
```

This reads our YANG file and shows its structure in an easier-to-read tree format.

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

Next, we transform our YANG model into Go structs and types that we can use in Go applications for configuration management and validation.

The `ygot` generator reads YANG files and creates corresponding Go code with proper type safety, validation methods, and JSON marshaling/unmarshaling capabilities. We'll generate a Go package that represents our network device model.

If you don't have it installed, you can get it with:

```bash
go install github.com/openconfig/ygot/generator@latest
```

Generate Go code. Most flags are optional, but we will use some to customize the output:

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

This generates Go structs and getter methods in the `network` package, with a root `Device` container, all derived from [`base.yang`](base.yang).

The generated Go code in `pkg/network.go` includes:
- A `Device` struct as the root container

```go
// Device represents the /device YANG schema element.
type Device struct {
	Interface *NetworkDevice_Interface
}
```

- A `NetworkDevice_Interface` struct for interface configuration


```go
// NetworkDevice_Interface represents the /network-device/interface YANG schema element.
type NetworkDevice_Interface struct {
	Bandwidth *uint32
	Mtu       *uint16
	Name      *string
}
```

- Type-safe fields and validation methods for enforcing YANG constraints


---

## 4. Create and Populate a YANG Instance

Let's write a Go program that constructs the generated structs, populates them with network configuration values, and prints the result as JSON.

We'll create a `Device` instance, set its interface properties (name, MTU, priority) -> ([`build/main.go`](build/main.go))


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

And use `ygot` to convert it to JSON.

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

You can compile and run the build example with `go run build/main.go`, which creates a network device configuration and outputs it as JSON.

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

Now let's look at how to convert JSON configuration data into Go structs using the generated code from our YANG model.

We'll take a JSON string with network interface configuration, convert it into our `Device` struct, and access the values through getter methods. This shows how to go from JSON back to Go structs -> [`parse/main.go`](parse/main.go)

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

We can run the parsing example with `go run parse/main.go`, which takes JSON input and converts it back into Go structs, then shows the parsed values.

Output:

```bash
Interface Name: eth0
MTU: 1500
```

---

## 6. Validate Instance Values

Next, we'll see how `ygot` enforces YANG constraints in Go by validating data and reporting errors when values don't meet the model's requirements.

We'll create examples with both valid and invalid priority values to show how the `Validate()` method enforces the range rules defined in our YANG model (1-5 or 10-15). This shows how YANG's rules become runtime validation in Go.

### Priority Level Validation

Our model defines valid priority levels as ranges `1..5` (low priority) or `10..15` (high priority). We test this validation by creating instances with both valid and invalid priority values -> [`validate/main.go`](validate/main.go)

First, from a JSON input string:

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

And then from a Go struct instance:

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

Run the validation example using `go run validate/main.go` to see how YANG constraints are enforced. Invalid configurations will trigger error messages, demonstrating the model's validation in action.

Output:

```bash
ERROR: ...: schema "priority": value 7 is outside specified ranges
ERROR: ...: schema "priority": value 25 is outside specified ranges
```

## 7. Change a YANG Model

You can modify existing YANG models using [deviation statements](https://datatracker.ietf.org/doc/html/rfc7950#section-5.6.3) without editing the original model files.

YANG deviations let you override or modify parts of an imported model. We'll add a pattern rule to interface names, limiting them to common network interface names (ethX, wlanX). This shows how to adapt existing models to your specific needs.


### Add Interface Name Pattern Restriction

The deviation statement targets a specific leaf in the imported model and replaces its type definition with a more restrictive pattern. This creates a specialized version of the model without changing the original.

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

This restricts valid interface names to `ethX` or `wlanX` patterns. Re-run the code generation, including the [`deviation.yang`](deviation.yang) file:

```bash
generator -path=. \
  # ...
  base.yang \
  deviation.yang
```

This regenerates the Go bindings, including both the base model and the deviation, applying the interface name restrictions to the generated code.

### Validate Pattern Restrictions

After regenerating the Go code with the deviation, we test various interface names to see which ones pass validation. The pattern `eth[0-9]+|wlan[0-9]+` only allows ethernet and wireless interface names with numbers -> [`deviation/main.go`](deviation/main.go)

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

Invoke the `Validate()` method provided by `ygot`.

```go
func main() {
  // ...
  err := device.Validate()
  if err != nil {
    fmt.Printf("ERROR: Built instance is not valid: %v\n", err)
  }
}
```

Run the deviation example with `go run deviation/main.go`, to test both valid interface names (eth0, wlan1) and an invalid one (lo0) to show the pattern rule enforcement.

Output:

```bash
ERROR: ...: schema "name": "lo0" does not match regular expression pattern "^(eth[0-9]+|wlan[0-9]+)$"
```

## 8. Extend a YANG Model

Finally, you can extend existing YANG models by adding new data elements using `augment` statements.

YANG augments let you add new leaves, containers, or other elements to existing models without changing the original. We'll add operational status and bandwidth information to our interface model, showing how to extend configuration models with operational data.


### Add Operational Status and Bandwidth

The `augment` statement extends the existing `/net:interface` container by adding new operational fields: a `status` leaf that demonstrates YANG's union types (combining enumerations and pattern-restricted strings), and a `bandwidth` leaf with a defined numeric range. This highlights YANG's ability to enhance models with advanced type constructs and additional operational data.


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

Re-run the code generation, including the [`augment.yang`](augment.yang) file:

```bash
generator -path=. \
  # ...
  base.yang \
  deviation.yang \
  augment.yang
```

This regenerates the Go code including the base model, deviation restrictions, and augmented fields, creating a complete unified Go interface.

### Use Extended Configuration

After regenerating the Go code with the augmented model, we can now set both the original fields (name, MTU, priority) and the new augmented fields (status, bandwidth). The generated Go code seamlessly combines all fields into a single struct, showing how YANG's modular extensions become unified Go types -> [`augment/main.go`](augment/main.go)

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

Leverage `ygot` JSON emission feature to produce configuration data that conforms to the RFC 7951 JSON encoding standard.

```go
func main() {
  // ...
  jsonOutput, _ := ygot.EmitJSON(iface, &ygot.EmitJSONConfig{Format: ygot.RFC7951})
  fmt.Printf("%s\n", jsonOutput)
}
```

Run the example using `go run augment/main.go` to display the interface configuration, including both the original fields and the newly added operational data (`status` and `bandwidth`).

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
- [Clixon](https://github.com/clicon/clixon): YANG-based configuration manager, with interactive CLI, NETCONF and RESTCONF interfaces, an embedded database and transaction mechanism
- [pyang](https://github.com/mbj4668/pyang): YANG validator and code generator

### OpenConfig Resources  
- [OpenConfig Paths](https://openconfig.net/projects/models/paths/): Standard configuration paths
- [OpenConfig Tree View](https://openconfig.net/projects/models/schemadocs/): Schema documentation
- [ygot](https://github.com/openconfig/ygot): YANG Go Tools library