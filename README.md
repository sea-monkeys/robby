# Robby 🤖
> Create an AI Agent with Docker Model Runner and Docker MCP Toolkit

Je vais vous fournir une explication en anglais pour le README concernant ce projet "Robby".

## Project Purpose Explanation

**Robby** is not intended to become a full-fledged framework. It is designed as a lightweight library that implements design patterns to help developers build **AI Agents** using the **OpenAI Go SDK** in conjunction with **[Docker Model Runner](https://docs.docker.com/model-runner/)** and **[Docker MCP Toolkit](https://docs.docker.com/ai/mcp-catalog-and-toolkit/toolkit/)**. 

The goal is not to hide or abstract away the OpenAI SDK usage, but rather to provide helpful utilities that simplify the developer experience. **Robby** offers convenience methods and patterns that make it easier to work with these technologies while still allowing direct access to the underlying SDK functionality when needed.

This approach gives developers the flexibility to leverage the full power of the OpenAI Go SDK while benefiting from the streamlined workflows that Robby provides for common AI agent implementation patterns.

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