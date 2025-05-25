package robby

import (
	"context"
	"os/exec"

	mcp_golang "github.com/metoro-io/mcp-golang"

	"github.com/openai/openai-go"
)

type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
	Text        string `json:"text,omitempty"`
	Blob        string `json:"blob,omitempty"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Message struct {
	Role    string  `json:"role"`
	Content Content `json:"content"`
}

type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Arguments   []map[string]any `json:"arguments"`
	Messages    []Message        `json:"messages,omitempty"` // Optional, for storing messages related to the prompt
}

type Agent struct {
	ctx       context.Context
	dmrClient openai.Client
	Params    openai.ChatCompletionNewParams
	EmbeddingParams openai.EmbeddingNewParams

	Tools     []openai.ChatCompletionToolParam
	ToolCalls []openai.ChatCompletionMessageToolCall

	Resources []Resource
	Prompts   []Prompt

	Store MemoryVectorStore

	mcpClient *mcp_golang.Client
	mcpCmd    *exec.Cmd

	lastError error
}

type AgentOption func(*Agent)

// NewAgent creates a new Agent instance with the provided options.
// It applies all the options to the Agent and returns it.
// If any option sets an error, it returns the error instead of the Agent.
// The Agent can be configured with various options such as DMR client, parameters, tools, and memory.
func NewAgent(options ...AgentOption) (*Agent, error) {

	agent := &Agent{}
	// Apply all options
	for _, option := range options {
		option(agent)
	}
	if agent.lastError != nil {
		return nil, agent.lastError
	}
	return agent, nil
}
