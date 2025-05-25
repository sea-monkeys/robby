package robby

import (
	"errors"

	"github.com/openai/openai-go"
)

// ToolsCompletion handles the tool calls completion request using the DMR client.
// It sends the parameters set in the Agent and returns the detected tool calls or an error.
// It is a synchronous operation that waits for the completion to finish.
func (agent *Agent) ToolsCompletion() ([]openai.ChatCompletionMessageToolCall, error) {

	agent.Params.Tools = agent.Tools

	completion, err := agent.dmrClient.Chat.Completions.New(agent.ctx, agent.Params)
	if err != nil {
		return nil, err
	}
	detectedToolCalls := completion.Choices[0].Message.ToolCalls
	if len(detectedToolCalls) == 0 {
		return nil, errors.New("no tool calls detected")
	}
	agent.ToolCalls = detectedToolCalls

	return detectedToolCalls, nil
}
