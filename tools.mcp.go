package robby

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/openai/openai-go"
)

// ExecuteMCPToolCalls executes the tool calls detected by the Agent using the MCP client.
// It takes no additional parameters as it uses the Agent's context and MCP client.
// It returns a slice of responses from the executed tools or an error if any tool call fails.
// It also appends the tool responses to the Agent's messages for further processing.
// This function is specifically designed to work with the MCP toolkit, which allows for tool calls
// to be executed in a remote environment using the MCP protocol.
// It assumes that the Agent has been initialized with an MCP client and the necessary context.
// The function iterates over the Agent's ToolCalls, unmarshals the arguments for each tool call,
// and then calls the tool using the MCP client. The responses are collected and returned.
// If no tool responses are found, it returns an error.
// It is important to ensure that the MCP client is properly configured and connected to the MCP server
// before calling this function, as it relies on the MCP protocol for executing tool calls.
// It is a synchronous operation that waits for the completion of each tool call.
func (agent *Agent) ExecuteMCPToolCalls() ([]string, error) {

	responses := []string{}
	for _, toolCall := range agent.ToolCalls {

		var args map[string]any
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
		if err != nil {
			return nil, err
		}

		// Call the tool with the arguments thanks to the MCP client
		toolResponse, err := agent.mcpClient.CallTool(agent.ctx, toolCall.Function.Name, args)
		if err != nil {
			responses = append(responses, fmt.Sprintf("%v", err))
		} else {
			if toolResponse != nil && len(toolResponse.Content) > 0 && toolResponse.Content[0].TextContent != nil {

				agent.Params.Messages = append(
					agent.Params.Messages,
					openai.ToolMessage(
						toolResponse.Content[0].TextContent.Text,
						toolCall.ID,
					),
				)
				responses = append(responses, toolResponse.Content[0].TextContent.Text)
			}
		}

	}
	if len(responses) == 0 {
		return nil, errors.New("no tool responses found")
	}
	return responses, nil
}
