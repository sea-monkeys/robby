# Robby 🤖
> Create an AI Agent with Docker Model Runner and Docker MCP Toolkit


## Chat Agent

```go
bob, err := NewAgent(
    WithDMRClient(
        context.Background(),
        "http://model-runner.docker.internal/engines/llama.cpp/v1/",
    ),
    WithParams(
        openai.ChatCompletionNewParams{
            Model: "ai/qwen2.5:latest",
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.SystemMessage("You are a pizza expert"),
                openai.UserMessage("What is the best pizza in the world?"),
            },
            Temperature: openai.Opt(0.9),
        },
    ),
)
// ...
bob.ChatCompletionStream(func(self *Agent, content string, err error) error{
    fmt.Print(content)
    return nil	
})
```

## Tools Agent

```go
bob, err := NewAgent(
    WithDMRClient(
        context.Background(),
        "http://model-runner.docker.internal/engines/llama.cpp/v1/",
    ),
    WithParams(
        openai.ChatCompletionNewParams{
            Model: "ai/qwen2.5:latest",
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.UserMessage(`
                Search information about Hawaiian pizza.(only 3 results)
                Search information about Mexican pizza.(only 3 results)
                `),
            },
            Temperature:       openai.Opt(0.0),
            ParallelToolCalls: openai.Bool(true),
        },
    ),
    WithMCPToolkitClient(WithDocker()),
    WithTools([]string{"brave_web_search"}),
)

results, err := bob.ToolsCompletion() 

fmt.Println("Tool Calls Done:")
for idx, result := range results {
    fmt.Println(fmt.Sprintf("%d.", idx), result)
}
```
