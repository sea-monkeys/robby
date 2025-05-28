package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

//var callbackMutex sync.Mutex // Mutex to serialize access to AgentCallback
var callbackMutex sync.RWMutex // Use RWMutex for concurrent reads

func main() {

	Bob, _ := robby.NewAgent(
		robby.WithDMRClient(
			context.Background(),
			"http://model-runner.docker.internal/engines/llama.cpp/v1/",
		),
		// STEP 1: Define Completion PArameters
		robby.WithParams(
			openai.ChatCompletionNewParams{
				//Model:       "ai/qwen2.5:0.5B-F16",
				Model:       "ai/qwen2.5:3B-F16",
				Messages:    []openai.ChatCompletionMessageParamUnion{},
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

	/* NOTE:
		Explanation
		Mutex Initialization: A sync.Mutex named callbackMutex is declared globally.
		Locking: callbackMutex.Lock() ensures that only one request is processed at a time.
		Unlocking: defer callbackMutex.Unlock() releases the lock after the request is processed.
		This approach guarantees serialized handling of requests, ensuring that the second request is processed only after the first one is completed.

	*/

	// Serialize access to AgentCallback using a mutex
	Bob.AgentCallback = func(taskRequest robby.TaskRequest) (robby.TaskResponse, error) {
		callbackMutex.Lock()         // Lock the mutex before processing the request
		defer callbackMutex.Unlock() // Unlock the mutex after processing the request

		// According to A2A spec, user message is in taskRequest.Message.Parts[0].Text
		userMessage := taskRequest.Params.Message.Parts[0].Text

		fmt.Println("🟢 Processing task request:", taskRequest.ID)
		fmt.Println("🔵 UserMessage:", userMessage)
		fmt.Println("🟡 TaskRequest Metadata:", taskRequest.Params.MetaData)

		// STEP 4: Process the task request based on the ID
		switch taskRequest.Params.MetaData["skill"] {
		case "ask_for_something":

			Bob.Params.Messages = []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are Bob, a simple A2A agent. You can answer questions."),
				openai.UserMessage(userMessage),
			}

		case "say_hello_world":

			Bob.Params.Messages = []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are Bob, a simple A2A agent. You can answer questions."),
				openai.UserMessage("Say hello world to " + userMessage + " from Bob, with emojis."),
			}

		default:

			Bob.Params.Messages = []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are Bob, a simple A2A agent. You can answer questions."),
				openai.UserMessage("Be nice, and explain that " + taskRequest.ID + " is not a valid task ID."),
			}

		}

		// STEP 5: Generate a response using the DMR client + Chat Completion
		responseText, err := Bob.ChatCompletion()
		if err != nil {
			return robby.TaskResponse{}, fmt.Errorf("error generating response: %w", err)
		}

		fmt.Println("🤖 Generated response:", responseText)

		// Formulate response in A2A Task format
		// We'll return a Task object with final state = 'completed' and agent message
		// STEP 6: Create and return the response task
		responseTask := robby.TaskResponse{
			ID:             taskRequest.ID, // use the same task ID
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
		fmt.Println("🟩 Response Task Id:", responseTask.ID)

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
