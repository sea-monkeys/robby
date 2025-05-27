package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

func main() {

	Bob, _ := robby.NewAgent(
		robby.WithDMRClient(
			context.Background(),
			"http://model-runner.docker.internal/engines/llama.cpp/v1/",
		),
		// STEP 1: Define Completion PArameters
		robby.WithParams(
			openai.ChatCompletionNewParams{
				Model: "ai/qwen2.5:0.5B-F16",
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage("Your name is Bob, You are a simple A2A agent server"),
				},
				Temperature: openai.Opt(0.5),
			},
		),
		// STEP 2: Define the settings and skills of the A2A server (for the discovery)
		robby.WithA2AServer(
			robby.A2AServerSettings{
				"name":        "Bob",
				"description": "A simple A2A agent server",
				"url":         "http://localhost:8080",
				"version":     "0.0.0",
			},
			robby.A2AServerSkills{
				map[string]any{
					"id":          "ask_for_something",
					"name":        "Ask for something",
					"description": "Bob is using a small language model to answer questions",
				},
				map[string]any{
					"id":          "say_hello_world",
					"name":        "Say Hello World",
					"description": "Bob can say hello world",
				},
			},
		),
	)


	// STEP 3: Define the Agent Callback function (triggered by A2A Task requests)
	Bob.AgentCallback = func(taskRequest robby.TaskRequest) (robby.TaskResponse, error) {
		// According to A2A spec, user message is in taskRequest.Message.Parts[0].Text
		userMessage := taskRequest.Message.Parts[0].Text
		fmt.Println("✋ Received task request:", taskRequest.ID, "with message:", userMessage, "from role:", taskRequest.Message.Role)
		fmt.Println(robby.TaskRequestToJSONString(taskRequest)) // Print task request in JSON format

		// STEP 4: Process the task request based on the ID
		switch taskRequest.ID {
		case "ask_for_something":
			Bob.Params.Messages = append(
				Bob.Params.Messages,
				openai.UserMessage(userMessage),
			)

		case "say_hello_world":
			Bob.Params.Messages = append(
				Bob.Params.Messages,
				openai.UserMessage("Say hello world to "+userMessage+" from Bob, with emojis."),
			)

		default:
			Bob.Params.Messages = append(
				Bob.Params.Messages,
				openai.UserMessage("Be nice, and explain that "+taskRequest.ID+" is not a valid task ID."),
			)
		}

		// STEP 5: Generate a response using the DMR client + Chat Completion 
		responseText, err := Bob.ChatCompletion()
		if err != nil {
			return robby.TaskResponse{}, fmt.Errorf("error generating response: %w", err)
		}

		// Formulate response in A2A Task format
		// We'll return a Task object with final state = 'completed' and agent message
		// STEP 6: Create and return the response task
		responseTask := robby.TaskResponse{
			ID:     taskRequest.ID, // use the same task ID
			Status: robby.TaskStatus{State: "completed"},
			Messages: []robby.AgentMessage{
				taskRequest.Message, // include original user message in history
				{
					Role: "agent", // agent's response
					Parts: []robby.TextPart{
						{Text: responseText}, // agent message content as TextPart
					},
				},
			},
		}

		return responseTask, nil

	}
	// TODO: make an example with streaming responses

	fmt.Println("Agent Bob initialized with A2A server settings.")
	fmt.Println("Agent Name:", Bob.AgentCard.Name)

	err := Bob.StartA2AServer("0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error starting A2A Bob server:", err)
		return
	}

}
