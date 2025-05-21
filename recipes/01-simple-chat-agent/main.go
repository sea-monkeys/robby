package main

import (
	"context"

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
					openai.SystemMessage("Your name is Bob, You are a pizza expert"),
					openai.UserMessage("What is the best pizza in the world?"),
				},
				Temperature: openai.Opt(0.9),
			},
		),
	)

	response, err := bob.ChatCompletion()
	if err != nil {
		panic(err)
	}
	// Print the response
	println(response)
}
