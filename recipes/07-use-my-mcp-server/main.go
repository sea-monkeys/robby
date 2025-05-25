package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

// Create a custom MCP server command optio
func WithMyMCPServer() robby.STDIOCommandOption {
	return robby.STDIOCommandOption{
		"docker",
		"run",
		"-i",
		"--rm",
		"k33g/mcp-demo:with-agents",
		"--debug",
		"--plugins",
		"./plugins",
	}
}

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
						Add 5 and 3
						Subtract 10 from 4
						Multiply 2 and 6
						Divide 8 by 2
					`),
				},
				Temperature:       openai.Opt(0.0),
				ParallelToolCalls: openai.Bool(true),
			},
		),
		robby.WithMCPClient(WithMyMCPServer()),
		robby.WithMCPTools([]string{"add", "subtract", "multiply", "divide"}),
		//robby.WithMCPTools([]string{}),
		robby.WithMCPResources([]string{}),
		robby.WithMCPPrompts([]string{}),
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

	fmt.Println("--------------------------------------")
	fmt.Println("Resources List:")
	for _, resource := range bob.Resources {
		fmt.Printf("URI: %s, Name: %s, Description: %s, MimeType: %s\n",
			resource.URI, resource.Name, resource.Description, resource.MimeType)
	}
	fmt.Println("--------------------------------------")
	fmt.Println("Resource Read:")

	addRsrc, err := bob.ReadResource("info:///calculator")
	if err != nil {
		fmt.Println("Error reading resource:", err)
		return
	}
	fmt.Println(addRsrc.Description, ":")
	fmt.Println(addRsrc.Text) // Read the resource for calculator information
	fmt.Println("--------------------------------------")
	fmt.Println("Prompts List:")
	for _, prompt := range bob.Prompts {
		fmt.Println("Name: ", prompt.Name)
		fmt.Println("Description:", prompt.Description)
		fmt.Println("Args:", prompt.Arguments)

	}
	fmt.Println("--------------------------------------")
	fmt.Println("Prompt Get:")

	args := map[string]any{
		"operation": "add",
		"a":         5,
		"b":         3,
	}

	prompt, err := bob.GetPrompt("calculator_prompt", args)
	if err != nil {
		fmt.Println("Error getting prompt:", err)
		return
	}
	fmt.Println("Prompt Name:", prompt.Name)
	fmt.Println("Prompt Description:", prompt.Description)
	fmt.Println(prompt.Messages[0].Role, ":", prompt.Messages[0].Content.Text)
	fmt.Println("======================================")

	riker, _ := robby.NewAgent(
		robby.WithDMRClient(
			context.Background(),
			"http://model-runner.docker.internal/engines/llama.cpp/v1/",
		),
		robby.WithParams(
			openai.ChatCompletionNewParams{
				Model:             "ai/qwen2.5:0.5B-F16",
				Messages:          []openai.ChatCompletionMessageParamUnion{},
				Temperature:       openai.Opt(0.0),
				ParallelToolCalls: openai.Bool(true),
			},
		),
		robby.WithMCPClient(WithMyMCPServer()),
		robby.WithMCPTools([]string{"add", "subtract", "multiply", "divide"}),
		robby.WithMCPResources([]string{}),
		robby.WithMCPPrompts([]string{}),
	)

	systemInstructions, err := riker.ReadResource("info:///calculator")
	if err != nil {
		fmt.Println("Error reading resource:", err)
		return
	}
	userInstructions, err := riker.GetPrompt("calculator_prompt", map[string]any{
		"operation": "add",
		"a":         25,
		"b":         50,
	})
	if err != nil {
		fmt.Println("Error getting prompt:", err)
		return
	}

	riker.Params.Messages = []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemInstructions.Text),
		openai.UserMessage(userInstructions.Messages[0].Content.Text),
	}

	toolCalls, err = riker.ToolsCompletion()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Number of Tool Calls:\n", len(toolCalls))
	toolCallsJSON, _ = riker.ToolCallsToJSON()
	fmt.Println("Tool Calls:\n", toolCallsJSON)
	results, err = riker.ExecuteMCPToolCalls()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("")
	fmt.Println("Result of the MCP tool calls execution:")
	for _, result := range results {
		fmt.Println(result)
	}
	fmt.Println("--------------------------------------")

}
