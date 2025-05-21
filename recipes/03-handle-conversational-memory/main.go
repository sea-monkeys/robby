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
					openai.SystemMessage("You are a Star Trek expert"),
					openai.UserMessage("Who is James Kirk?"),
				},
				Temperature: openai.Opt(0.9),
			},
		),
	)

	response, err := bob.ChatCompletionStream(func(self *robby.Agent, content string, err error) error {
		fmt.Print(content)
		return nil
	})

	if condition := err != nil; condition {
		panic(err)
	}

	// Add assistant response to memory
	bob.Params.Messages = append(bob.Params.Messages, openai.AssistantMessage(response))

	// Add new user question
	bob.Params.Messages = append(bob.Params.Messages, openai.UserMessage("Who is his best friend?"))

	fmt.Println("")
	fmt.Println("")

	_, err = bob.ChatCompletionStream(func(self *robby.Agent, content string, err error) error {
		fmt.Print(content)
		return nil
	})

	if condition := err != nil; condition {
		panic(err)
	}

}
