# Robby Recipes 🧑‍🍳

This directory contains practical examples demonstrating how to use **Robby** to build AI agents with different capabilities. Each recipe is a self-contained example that showcases specific features and patterns.

## Prerequisites

Before running these recipes, make sure you have:

1. **Docker Model Runner** running and accessible at `http://model-runner.docker.internal/engines/llama.cpp/v1/`
2. A compatible model loaded (examples use `ai/qwen2.5:0.5B-F16` or `ai/qwen2.5:latest`)
3. Go 1.24.0 or later
4. For MCP examples: **Docker MCP Toolkit** running on port 8811

## Recipes Overview

### 📝 [01 - Simple Chat Agent](./01-simple-chat-agent/)
**Basic chat completion with a single response**

Learn the fundamentals of creating an AI agent that can answer questions with a single, non-streaming response.

```go
bob, _ := robby.NewAgent(
    robby.WithDMRClient(ctx, "http://model-runner.docker.internal/engines/llama.cpp/v1/"),
    robby.WithParams(openai.ChatCompletionNewParams{
        Model: "ai/qwen2.5:0.5B-F16",
        Messages: []openai.ChatCompletionMessageParamUnion{
            openai.SystemMessage("Your name is Bob, You are a pizza expert"),
            openai.UserMessage("What is the best pizza in the world?"),
        },
        Temperature: openai.Opt(0.9),
    }),
)

response, err := bob.ChatCompletion()
```

**Key concepts:** Basic agent creation, chat completion, system and user messages

---

### 🌊 [02 - Simple Chat Stream Agent](./02-simple-chat-stream-agent/)
**Real-time streaming chat responses**

Enhance user experience with streaming responses that display content as it's generated.

```go
bob.ChatCompletionStream(func(self *robby.Agent, content string, err error) error {
    fmt.Print(content)
    return nil
})
```

**Key concepts:** Streaming responses, callback functions, real-time output

---

### 🧠 [03 - Handle Conversational Memory](./03-handle-conversational-memory/)
**Maintaining conversation context across multiple exchanges**

Build agents that remember previous interactions and maintain conversational context.

```go
// Add assistant response to memory
bob.Params.Messages = append(bob.Params.Messages, openai.AssistantMessage(response))

// Add new user question
bob.Params.Messages = append(bob.Params.Messages, openai.UserMessage("Who is his best friend?"))
```

**Key concepts:** Message history, conversational memory, context preservation

---

### 📋 [04 - JSON Output](./04-json-output/)
**Structured data extraction using JSON schemas**

Extract structured information from unstructured text using JSON schema validation.

```go
schema := map[string]any{
    "type": "array",
    "items": map[string]any{
        "type": "object",
        "properties": map[string]any{
            "title": map[string]any{"type": "string"},
            "url": map[string]any{"type": "string"},
            "description": map[string]any{"type": "string"},
        },
        "required": []string{"title", "url", "description"},
    },
}

ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
    OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
        JSONSchema: schemaParam,
    },
}
```

**Key concepts:** JSON schema, structured output, data extraction

---

### 🔧 [05 - Tool Calls](./05-tool-calls/)
**Function calling and tool execution**

Enable your agent to execute custom functions and tools to perform specific tasks.

```go
addTool := openai.ChatCompletionToolParam{
    Function: openai.FunctionDefinitionParam{
        Name:        "add",
        Description: openai.String("add two numbers"),
        Parameters: openai.FunctionParameters{
            "type": "object",
            "properties": map[string]interface{}{
                "a": map[string]string{"type": "number"},
                "b": map[string]string{"type": "number"},
            },
            "required": []string{"a", "b"},
        },
    },
}

results, err := bob.ExecuteToolCalls(map[string]func(any) (any, error){
    "add": func(args any) (any, error) {
        a := args.(map[string]any)["a"].(float64)
        b := args.(map[string]any)["b"].(float64)
        return a + b, nil
    },
})
```

**Key concepts:** Tool definition, function calling, parallel tool execution, custom function implementation

---

### 🐳 [06 - MCP Tool Calls](./06-mcp-tool-calls/)
**Model Context Protocol (MCP) integration**

Integrate with external services and APIs through the Model Context Protocol using Docker MCP Toolkit.

```go
bob, _ := robby.NewAgent(
    robby.WithDMRClient(ctx, "http://model-runner.docker.internal/engines/llama.cpp/v1/"),
    robby.WithParams(params),
    robby.WithMCPClient(robby.WithDockerMCPToolkit()),
    robby.WithMCPTools([]string{"fetch"}),
)

results, err := bob.ExecuteMCPToolCalls()
```

**Key concepts:** MCP integration, external tool access, Docker MCP Toolkit, web fetching

## Running the Examples

Each recipe is a standalone Go module. To run any example:

1. Navigate to the recipe directory:
   ```bash
   cd recipes/01-simple-chat-agent
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the example:
   ```bash
   go run main.go
   ```

## Common Patterns

### Agent Configuration
All recipes follow a consistent pattern for agent creation:

```go
bob, _ := robby.NewAgent(
    robby.WithDMRClient(context.Background(), "http://model-runner.docker.internal/engines/llama.cpp/v1/"),
    robby.WithParams(openai.ChatCompletionNewParams{...}),
    // Additional options...
)
```

### Error Handling
Most examples use simplified error handling for clarity. In production code, always handle errors appropriately:

```go
bob, err := robby.NewAgent(...)
if err != nil {
    log.Fatal(err)
}
```

### Model Selection
Examples use different models based on complexity:
- `ai/qwen2.5:0.5B-F16` - Lightweight model for simple tasks
- `ai/qwen2.5:latest` - Full model for complex tasks

## Contributing

Found an issue or want to add a new recipe? Check out the main project repository and submit a pull request!
