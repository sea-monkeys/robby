# Robby 🤖
> Create an AI Agent with Docker Model Runner and Docker MCP Toolkit

## Install

```bash
go get github.com/sea-monkeys/robbyby@v0.0.0
```

## Use

```golang
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
				Model: "ai/qwen2.5:latest",
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage("Your name is Bob, You are a pizza expert"),
					openai.UserMessage("What is the best pizza in the world?"),
				},
				Temperature: openai.Opt(0.9),
			},
		),
	)

	response, _ := bob.ChatCompletionStream(func(self *robby.Agent, content string, err error) error {
		fmt.Print(content)
		return nil
	})

	// Add the assistant message to the messages to keep the conversation going
	bob.Params.Messages = append(bob.Params.Messages, openai.AssistantMessage(response))

}
```