# Use the Sea Flea MCP Server
> https://github.com/sea-monkeys/sea-flea

## Setup

### Build the WASM plugin

```bash
cd mcp-server/plugins/calc
./build.plugin.sh
```
> this will produce `calculator.wasm` at the root of the `plugins` directory

### Build a Docker image to package Sea Flea with the WASM plugin


```bash
cd mcp-server
./build.image.sh
```
> this will produce this image `k33g/mcp-demo:with-agents`

