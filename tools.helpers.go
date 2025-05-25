package robby

import (
	"encoding/json"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/openai/openai-go"
)

// ToolCallsToJSONString converts a slice of openai.ChatCompletionMessageToolCall to a JSON string.
// It extracts the tool call ID and function arguments, converting them to a generic interface
// for JSON marshaling. The resulting JSON string is formatted with indentation for readability.
// If the tool calls are empty, it returns an empty JSON array.
// If any error occurs during the conversion, it returns an error.
// The function is useful for logging or storing tool calls in a structured format.
// It returns a JSON string representation of the tool calls.
func ToolCallsToJSONString(tools []openai.ChatCompletionMessageToolCall) (string, error) {
	var jsonData []any

	// Convert tools to generic interface
	for _, tool := range tools {
		var args any
		if err := json.Unmarshal([]byte(tool.Function.Arguments), &args); err != nil {
			return "", err
		}

		jsonData = append(jsonData, map[string]any{
			"id": tool.ID,
			"function": map[string]any{
				"name":      tool.Function.Name,
				"arguments": args,
			},
		})
	}

	// Marshal back to JSON with indentation
	jsonString, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}

// convertToOpenAITools converts a slice of mcp_golang.ToolRetType to a slice of openai.ChatCompletionToolParam.
// It extracts the tool name, description, and input schema from each tool.
// The input schema is expected to be a map[string]any, which is converted to the appropriate OpenAI function parameters format.
// The resulting slice of openai.ChatCompletionToolParam can be used in OpenAI API calls for tool usage.
// It returns a slice of openai.ChatCompletionToolParam containing the converted tools.
func convertToOpenAITools(tools []mcp_golang.ToolRetType) []openai.ChatCompletionToolParam {
	openAITools := make([]openai.ChatCompletionToolParam, len(tools))

	for i, tool := range tools {
		schema := tool.InputSchema.(map[string]any)
		openAITools[i] = openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        tool.Name,
				Description: openai.String(*tool.Description),
				Parameters: openai.FunctionParameters{
					"type":       "object",
					"properties": schema["properties"],
					"required":   schema["required"],
				},
			},
		}
	}
	return openAITools
}
