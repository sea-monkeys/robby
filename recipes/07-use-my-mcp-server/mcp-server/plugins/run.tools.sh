#!/bin/bash

echo "=== Testing Tools Information ==="
extism call calculator.wasm tools_information \
  --log-level "info" \
  --wasi
echo ""

echo "=== Testing Add Tool ==="
extism call calculator.wasm  add \
  --input '{"a":10, "b":32}' \
  --log-level "info" \
  --wasi
echo ""

echo "=== Testing Subtract Tool ==="
extism call calculator.wasm  subtract \
  --input '{"a":50, "b":8}' \
  --log-level "info" \
  --wasi
echo ""

echo "=== Testing Multiply Tool ==="
extism call calculator.wasm  multiply \
  --input '{"a":2, "b":21}' \
  --log-level "info" \
  --wasi
echo ""

echo "=== Testing Divide Tool ==="
extism call calculator.wasm  divide \
  --input '{"a":84, "b":2}' \
  --log-level "info" \
  --wasi
echo ""

