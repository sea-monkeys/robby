package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

// Use RWMutex for better concurrency - allows multiple reads but exclusive writes
var callbackMutex sync.RWMutex

func main() {
	Bob, _ := robby.NewAgent(
		robby.WithDMRClient(
			context.Background(),
			"http://model-runner.docker.internal/engines/llama.cpp/v1/",
		),
		robby.WithParams(
			openai.ChatCompletionNewParams{
				Model:       "ai/qwen2.5:1.5B-F16",
				Messages:    []openai.ChatCompletionMessageParamUnion{},
				Temperature: openai.Opt(0.5),
			},
		),
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

	// Improved AgentCallback with better error handling and logging
	Bob.AgentCallback = func(taskRequest robby.TaskRequest) (robby.TaskResponse, error) {
		fmt.Printf("🟢 Processing task request: %s\n", taskRequest.ID)

		// Lock only for the duration of processing this specific request
		//callbackMutex.Lock()
		//defer callbackMutex.Unlock()

		// Extract user message
		userMessage := taskRequest.Params.Message.Parts[0].Text
		fmt.Printf("🔵 UserMessage: %s\n", userMessage)
		fmt.Printf("🟡 TaskRequest Metadata: %v\n", taskRequest.Params.MetaData)

		// Create a fresh context for each request to avoid conflicts
		ctx := context.Background()

		// Create a new agent instance for this request to avoid state conflicts
		// Or alternatively, just update the messages without creating new agent
		var systemMessage, userPrompt string

		fmt.Println("🟪", userMessage)

		switch taskRequest.Params.MetaData["skill"] {
		case "ask_for_something":
			systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
			userPrompt = userMessage

		case "greetings":
		
			systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
			userPrompt = "Greetings to " + userMessage + " with emojis and use his name."

		default:
			systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
			userPrompt = "Be nice, and explain that " + fmt.Sprintf("%v", taskRequest.Params.MetaData["skill"]) + " is not a valid task ID."
		}

		// Create a new agent instance for this specific request to avoid state conflicts
		requestAgent, err := robby.NewAgent(
			robby.WithDMRClient(ctx, "http://model-runner.docker.internal/engines/llama.cpp/v1/"),
			robby.WithParams(
				openai.ChatCompletionNewParams{
					Model: "ai/qwen2.5:3B-F16",
					Messages: []openai.ChatCompletionMessageParamUnion{
						openai.SystemMessage(systemMessage),
						openai.UserMessage(userPrompt),
					},
					Temperature: openai.Opt(0.5),
				},
			),
		)

		if err != nil {
			return robby.TaskResponse{}, fmt.Errorf("error creating request agent: %w", err)
		}

		// Generate response
		responseText, err := requestAgent.ChatCompletion()
		if err != nil {
			return robby.TaskResponse{}, fmt.Errorf("error generating response: %w", err)
		}

		fmt.Printf("🤖 Generated response: %s\n", responseText)

		// Create response task
		responseTask := robby.TaskResponse{
			ID:             taskRequest.ID,
			JSONRpcVersion: "2.0",
			Result: robby.Result{
				Status: robby.TaskStatus{
					State: "completed",
				},
				History: []robby.AgentMessage{
					{
						Role: "assistant",
						Parts: []robby.TextPart{
							{
								Text: responseText,
								Type: "text",
							},
						},
					},
				},
				Kind:     "task",
				Metadata: map[string]any{},
			},
		}

		fmt.Printf("🟩 Response Task Id: %s\n", responseTask.ID)
		return responseTask, nil
	}

	fmt.Println("Agent Bob initialized with A2A server settings.")
	fmt.Printf("Agent Name: %s\n", Bob.AgentCard.Name)
	fmt.Println("Starting server on 0.0.0.0:8080...")

	err := Bob.StartA2AServer("0.0.0.0:8080")
	if err != nil {
		fmt.Printf("Error starting A2A Bob server: %v\n", err)
		return
	}
}
