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
	ctx             context.Context
	dmrClient       openai.Client
	Params          openai.ChatCompletionNewParams
	EmbeddingParams openai.EmbeddingNewParams

	Tools     []openai.ChatCompletionToolParam
	ToolCalls []openai.ChatCompletionMessageToolCall

	Resources []Resource
	Prompts   []Prompt

	Store MemoryVectorStore

	// --- MCP Client ---
	mcpClient *mcp_golang.Client
	mcpCmd    *exec.Cmd
	// QUESTION: is it a good idea to make an agent a MCP Server?

	// --- A2A Server ---
	AgentCard AgentCard

	AgentCallback func(taskRequest TaskRequest) (TaskResponse, error)

	// --- A2A Client ---

	lastError error
}

// --- BEGIN: A2A Protocol ---

// AgentCard represents the metadata for this agent
type AgentCard struct {
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	URL          string           `json:"url"`
	Version      string           `json:"version"`
	Capabilities map[string]any   `json:"capabilities"`
	Skills       []map[string]any `json:"skills,omitempty"` // Optional, for storing skills related to the agent
}

type AgentMessageParams struct {
	Message  AgentMessage   `json:"message"`
	MetaData map[string]any `json:"metadata,omitempty"` // Optional, for additional metadata
}

// QUESTION: This is a message, shoul I change the name
// REF: https://google-a2a.github.io/A2A/specification/#92-basic-execution-synchronous-polling-style
// TaskRequest represents an incoming A2A task request
type TaskRequest struct {
	JSONRpcVersion string             `json:"jsonrpc"` // Should be "2.0"
	ID             string             `json:"id"`
	Params         AgentMessageParams `json:"params"`
	Method         string             `json:"method,omitempty"` // Optional, for specifying the method of the task
}

// Message represents a message structure
type AgentMessage struct {
	Role  string     `json:"role,omitempty"`
	Parts []TextPart `json:"parts"`
	MessageID string `json:"messageId,omitempty"` // Optional, for storing message ID
	TaskID     string `json:"taskId,omitempty"`    // Optional, for storing task ID
	ContextID string `json:"contextId,omitempty"` // Optional, for storing context ID
}

// TextPart represents a text part of a message
type TextPart struct {
	Text string `json:"text"`
	Type string `json:"type"` // Should be "text" for text parts
}

// TaskStatus represents the status of a task
type TaskStatus struct {
	State string `json:"state"`
}

// TODO: make the response compliant with the A2A protocol
// REF: https://google-a2a.github.io/A2A/specification/#92-basic-execution-synchronous-polling-style

type Artifact struct {
	ArtifactID string     `json:"artifactId"`
	Name       string     `json:"name"`
	Parts      []TextPart `json:"parts"` // Parts of the artifact, e.g., text, images, etc.
}

type Result struct {
	ID        string          `json:"id"`
	ContextID string          `json:"contextId"`
	Status    TaskStatus      `json:"status"`
	Artifacts []Artifact      `json:"artifacts,omitempty"` // Optional, for storing artifacts related to the task
	History   []AgentMessage  `json:"history,omitempty"`   // Optional, for storing message history related to the task
	Kind      string          `json:"kind"` // Should be "task"
	Metadata  map[string]any `json:"metadata,omitempty"` // Optional, for additional metadata
}

// TaskResponse represents the response task structure
type TaskResponse struct {
	JSONRpcVersion string             `json:"jsonrpc"` // Should be "2.0"
	ID       string         `json:"id"`
	Result   Result         `json:"result"` // The result of the task execution

}


// --- END: A2A Protocol ---

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
