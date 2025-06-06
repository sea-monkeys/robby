package robby

import (
	"context"
	"fmt"
	"testing"

	"github.com/openai/openai-go"
)

func TestChat(t *testing.T) {
	
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
					openai.UserMessage("[Brief] What is the best pizza in the world?"),
				},
				Temperature: openai.Opt(0.9),
			},
		),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	response, err := bob.ChatCompletionStream(func(self *Agent, content string, err error) error{
		fmt.Print(content)
		return nil	
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Add the assistant message to the messages to keep the conversation going
	bob.Params.Messages = append(bob.Params.Messages, openai.AssistantMessage(response))

	

}
