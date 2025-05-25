package robby

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/openai/openai-go"
)

// ExecuteToolCalls executes the tool calls detected by the Agent.
// It takes a map of tool implementations where the key is the tool name and the value is a function that implements the tool.
// Each tool function should accept a map of arguments and return a response or an error.
// The function returns a slice of responses from the executed tools or an error if any tool call fails.
// It also appends the tool responses to the Agent's messages for further processing
func (agent *Agent) ExecuteToolCalls(toolsImpl map[string]func(any) (any, error)) ([]string, error) {
	responses := []string{}
	for _, toolCall := range agent.ToolCalls {
		// Check if the tool is implemented
		toolFunc, ok := toolsImpl[toolCall.Function.Name]
		if !ok {
			return nil, fmt.Errorf("tool %s not implemented", toolCall.Function.Name)
		}

		var args map[string]any
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
		if err != nil {
			return nil, err
		}

		// Call the tool with the arguments
		toolResponse, err := toolFunc(args)
		if err != nil {
			responses = append(responses, fmt.Sprintf("%v", err))
		} else {
			responses = append(responses, fmt.Sprintf("%v", toolResponse))
			agent.Params.Messages = append(
				agent.Params.Messages,
				openai.ToolMessage(
					fmt.Sprintf("%v", toolResponse),
					toolCall.ID,
				),
			)
		}
	}
	if len(responses) == 0 {
		return nil, errors.New("no tool responses found")
	}
	return responses, nil
}

// ToolCallsToJSON converts the Agent's tool calls to a JSON string.
// If there are no tool calls, it returns an empty JSON array.
// It uses the ToolCallsToJSONString function to convert the tool calls to a JSON string format.
func (agent *Agent) ToolCallsToJSON() (string, error) {
	if len(agent.ToolCalls) == 0 {
		return "[]", nil
	}
	return ToolCallsToJSONString(agent.ToolCalls)
}
