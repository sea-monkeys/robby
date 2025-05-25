package robby

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// WithDMRClient initializes the Agent with a DMR client using the provided context and base URL.
func WithDMRClient(ctx context.Context, baseURL string) AgentOption {
	return func(agent *Agent) {
		agent.ctx = ctx
		agent.dmrClient = openai.NewClient(
			option.WithBaseURL(baseURL),
			option.WithAPIKey(""),
		)
	}
}

// WithParams sets the parameters for the Agent's chat completion requests.
func WithParams(params openai.ChatCompletionNewParams) AgentOption {
	return func(agent *Agent) {
		agent.Params = params
	}
}

// WithEmbeddingParams sets the parameters for the Agent's embedding requests.
func WithEmbeddingParams(embeddingParams openai.EmbeddingNewParams) AgentOption {
	return func(agent *Agent) {
		agent.EmbeddingParams = embeddingParams
	}
}

// WithTools sets the tools for the Agent's chat completion requests.
// It allows the Agent to use specific tools during the chat completion process.
func WithTools(tools []openai.ChatCompletionToolParam) AgentOption {
	return func(agent *Agent) {
		agent.Tools = tools
	}
}
