package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

func main() {
	Riker, _ := robby.NewAgent(
		robby.WithOllamaClient(
			context.Background(),
			"http://host.docker.internal:11434/v1",
		),
		robby.WithParams(
			openai.ChatCompletionNewParams{
				Model: "qwen2.5:0.5B",
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage("Your name is Rikker, You are a pizza expert"),
					openai.UserMessage("What is the best pizza in the world?"),
				},
				Temperature: openai.Opt(0.9),
			},
		),
	)

	agentBaseURL := "http://0.0.0.0:8080"
	agentCard, err := Riker.Ping(agentBaseURL)
	if err != nil {
		fmt.Println("Error pinging agent:", err)
		return
	}
	jsonAgentCard, err := robby.AgentCardToJSONString(agentCard)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Agent Card:", jsonAgentCard)

	taskRequest := robby.TaskRequest{
		ID: uuid.NewString(),
		Method: "message/send",
		Params: robby.AgentMessageParams{
			Message: robby.AgentMessage{
				Role: "user",
				Parts: []robby.TextPart{
					{
						Text: "What is the best pizza in the world?",
					},
				},
			},
			// NOTE: I don't know how to query a specific agent skill
			MetaData: map[string]any{
				"skill": "ask_for_something",
			},
		},
	}
	taskResponse, err := Riker.SendToAgent(agentBaseURL, taskRequest)
	if err != nil {
		fmt.Println("Error sending task request:", err)
		return
	}

	jsonTaskResponse, err := robby.TaskResponseToJSONString(taskResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Task Response JSON:", jsonTaskResponse)

	/*
		Riker.ChatCompletionStream(func(self *robby.Agent, content string, err error) error {
			fmt.Print(content)
			return nil
		})
	*/

}
