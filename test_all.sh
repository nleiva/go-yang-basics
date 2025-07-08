#!/bin/bash

echo "=========================================="
echo "Go YANG Basics - Testing All Examples"
echo "=========================================="

echo ""
echo "0. Install dependencies:"
echo "--------------------------------------"
go install github.com/openconfig/goyang@latest
go install github.com/openconfig/ygot/generator@latest
./generate.sh

echo ""
echo "1. Inspecting YANG model with goyang:"
echo "--------------------------------------"
goyang base.yang

echo ""
echo "2. Building and configuring network interface:"
echo "-----------------------------------------------"
go run build/main.go

echo ""
echo "3. Parsing JSON into Go structs:"
echo "---------------------------------"
go run parse/main.go

echo ""
echo "4. Validating YANG constraints:"
echo "--------------------------------"
go run validate/main.go

echo ""
echo "5. Testing deviation constraints:"
echo "---------------------------------"
go run deviation/main.go

echo ""
echo "6. Using augmented model:"
echo "-------------------------"
go run augment/main.go

echo ""
echo "=========================================="
echo "All examples completed successfully!"
echo "=========================================="
