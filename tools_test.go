package robby

import (
	"context"
	"fmt"
	"testing"

	"github.com/openai/openai-go"
)

func TestTools(t *testing.T) {

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
		/*
			WithMCPToolkitClient([]string{
				"docker",
				"run",
				"-i",
				"--rm",
				"alpine/socat",
				"STDIO",
				"TCP:host.docker.internal:8811",
			}),
		*/
		WithMCPTools([]string{"brave_web_search"}),
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

	results, err := bob.ExecuteMCPToolCalls()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Tools Done:")
	for idx, result := range results {
		fmt.Println("---------------------------------")
		fmt.Println(fmt.Sprintf("%d.", idx), result)
	}

}
