#!/bin/bash
tinygo build -scheduler=none --no-debug \
  -o ../calculator.wasm \
  -target wasi .
