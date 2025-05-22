# Robby.go API Documentation

Complete documentation for all types, functions, and methods in the Robby library.

## Table of Contents

- [Core Types](#core-types)
- [Agent Creation](#agent-creation)
- [Configuration Options](#configuration-options)
- [Chat Methods](#chat-methods)
- [Tool Methods](#tool-methods)
- [MCP Methods](#mcp-methods)
- [Utility Functions](#utility-functions)

---

## Core Types

### `Agent`

The main structure representing an AI agent instance.

```go
type Agent struct {
    ctx       context.Context                        // Context for operations
    dmrClient openai.Client                         // Docker Model Runner client
    Params    openai.ChatCompletionNewParams        // Chat completion parameters
    Tools     []openai.ChatCompletionToolParam      // Available tools
    ToolCalls []openai.ChatCompletionMessageToolCall // Current tool calls
    mcpClient *mcp_golang.Client                    // MCP client instance
    mcpCmd    *exec.Cmd                             // MCP command process
    lastError error                                 // Last error encountered
}
```

**Fields:**
- `ctx`: Context used for all operations and cancellation
- `dmrClient`: OpenAI-compatible client for Docker Model Runner
- `Params`: Chat completion parameters (messages, model, temperature, etc.)
- `Tools`: Array of available tools/functions the agent can use
- `ToolCalls`: Current tool calls detected by the model
- `mcpClient`: Client for Model Context Protocol operations
- `mcpCmd`: Command process for MCP server communication
- `lastError`: Internal error tracking during agent creation

### `AgentOption`

Function type for configuring agent options during creation.

```go
type AgentOption func(*Agent)
```

### `STDIOCommandOption`

Type alias for command arguments used in MCP client setup.

```go
type STDIOCommandOption []string
```

---

## Agent Creation

### `NewAgent(options ...AgentOption) (*Agent, error)`

Creates a new Agent instance with the specified configuration options.

**Parameters:**
- `options`: Variable number of `AgentOption` functions to configure the agent

**Returns:**
- `*Agent`: Configured agent instance
- `error`: Error if agent creation fails

**Example:**
```go
agent, err := NewAgent(
    WithDMRClient(ctx, "http://model-runner.docker.internal/engines/llama.cpp/v1/"),
    WithParams(openai.ChatCompletionNewParams{
        Model: "ai/qwen2.5:latest",
        Messages: []openai.ChatCompletionMessageParamUnion{
            openai.UserMessage("Hello, world!"),
        },
    }),
)
```

**Usage Notes:**
- Apply options in order - later options may override earlier ones
- Returns error if any option fails during application
- Agent is ready to use immediately after successful creation

---

## Configuration Options

### `WithDMRClient(ctx context.Context, baseURL string) AgentOption`

Configures the agent to use Docker Model Runner as the LLM backend.

**Parameters:**
- `ctx`: Context for the client operations
- `baseURL`: Base URL of the Docker Model Runner API endpoint

**Returns:**
- `AgentOption`: Configuration function for agent creation

**Example:**
```go
WithDMRClient(
    context.Background(),
    "http://model-runner.docker.internal/engines/llama.cpp/v1/"
)
```

**Usage Notes:**
- Sets up OpenAI-compatible client with empty API key
- Required for all agent operations
- Context is stored and used for all subsequent operations

### `WithParams(params openai.ChatCompletionNewParams) AgentOption`

Sets the chat completion parameters for the agent.

**Parameters:**
- `params`: OpenAI chat completion parameters including model, messages, temperature, etc.

**Returns:**
- `AgentOption`: Configuration function for agent creation

**Example:**
```go
WithParams(openai.ChatCompletionNewParams{
    Model: "ai/qwen2.5:latest",
    Messages: []openai.ChatCompletionMessageParamUnion{
        openai.SystemMessage("You are a helpful assistant"),
        openai.UserMessage("What is the capital of France?"),
    },
    Temperature: openai.Opt(0.7),
    MaxTokens: openai.Opt(150),
})
```

**Usage Notes:**
- Can be modified after agent creation by directly accessing `agent.Params`
- Messages array is automatically updated during conversations
- Temperature and other parameters affect response generation

### `WithTools(tools []openai.ChatCompletionToolParam) AgentOption`

Adds custom tools/functions that the agent can call.

**Parameters:**
- `tools`: Array of tool definitions following OpenAI function calling format

**Returns:**
- `AgentOption`: Configuration function for agent creation

**Example:**
```go
WithTools([]openai.ChatCompletionToolParam{
    {
        Function: openai.FunctionDefinitionParam{
            Name:        "calculate",
            Description: openai.String("Perform mathematical calculations"),
            Parameters: openai.FunctionParameters{
                "type": "object",
                "properties": map[string]interface{}{
                    "expression": map[string]string{
                        "type": "string",
                        "description": "Mathematical expression to evaluate",
                    },
                },
                "required": []string{"expression"},
            },
        },
    },
})
```

**Usage Notes:**
- Tools are automatically included in chat completion requests
- Must implement corresponding functions in `ExecuteToolCalls`
- Supports parallel tool calling when enabled in parameters

### `WithMCPClient(command STDIOCommandOption) AgentOption`

Configures the agent to use Model Context Protocol (MCP) for external tool access.

**Parameters:**
- `command`: Command and arguments to start the MCP server process

**Returns:**
- `AgentOption`: Configuration function for agent creation

**Example:**
```go
WithMCPClient(WithDockerMCPToolkit())
```

**Usage Notes:**
- Starts external MCP server process
- Sets up STDIO communication channel
- Initializes MCP client for tool discovery and execution
- Process is managed automatically

### `WithMCPTools(tools []string) AgentOption`

Filters and enables specific tools from the MCP server.

**Parameters:**
- `tools`: Array of tool names to enable from available MCP tools

**Returns:**
- `AgentOption`: Configuration function for agent creation

**Example:**
```go
WithMCPTools([]string{"fetch", "brave_web_search", "file_read"})
```

**Usage Notes:**
- Requires `WithMCPClient` to be configured first
- Only specified tools are made available to the agent
- Tool names must match exactly those provided by MCP server

### `WithDockerMCPToolkit() STDIOCommandOption`

Pre-configured command for Docker MCP Toolkit integration.

**Returns:**
- `STDIOCommandOption`: Command configuration for Docker-based MCP toolkit

**Example:**
```go
WithMCPClient(WithDockerMCPToolkit())
```

**Usage Notes:**
- Uses Docker container with socat for communication
- Connects to MCP toolkit running on host.docker.internal:8811
- Requires Docker MCP Toolkit to be running

### `WithSocatMCPToolkit() STDIOCommandOption`

Pre-configured command for direct socat MCP Toolkit integration.

**Returns:**
- `STDIOCommandOption`: Command configuration for socat-based MCP toolkit

**Example:**
```go
WithMCPClient(WithSocatMCPToolkit())
```

**Usage Notes:**
- Uses local socat command for communication
- Requires socat to be installed on the system
- Connects directly to host.docker.internal:8811

---

## Chat Methods

### `ChatCompletion() (string, error)`

Performs a single, non-streaming chat completion request.

**Returns:**
- `string`: Complete response content from the model
- `error`: Error if completion fails

**Example:**
```go
response, err := agent.ChatCompletion()
if err != nil {
    log.Fatal(err)
}
fmt.Println("Response:", response)
```

**Usage Notes:**
- Blocks until complete response is received
- Returns only the text content of the first choice
- Suitable for simple question-answer scenarios

### `ChatCompletionStream(callback func(self *Agent, content string, err error) error) (string, error)`

Performs a streaming chat completion with real-time content delivery.

**Parameters:**
- `callback`: Function called for each content chunk received
  - `self`: Reference to the agent instance
  - `content`: Content chunk received from the stream
  - `err`: Error if chunk processing fails
  - Returns `error` to stop streaming early

**Returns:**
- `string`: Complete accumulated response content
- `error`: Error if streaming fails

**Example:**
```go
response, err := agent.ChatCompletionStream(func(self *robby.Agent, content string, err error) error {
    if err != nil {
        return err
    }
    fmt.Print(content) // Print each chunk as received
    return nil
})
```

**Usage Notes:**
- Provides real-time response streaming
- Callback can return error to abort streaming
- Accumulates and returns complete response
- Better user experience for long responses

---

## Tool Methods

### `ToolsCompletion() ([]openai.ChatCompletionMessageToolCall, error)`

Detects and extracts tool calls from the model's response.

**Returns:**
- `[]openai.ChatCompletionMessageToolCall`: Array of detected tool calls
- `error`: Error if completion fails or no tools detected

**Example:**
```go
toolCalls, err := agent.ToolsCompletion()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Detected %d tool calls\n", len(toolCalls))
```

**Usage Notes:**
- Automatically adds tools to completion request
- Stores detected tool calls in `agent.ToolCalls`
- Returns error if no tool calls are detected
- Required before executing tool calls

### `ExecuteToolCalls(toolsImpl map[string]func(any) (any, error)) ([]string, error)`

Executes detected tool calls using provided implementations.

**Parameters:**
- `toolsImpl`: Map of tool names to implementation functions
  - Key: Tool name (string)
  - Value: Function that takes arguments and returns result and error

**Returns:**
- `[]string`: Array of tool execution results as strings
- `error`: Error if any tool execution fails

**Example:**
```go
results, err := agent.ExecuteToolCalls(map[string]func(any) (any, error){
    "add": func(args any) (any, error) {
        argMap := args.(map[string]any)
        a := argMap["a"].(float64)
        b := argMap["b"].(float64)
        return a + b, nil
    },
    "multiply": func(args any) (any, error) {
        argMap := args.(map[string]any)
        a := argMap["a"].(float64)
        b := argMap["b"].(float64)
        return a * b, nil
    },
})
```

**Usage Notes:**
- Requires `ToolsCompletion()` to be called first
- Arguments are passed as `map[string]any` - type assertion required
- Results are automatically added to conversation history
- All detected tools must have implementations provided

### `ToolCallsToJSON() (string, error)`

Converts current tool calls to formatted JSON string.

**Returns:**
- `string`: JSON representation of tool calls
- `error`: Error if JSON serialization fails

**Example:**
```go
jsonStr, err := agent.ToolCallsToJSON()
if err != nil {
    log.Fatal(err)
}
fmt.Println("Tool calls JSON:", jsonStr)
```

**Usage Notes:**
- Returns "[]" if no tool calls are present
- Includes tool ID, function name, and parsed arguments
- Useful for debugging and logging tool interactions

---

## MCP Methods

### `ExecuteMCPToolCalls() ([]string, error)`

Executes detected tool calls using MCP client.

**Returns:**
- `[]string`: Array of tool execution results from MCP server
- `error`: Error if MCP tool execution fails

**Example:**
```go
results, err := agent.ExecuteMCPToolCalls()
if err != nil {
    log.Fatal(err)
}
for i, result := range results {
    fmt.Printf("Tool %d result: %s\n", i, result)
}
```

**Usage Notes:**
- Requires MCP client to be configured
- Requires `ToolsCompletion()` to be called first
- Automatically handles argument parsing and result formatting
- Results are added to conversation history
- Tools are executed by external MCP server

---

## Utility Functions

### `ToolCallsToJSONString(tools []openai.ChatCompletionMessageToolCall) (string, error)`

Standalone utility to convert tool calls array to JSON string.

**Parameters:**
- `tools`: Array of tool calls to convert

**Returns:**
- `string`: Formatted JSON string representation
- `error`: Error if JSON serialization fails

**Example:**
```go
toolCalls, _ := agent.ToolsCompletion()
jsonStr, err := robby.ToolCallsToJSONString(toolCalls)
if err != nil {
    log.Fatal(err)
}
fmt.Println(jsonStr)
```

**Usage Notes:**
- Independent of agent instance
- Useful for external tool call processing
- Produces indented, readable JSON output
- Parses function arguments into proper JSON objects

### `convertToOpenAITools(tools []mcp_golang.ToolRetType) []openai.ChatCompletionToolParam`

Internal utility that converts MCP tool definitions to OpenAI format.

**Parameters:**
- `tools`: Array of MCP tool definitions

**Returns:**
- `[]openai.ChatCompletionToolParam`: Array of OpenAI-compatible tool definitions

**Usage Notes:**
- Internal function used by MCP integration
- Converts MCP JSON schema to OpenAI function parameters
- Handles tool name, description, and parameter mapping
- Not intended for direct external use

---

## Error Handling

Most functions return errors that should be handled appropriately:

```go
// Basic error handling
result, err := agent.ChatCompletion()
if err != nil {
    log.Printf("Chat completion failed: %v", err)
    return
}

// Error handling with context
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

agent, err := NewAgent(
    WithDMRClient(ctx, baseURL),
    WithParams(params),
)
if err != nil {
    log.Fatalf("Failed to create agent: %v", err)
}
```

## Best Practices

1. **Always handle errors** - All functions return meaningful errors
2. **Use contexts** - Provide appropriate contexts for cancellation and timeouts
3. **Manage conversation history** - Update `agent.Params.Messages` for multi-turn conversations
4. **Resource cleanup** - MCP processes are managed automatically but monitor for resource leaks
5. **Tool validation** - Ensure all tool implementations handle edge cases and invalid inputs

---

## Thread Safety

**Important**: Agent instances are **not thread-safe**. Each goroutine should use its own agent instance or implement proper synchronization when sharing agents across goroutines.