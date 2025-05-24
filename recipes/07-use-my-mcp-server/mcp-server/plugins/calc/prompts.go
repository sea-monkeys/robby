package main

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)

// -------------------------------------------------
//	Prompts
// -------------------------------------------------
type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Arguments   []map[string]any `json:"arguments"`
}

//go:export prompts_information
func PromptsInformation() {

	calculatorPrompt := Prompt{
		Name:        "calculator_prompt",
		Description: "A prompt to perform calculations",
		Arguments: []map[string]any{
			{
				"name":        "operation",
				"description": "The operation to perform (add, subtract, multiply, divide)",
				"type":        "string",
			},
			{
				"name":        "a",
				"description": "First number for the operation",
				"type":        "number",
			},
			{
				"name":        "b",
				"description": "Second number for the operation",
				"type":        "number",
			},
		},
	}

	prompts := []Prompt{calculatorPrompt}

	jsonData, _ := json.Marshal(prompts)
	pdk.OutputString(string(jsonData))
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Message struct {
	Role    string  `json:"role"`
	Content Content `json:"content"`
}

//go:export calculator_prompt
func CalculatorPrompt() {
	type Arguments struct {
		Operation string  `json:"operation"`
		A         float64 `json:"a"`
		B         float64 `json:"b"`
	}
	arguments := pdk.InputString()
	var args Arguments
	json.Unmarshal([]byte(arguments), &args)

	promptText := "Please perform the following calculation: " + args.Operation + " " + fmt.Sprintf("%f", args.A) + " and " + fmt.Sprintf("%f", args.B) + "."

	messages := []Message{
		{
			Role: "user",
			Content: Content{
				Type: "text",
				Text: promptText,
			},
		},
	}

	jsonData, _ := json.Marshal(messages)
	pdk.OutputString(string(jsonData))
}
