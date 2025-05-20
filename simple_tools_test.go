package robby

import (
	"context"
	"fmt"
	"testing"

	"github.com/openai/openai-go"
)

func TestSimpleTools(t *testing.T) {

	sayHelloTool := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "say_hello",
			Description: openai.String("Say hello to the given person name"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]string{
						"type": "string",
					},
				},
				"required": []string{"name"},
			},
		},
	}

	vulcanSaluteTool := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "vulcan_salute",
			Description: openai.String("Give a vulcan salute to the given person name"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]string{
						"type": "string",
					},
				},
				"required": []string{"name"},
			},
		},
	}	


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
						Say hello to Bob.
						Give a vulcan salute to James Kirk.
						Say hello to Spock.
					`),
				},
				Temperature:       openai.Opt(0.0),
				ParallelToolCalls: openai.Bool(true),
			},
		),
		WithTools([]openai.ChatCompletionToolParam{
			sayHelloTool,
			vulcanSaluteTool,
		}),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	toolCalls, err := bob.ToolsCompletion() // This add the Tools to the agent.Params
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	toolCallsJSON, _ := ToolCallsToJSONString(toolCalls)
	fmt.Println("Tool Calls:\n", toolCallsJSON)

	if len(toolCalls) == 0 {
		fmt.Println("No tools found.")
		return
	}
	results, err := bob.ExecuteToolCalls(map[string]func(any) any{
		"say_hello": func(args any) any {
			name := args.(map[string]any)["name"].(string)
			return fmt.Sprintf("Hello, %s!", name)
		},
		"vulcan_salute": func(args any) any {
			name := args.(map[string]any)["name"].(string)
			return fmt.Sprintf("Live long and prosper, %s!", name)
		},
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Tools Done:")
	for _, result := range results {
		fmt.Println("---------------------------------")
		fmt.Println(result)
	}


}
