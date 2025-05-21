package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/sea-monkeys/robby"
)

func main() {

	content := `
	Title: WebAssembly
	Description: WebAssembly (abbreviated <strong>Wasm</strong>) is a binary instruction format for a stack-based virtual machine. <strong>Wasm</strong> is designed as a portable compilation target for programming languages, enabling deployment on the web for client and server applications.
	URL: https://webassembly.org/

	Title: WebAssembly - Wikipedia
	Description: WebAssembly (<strong>Wasm</strong>) defines a portable binary-code format and a corresponding text format for executable programs as well as software interfaces for facilitating communication between such programs and their host environment.
	URL: https://en.wikipedia.org/wiki/WebAssembly

	Title: Wasm 2.0 Completed - WebAssembly
	Description: For the most up-to-date version of the current specification, we recommend looking at the documents hosted on our GitHub page. This always includes the <strong>latest</strong> fixes and offers multiple different formats for reading and browsing. For those who are not following the evolution of <strong>Wasm</strong> as closely, ...
	URL: https://webassembly.org/news/2025-03-20-wasm-2.0/ Title: Introduction · WASI.dev
	Description: The WebAssembly System Interface (<strong>WASI</strong>) is a group of standards-track API specifications for software compiled to the W3C WebAssembly (Wasm) standard. <strong>WASI</strong> is designed to provide a secure standard interface for applications that can be compiled to Wasm from any language, and that may run ...
	URL: https://wasi.dev/
	`

	schema := map[string]any{
		"type": "array",
		"items": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"title": map[string]any{
					"type": "string",
				},
				"url": map[string]any{
					"type": "string",
				},
				"description": map[string]any{
					"type": "string",
				},
			},
			"required": []string{"title", "url", "description"},
		},
	}

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "search_results",
		Description: openai.String("Notable information about search results"),
		Schema:      schema,
		Strict:      openai.Bool(true),
	}

	// This agent will use the JSON schema to parse the results
	// and return the results in a JSON format.
	// The JSON schema is defined in the schemaParam variable
	// and is passed to the agent as a parameter.
	bob, _ := robby.NewAgent(
		robby.WithDMRClient(
			context.Background(),
			"http://model-runner.docker.internal/engines/llama.cpp/v1/",
		),
		robby.WithParams(
			openai.ChatCompletionNewParams{
				Model: "ai/qwen2.5:latest",
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(content),
					openai.UserMessage("give me the list of the results."),
				},
				Temperature: openai.Opt(0.0),
				ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
					OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
						JSONSchema: schemaParam,
					},
				},
			},
		),
	)

	jsonResults, err := bob.ChatCompletion()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("📝 JSON Results:\n", jsonResults)
}
