package robby

import (
	"context"
	"fmt"
	"testing"

	"github.com/openai/openai-go"
)

func TestToolCallsToJSON(t *testing.T) {
	// Create an agent with tools
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

	// Create the agent
	agent, err := NewAgent(
		WithDMRClient(
			context.Background(),
			"http://model-runner.docker.internal/engines/llama.cpp/v1/",
		),
		WithParams(
			openai.ChatCompletionNewParams{
				Model: "ai/qwen2.5:0.5B-F16",
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.UserMessage(`
						Say hello to Bob.
						Give a vulcan salute to James Kirk.
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
		t.Fatalf("Failed to create agent: %v", err)
	}

	// Test 1: No tool calls yet
	jsonStr, err := agent.ToolCallsToJSON()
	if err != nil {
		t.Errorf("ToolCallsToJSON failed on empty tool calls: %v", err)
	}
	if jsonStr != "[]" {
		t.Errorf("Expected empty array JSON '[]', got: %s", jsonStr)
	}

	toolsCall, err := agent.ToolsCompletion()
	if err != nil {
		t.Fatalf("Failed to get tool calls: %v", err)
	}
	
	jsonStr, err = agent.ToolCallsToJSON()
	if err != nil {
		t.Errorf("ToolCallsToJSON failed after tool calls: %v", err)
	}
	if jsonStr == "[]" {
		t.Errorf("Expected non-empty JSON, got: %s", jsonStr)
	}

	fmt.Println("📝 Tool Calls JSON:\n", jsonStr)

	newJsonStr, err := ToolCallsToJSONString(toolsCall)
	if err != nil {
		t.Errorf("ToolCallsToJSONString failed: %v", err)
	}

	fmt.Println("📝 Tool Calls JSON from ToolCallsToJSONString:\n", newJsonStr)

	if newJsonStr != jsonStr {
		t.Errorf("ToolCallsToJSONString and ToolCallsToJSON mismatch: %s != %s", newJsonStr, jsonStr)
	}

}
