# Robby 🤖

> **Lightweight Go library[pattern] for building AI Agents with Docker Model Runner and Docker MCP Toolkit**

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Required-2496ED?style=flat&logo=docker)](https://docs.docker.com/model-runner/)

**Robby** is a lightweight Go library[pattern] that provides design patterns and utilities for building **AI Agents** using the **OpenAI Go SDK** with **[Docker Model Runner](https://docs.docker.com/model-runner/)** and **[Docker MCP Toolkit](https://docs.docker.com/ai/mcp-catalog-and-toolkit/toolkit/)**. 

Rather than hiding the OpenAI SDK, Robby enhances your development experience with convenient patterns while preserving full access to the underlying SDK functionality.

## ✨ Features

- **🚀 Simple Agent Creation** - Minimal boilerplate with powerful configuration options
- **💬 Streaming & Non-Streaming Chat** - Real-time responses with callback support
- **🧠 Conversational Memory** - Built-in patterns for maintaining conversation context
- **🔧 Tool Integration** - Easy function calling with custom tool implementations
- **🐳 MCP Protocol Support** - Seamless integration with Model Context Protocol services
- **📋 Structured Output** - JSON schema validation for reliable data extraction
- **⚡ Docker Model Runner** - Local LLM execution without external API dependencies
- **🛠️ Developer Friendly** - Comprehensive examples and clear documentation

## 🚀 Quick Start

### Prerequisites

1. **Docker Model Runner** - [Installation Guide](https://docs.docker.com/model-runner/)
2. **Go 1.24+** - [Download Go](https://golang.org/dl/)
3. **Compatible Model** - Any model supported by Docker Model Runner

### Installation

```bash
go get github.com/sea-monkeys/robby@v0.0.1
```

### Basic Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/openai/openai-go"
    "github.com/sea-monkeys/robby"
)

func main() {
    // Create an AI agent
    bob, err := robby.NewAgent(
        robby.WithDMRClient(
            context.Background(),
            "http://model-runner.docker.internal/engines/llama.cpp/v1/",
        ),
        robby.WithParams(
            openai.ChatCompletionNewParams{
                Model: "ai/qwen2.5:latest",
                Messages: []openai.ChatCompletionMessageParamUnion{
                    openai.SystemMessage("You are a helpful assistant"),
                    openai.UserMessage("What is the capital of France?"),
                },
                Temperature: openai.Opt(0.7),
            },
        ),
    )
    if err != nil {
        panic(err)
    }

    // Get a response
    response, err := bob.ChatCompletion()
    if err != nil {
        panic(err)
    }

    fmt.Println("Assistant:", response)
}
```

## 📚 Examples & Recipes

Explore comprehensive examples in the [`recipes/`](./recipes/) directory:

| Recipe | Description | Key Features |
|--------|-------------|--------------|
| **[01 - Simple Chat](./recipes/01-simple-chat-agent/)** | Basic chat completion | Agent creation, single response |
| **[02 - Streaming Chat](./recipes/02-simple-chat-stream-agent/)** | Real-time streaming responses | Live content delivery, callbacks |
| **[03 - Memory Management](./recipes/03-handle-conversational-memory/)** | Conversation context handling | Multi-turn conversations, history |
| **[04 - JSON Output](./recipes/04-json-output/)** | Structured data extraction | Schema validation, data parsing |
| **[05 - Tool Calls](./recipes/05-tool-calls/)** | Custom function execution | Function calling, parallel execution |
| **[06 - MCP Integration](./recipes/06-mcp-tool-calls/)** | MCP External service integration | MCP tool calls |
| **[07 - MCP Integration](./recipes/07-use-my-mcp-server/)** | MCP External service integration | MCP tool calls with MCP Server sample |
| **[08 - RAG in memory](./recipes/08-rag-memory/)** | Retrieval Augmented Generation | Embeddings and Similarity search |

## 🔧 Core Concepts

### Agent Configuration

Robby uses a functional options pattern for clean, flexible configuration:

```go
agent, err := robby.NewAgent(
    // Docker Model Runner client
    robby.WithDMRClient(ctx, "http://model-runner.docker.internal/engines/llama.cpp/v1/"),
    
    // Chat parameters
    robby.WithParams(openai.ChatCompletionNewParams{
        Model: "ai/qwen2.5:latest",
        Messages: messages,
        Temperature: openai.Opt(0.7),
    }),
    
    // Custom tools
    robby.WithTools(customTools),
    
    // MCP integration
    robby.WithMCPClient(robby.WithDockerMCPToolkit()),
    robby.WithMCPTools([]string{"fetch", "brave_web_search"}),
)
```

### Streaming Responses

Enhance user experience with real-time content delivery:

```go
response, err := agent.ChatCompletionStream(func(self *robby.Agent, content string, err error) error {
    if err != nil {
        return err
    }
    fmt.Print(content) // Display content as it arrives
    return nil
})
```

### Function Calling

Enable your agent to execute custom functions:

```go
// Define tools
calculatorTool := openai.ChatCompletionToolParam{
    Function: openai.FunctionDefinitionParam{
        Name: "calculate",
        Description: openai.String("Perform mathematical calculations"),
        Parameters: openai.FunctionParameters{
            "type": "object",
            "properties": map[string]interface{}{
                "expression": map[string]string{
                    "type": "string",
                    "description": "Math expression to evaluate",
                },
            },
            "required": []string{"expression"},
        },
    },
}

// Create agent with tools
agent, _ := robby.NewAgent(
    robby.WithDMRClient(ctx, baseURL),
    robby.WithParams(params),
    robby.WithTools([]openai.ChatCompletionToolParam{calculatorTool}),
)

// Detect tool calls
toolCalls, err := agent.ToolsCompletion()

// Execute tools
results, err := agent.ExecuteToolCalls(map[string]func(any) (any, error){
    "calculate": func(args any) (any, error) {
        expression := args.(map[string]any)["expression"].(string)
        // Implement calculation logic
        return evaluateExpression(expression), nil
    },
})
```

### MCP Integration

Connect to external services through Model Context Protocol:

```go
agent, _ := robby.NewAgent(
    robby.WithDMRClient(ctx, baseURL),
    robby.WithParams(params),
    robby.WithMCPClient(robby.WithDockerMCPToolkit()),
    robby.WithMCPTools([]string{"fetch", "brave_web_search"}),
)

// Agent can now use web search and HTTP fetching
toolCalls, _ := agent.ToolsCompletion()
results, _ := agent.ExecuteMCPToolCalls()
```

## 🛠️ Advanced Features

### Conversation Memory

Maintain context across multiple exchanges:

```go
// Add assistant response to conversation history
agent.Params.Messages = append(agent.Params.Messages, openai.AssistantMessage(response))

// Continue the conversation
agent.Params.Messages = append(agent.Params.Messages, openai.UserMessage("Tell me more about that"))
```

### Structured Output

Extract structured data using JSON schemas:

```go
schema := map[string]any{
    "type": "object",
    "properties": map[string]any{
        "name": map[string]any{"type": "string"},
        "age": map[string]any{"type": "number"},
        "skills": map[string]any{
            "type": "array",
            "items": map[string]any{"type": "string"},
        },
    },
}

params := openai.ChatCompletionNewParams{
    Model: "ai/qwen2.5:latest",
    Messages: messages,
    ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
        OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
            JSONSchema: schemaParam,
        },
    },
}
```

### Tool Call Debugging

Monitor and debug tool interactions:

```go
toolCalls, _ := agent.ToolsCompletion()

// Convert to JSON for inspection
toolCallsJSON, _ := agent.ToolCallsToJSON()
fmt.Println("Tool Calls:", toolCallsJSON)

// Or use the standalone function
jsonStr, _ := robby.ToolCallsToJSONString(toolCalls)
```

## 🐳 Docker Integration

### Docker Model Runner Connection

The connection URL depends on where your application is running:

**🐳 Application running in a container (DevContainer, Docker, etc.):**
```go
robby.WithDMRClient(
    context.Background(),
    "http://model-runner.docker.internal/engines/llama.cpp/v1/",
)
```

**💻 Application running directly on your machine:**
```go
robby.WithDMRClient(
    context.Background(),
    "http://localhost:12434/engines/v1",
)
```

### Docker MCP Toolkit Connection

Robby provides two methods to connect to Docker MCP Toolkit:

#### Option 1: `WithDockerMCPToolkit()` (Recommended)
Uses Docker container with Alpine/socat to establish the connection:

```go
agent, _ := robby.NewAgent(
    robby.WithDMRClient(ctx, baseURL),
    robby.WithParams(params),
    robby.WithMCPClient(robby.WithDockerMCPToolkit()), // Uses Docker container
    robby.WithMCPTools([]string{"fetch", "brave_web_search"}),
)
```

**Advantages:**
- No local dependencies required
- Works in any Docker environment
- Self-contained solution

**Requirements:**
- Docker must be available
- Docker MCP Toolkit running on `host.docker.internal:8811`

#### Option 2: `WithSocatMCPToolkit()`
Uses local socat command for direct connection:

```go
agent, _ := robby.NewAgent(
    robby.WithDMRClient(ctx, baseURL),
    robby.WithParams(params),
    robby.WithMCPClient(robby.WithSocatMCPToolkit()), // Uses local socat
    robby.WithMCPTools([]string{"fetch", "brave_web_search"}),
)
```

**Advantages:**
- Faster connection (no Docker overhead)
- Direct system integration
- **Ideal for dockerizing your agent application** (avoids Docker-in-Docker complexity)

**Requirements:**
- `socat` must be installed on your system
- Docker MCP Toolkit running on `host.docker.internal:8811`

**Installing socat:**
```bash
# macOS
brew install socat

# Ubuntu/Debian
sudo apt-get install socat

# Alpine Linux
apk add socat
```

## 📖 API Documentation

### Core Types

- **`Agent`** - Main agent structure with context, client, and configuration
- **`AgentOption`** - Functional option for agent configuration
- **`STDIOCommandOption`** - Command configuration for MCP integration

### Key Methods

- **`NewAgent(options ...AgentOption)`** - Create configured agent instance
- **`ChatCompletion()`** - Single response completion
- **`ChatCompletionStream(callback)`** - Streaming response with callbacks
- **`ToolsCompletion()`** - Detect tool calls from model response
- **`ExecuteToolCalls(implementations)`** - Execute custom tool functions
- **`ExecuteMCPToolCalls()`** - Execute MCP protocol tools

For complete API documentation, see our [detailed function reference](./docs/api.md).

## 🧪 Testing

Run the test suite:

```bash
# Basic functionality tests
./simple.tools.tests.sh

# Chat functionality tests  
./chat.tests.sh

# MCP integration tests
./agent.mcp.tests.sh

# Tool call JSON formatting tests
./tool.call.json.test.sh
```

## 🎯 Use Cases

**Robby is perfect for:**

- **Chatbots & Virtual Assistants** - Build conversational AI with memory and tools
- **Content Generation** - Create articles, summaries, and structured content
- **Data Processing** - Extract and transform unstructured data into JSON
- **API Integration** - Connect AI agents to external services and databases
- **Workflow Automation** - Automate complex multi-step processes
- **Research Assistants** - Web search, data analysis, and report generation

## 🚧 Development Status

Robby is currently in **v0.0.1** - early development phase. The API may change as we gather feedback and add features. 

**Planned Features:**
- Enhanced error handling and recovery
- Built-in RAG (Retrieval Augmented Generation) support
- Advanced conversation management
- More MCP integrations

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **🐛 Bug Reports** - Create detailed issues with reproduction steps
2. **✨ Feature Requests** - Suggest new capabilities and improvements
3. **📖 Documentation** - Improve examples, guides, and API docs
4. **🧪 Testing** - Add test cases and improve coverage
5. **💡 Examples** - Contribute new recipes and use cases

### Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/sea-monkeys/robby.git
   cd robby
   ```

2. Start development environment:
   ```bash
   # Using DevContainer (recommended)
   code . # Open in VS Code with DevContainer

   # Or manually
   go mod download
   ```

3. Run tests:
   ```bash
   go test ./...
   ```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **OpenAI** - For the excellent Go SDK that powers our chat completions
- **Docker** - For Docker Model Runner and MCP Toolkit that enable local AI execution
- **Community** - For feedback, contributions, and real-world testing

## 📬 Support & Community

- **Documentation** - Explore the [`recipes/`](./recipes/) directory for examples
- **Issues** - Report bugs and request features on [GitHub Issues](https://github.com/sea-monkeys/robby/issues)
- **Discussions** - Join community discussions and get help
- **Examples** - Check out real-world implementations in our recipe collection

---

**Built with ❤️ by the Sea Monkeys team**

*Robby: Making AI agent development simple, powerful, and fun!*