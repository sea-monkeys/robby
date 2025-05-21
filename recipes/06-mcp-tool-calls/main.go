package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

func main() {


	bob, _ := robby.NewAgent(
		robby.WithDMRClient(
			context.Background(),
			"http://model-runner.docker.internal/engines/llama.cpp/v1/",
		),
		robby.WithParams(
			openai.ChatCompletionNewParams{
				Model: "ai/qwen2.5:0.5B-F16",
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.UserMessage(`
						get the content of https://raw.githubusercontent.com/sea-monkeys/robby/refs/heads/main/README.md
					`),
				},
				Temperature:       openai.Opt(0.0),
				ParallelToolCalls: openai.Bool(true),
			},
		),
		robby.WithMCPClient(robby.WithDockerMCPToolkit()),
		robby.WithMCPTools([]string{"fetch"}), // you must activate the fetch MCP server in Docker MCP Toolkit
	)

	// Generate the tools detection completion
	toolCalls, err := bob.ToolsCompletion()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Number of Tool Calls:\n", len(toolCalls))

	toolCallsJSON, _ := bob.ToolCallsToJSON()
	fmt.Println("Tool Calls:\n", toolCallsJSON)

	results, err := bob.ExecuteMCPToolCalls()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("")

	// Print the results of the tool calls
	fmt.Println("Result of the MCP tool calls execution:")
	for _, result := range results {
		fmt.Println(result)
	}
}
