package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

func main() {

	addTool := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "add",
			Description: openai.String("add two numbers"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]string{
						"type":        "number",
						"description": "The first number to add.",
					},
					"b": map[string]string{
						"type":        "number",
						"description": "The second number to add.",
					},
				},
				"required": []string{"a", "b"},
			},
		},
	}

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
						Add 34 and 56.
						Add 12 and 34.
						Add 1 and 2.
						Add 3 and 4.
					`),
				},
				Temperature:       openai.Opt(0.0),
				ParallelToolCalls: openai.Bool(true),
			},
		),
		robby.WithTools([]openai.ChatCompletionToolParam{
			addTool,
		}),
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

	results, err := bob.ExecuteToolCalls(map[string]func(any) (any, error){
		"add": func(args any) (any, error) {
			a := args.(map[string]any)["a"].(float64)
			b := args.(map[string]any)["b"].(float64)
			return a + b, nil
		},
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("")

	// Print the results of the tool calls
	fmt.Println("Results of the tool calls execution:")
	for _, result := range results {
		fmt.Println(result)
	}
}
