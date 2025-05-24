package main

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)

// -------------------------------------------------
//  Tools
// -------------------------------------------------
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string         `json:"type"`
	Required   []string       `json:"required"`
	Properties map[string]any `json:"properties"`
}

/*
  Defining a set of 4 tools:
  - addition
  - subtraction
  - multiplication
  - division
  Each tool has a name, description, and input schema.
  The input schema defines the type of input the tool expects, the required fields, and the properties of each field.
  The tools are defined in the ToolsInformation function, which is called by the MCP "runner" to get the tool information.
*/

//go:export tools_information
func ToolsInformation() {

	addition := Tool{
		Name:        "add",
		Description: "add two numbers",
		InputSchema: InputSchema{
			Type:     "object",
			Required: []string{"a", "b"},
			Properties: map[string]any{
				"a": map[string]any{
					"type":        "number",
					"description": "First number to add",
				},
				"b": map[string]any{
					"type":        "number",
					"description": "Second number to add",
				},
			},
		},
	}

	subtraction := Tool{
		Name:        "subtract",
		Description: "subtract two numbers",
		InputSchema: InputSchema{
			Type:     "object",
			Required: []string{"a", "b"},
			Properties: map[string]any{
				"a": map[string]any{
					"type":        "number",
					"description": "First number to subtract",
				},
				"b": map[string]any{
					"type":        "number",
					"description": "Second number to subtract",
				},
			},
		},
	}

	multiplication := Tool{
		Name:        "multiply",
		Description: "multiply two numbers",
		InputSchema: InputSchema{
			Type:     "object",
			Required: []string{"a", "b"},
			Properties: map[string]any{
				"a": map[string]any{
					"type":        "number",
					"description": "First number to multiply",
				},
				"b": map[string]any{
					"type":        "number",
					"description": "Second number to multiply",
				},
			},
		},
	}

	division := Tool{
		Name:        "divide",
		Description: "divide two numbers",
		InputSchema: InputSchema{
			Type:     "object",
			Required: []string{"a", "b"},
			Properties: map[string]any{
				"a": map[string]any{
					"type":        "number",
					"description": "First number to divide",
				},
				"b": map[string]any{
					"type":        "number",
					"description": "Second number to divide",
				},
			},
		},
	}

	tools := []Tool{addition, subtraction, multiplication, division}

	jsonData, _ := json.Marshal(tools)
	pdk.OutputString(string(jsonData))
}

/*
  Implementation of the four operarions
  - addition
  - subtraction
  - multiplication
  - division
  Each operation takes two numbers as input and returns the result.
  The input is expected to be in JSON format, with the keys "a" and "b" representing the two numbers.
  The result is returned as a string.
  The functions are exported using the //go:export directive, which allows them to be called from the MCP.
*/

//go:export add
func Add() {
	type Arguments struct {
		A float64 `json:"a"`
		B float64 `json:"b"`
	}
	arguments := pdk.InputString()
	var args Arguments
	json.Unmarshal([]byte(arguments), &args)

	result := args.A + args.B
	pdk.OutputString(fmt.Sprintf("%v", result))
}

//go:export subtract
func Subtract() {
	type Arguments struct {
		A float64 `json:"a"`
		B float64 `json:"b"`
	}
	arguments := pdk.InputString()
	var args Arguments
	json.Unmarshal([]byte(arguments), &args)

	result := args.A - args.B
	pdk.OutputString(fmt.Sprintf("%v", result))
}

//go:export multiply
func Multiply() {
	type Arguments struct {
		A float64 `json:"a"`
		B float64 `json:"b"`
	}
	arguments := pdk.InputString()
	var args Arguments
	json.Unmarshal([]byte(arguments), &args)

	result := args.A * args.B
	pdk.OutputString(fmt.Sprintf("%v", result))
}

//go:export divide
func Divide() {
	type Arguments struct {
		A float64 `json:"a"`
		B float64 `json:"b"`
	}
	arguments := pdk.InputString()
	var args Arguments
	json.Unmarshal([]byte(arguments), &args)

	if args.B == 0 {
		pdk.OutputString("Error: Division by zero")
		return
	}

	result := args.A / args.B
	pdk.OutputString(fmt.Sprintf("%v", result))
}
