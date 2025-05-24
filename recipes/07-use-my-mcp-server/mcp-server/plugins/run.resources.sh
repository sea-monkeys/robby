#!/bin/bash
echo "=== Testing Resources Information ==="
extism call calculator.wasm resources_information \
  --log-level "info" \
  --wasi
echo ""



