package robby

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
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

	Tools     []openai.ChatCompletionToolParam
	ToolCalls []openai.ChatCompletionMessageToolCall

	Resources []Resource
	Prompts   []Prompt

	mcpClient *mcp_golang.Client
	mcpCmd    *exec.Cmd

	lastError error
}

type AgentOption func(*Agent)

func WithDMRClient(ctx context.Context, baseURL string) AgentOption {
	return func(agent *Agent) {
		agent.ctx = ctx
		agent.dmrClient = openai.NewClient(
			option.WithBaseURL(baseURL),
			option.WithAPIKey(""),
		)
	}
}

func WithParams(params openai.ChatCompletionNewParams) AgentOption {
	return func(agent *Agent) {
		agent.Params = params
	}
}

func WithTools(tools []openai.ChatCompletionToolParam) AgentOption {
	return func(agent *Agent) {
		agent.Tools = tools
	}
}

type STDIOCommandOption []string

func WithDockerMCPToolkit() STDIOCommandOption {
	return STDIOCommandOption{
		"docker",
		"run",
		"-i",
		"--rm",
		"alpine/socat",
		"STDIO",
		"TCP:host.docker.internal:8811",
	}
}
func WithSocatMCPToolkit() STDIOCommandOption {
	return STDIOCommandOption{
		"socat",
		"STDIO",
		"TCP:host.docker.internal:8811",
	}
}

func WithMCPClient(command STDIOCommandOption) AgentOption {
	return func(agent *Agent) {

		cmd := exec.Command(
			command[0],
			command[1:]...,
		)

		stdin, err := cmd.StdinPipe()
		if err != nil {
			agent.lastError = fmt.Errorf("failed to get stdin pipe: %v", err)
			return
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			agent.lastError = fmt.Errorf("failed to get stdout pipe: %v", err)
			return
		}

		if err := cmd.Start(); err != nil {
			agent.lastError = fmt.Errorf("failed to start server: %v", err)
			return
		}

		clientTransport := stdio.NewStdioServerTransportWithIO(stdout, stdin)

		mcpClient := mcp_golang.NewClient(clientTransport)

		if _, err := mcpClient.Initialize(agent.ctx); err != nil {
			agent.lastError = fmt.Errorf("failed to initialize client: %v", err)
			return
		}
		agent.mcpClient = mcpClient
		agent.mcpCmd = cmd
	}
}

// WithMCPTools fetches the tools from the MCP server and sets them in the agent.
// It filters the tools based on the provided names and converts them to OpenAI format.
// It requires the MCP server to be running and accessible at the specified address.
// The tools are expected to be in the format defined by the MCP server.
// It returns an AgentOption that can be used to configure the agent.
// The tools are fetched using the MCP client and converted to OpenAI format.
func WithMCPTools(tools []string) AgentOption {
	return func(agent *Agent) {

		// Get the tools from the MCP client
		mcpTools, err := agent.mcpClient.ListTools(agent.ctx, nil)
		if err != nil {
			agent.lastError = err
			return
		}

		if len(tools) == 0 {
			// If no tools are specified, use all available tools
			// Convert the tools to OpenAI format
			agent.Tools = convertToOpenAITools(mcpTools.Tools)
		} else {
			filteredTools := []mcp_golang.ToolRetType{}
			for _, tool := range mcpTools.Tools {
				for _, t := range tools {
					if tool.Name == t {
						filteredTools = append(filteredTools, tool)
					}
				}
			}
			// Convert the tools to OpenAI format
			agent.Tools = convertToOpenAITools(filteredTools)
		}
	}
}

func WithMCPResources(resources []string) AgentOption {
	return func(agent *Agent) {
		// Get the resources from the MCP client
		mcpResources, err := agent.mcpClient.ListResources(agent.ctx, nil)
		if err != nil {
			agent.lastError = err
			return
		}
		resourcesList := []Resource{}
		if len(resources) == 0 {
			// If no resources are specified, use all available resources
			for _, resource := range mcpResources.Resources {
				resourcesList = append(resourcesList, Resource{
					URI:         resource.Uri,
					Name:        resource.Name,
					Description: *resource.Description,
					MimeType:    *resource.MimeType,
				})
			}

		} else {
			for _, resource := range mcpResources.Resources {
				for _, r := range resources {
					if resource.Name == r {
						resourcesList = append(resourcesList, Resource{
							URI:         resource.Uri,
							Name:        resource.Name,
							Description: *resource.Description,
							MimeType:    *resource.MimeType,
						})
					}
				}

			}
		}
		agent.Resources = resourcesList
	}
}

// TODO -> factorize the arguments conversion
func WithMCPPrompts(prompts []string) AgentOption {
	return func(agent *Agent) {
		mcpPrompts, err := agent.mcpClient.ListPrompts(agent.ctx, nil)
		if err != nil {
			agent.lastError = err
			return
		}
		promptsList := []Prompt{}
		if len(prompts) == 0 {
			// If no prompts are specified, use all available prompts
			for _, prompt := range mcpPrompts.Prompts {
				// Convert []mcp_golang.PromptSchemaArgument to []map[string]any
				args := make([]map[string]any, len(prompt.Arguments))
				for i, arg := range prompt.Arguments {
					// Marshal to JSON then unmarshal to map[string]any
					b, _ := json.Marshal(arg)
					_ = json.Unmarshal(b, &args[i])
				}
				promptsList = append(promptsList, Prompt{
					Name:        prompt.Name,
					Description: *prompt.Description,
					Arguments:   args,
				})
			}

		} else {
			for _, prompt := range mcpPrompts.Prompts {
				for _, p := range prompts {
					if prompt.Name == p {
						// Convert []mcp_golang.PromptSchemaArgument to []map[string]any
						args := make([]map[string]any, len(prompt.Arguments))
						for i, arg := range prompt.Arguments {
							// Marshal to JSON then unmarshal to map[string]any
							b, _ := json.Marshal(arg)
							_ = json.Unmarshal(b, &args[i])
						}
						promptsList = append(promptsList, Prompt{
							Name:        prompt.Name,
							Description: *prompt.Description,
							Arguments:   args,
						})
					}
				}

			}
		}
		agent.Prompts = promptsList
	}

}

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

func (agent *Agent) ChatCompletion() (string, error) {
	completion, err := agent.dmrClient.Chat.Completions.New(agent.ctx, agent.Params)

	if err != nil {
		return "", err
	}

	if len(completion.Choices) > 0 {
		return completion.Choices[0].Message.Content, nil
	} else {
		return "", errors.New("no choices found")

	}
}

func (agent *Agent) ChatCompletionStream(callBack func(self *Agent, content string, err error) error) (string, error) {
	response := ""
	stream := agent.dmrClient.Chat.Completions.NewStreaming(agent.ctx, agent.Params)
	var cbkRes error

	for stream.Next() {
		chunk := stream.Current()
		// Stream each chunk as it arrives
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			cbkRes = callBack(agent, chunk.Choices[0].Delta.Content, nil)
			response += chunk.Choices[0].Delta.Content
		}

		if cbkRes != nil {
			break
		}
	}
	if cbkRes != nil {
		return response, cbkRes
	}
	if err := stream.Err(); err != nil {
		return response, err
	}
	if err := stream.Close(); err != nil {
		return response, err
	}

	return response, nil
}

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

func (agent *Agent) ToolCallsToJSON() (string, error) {
	if len(agent.ToolCalls) == 0 {
		return "[]", nil
	}
	return ToolCallsToJSONString(agent.ToolCalls)
}

// TODO-> pagination, righ now, only resource text is returned
func (agent *Agent) ReadResource(uri string) (Resource, error) {
	mcpResourceResponse, err := agent.mcpClient.ReadResource(agent.ctx, uri)
	if err != nil {
		return Resource{}, fmt.Errorf("failed to read resource %s: %v", uri, err)
	}

	mcpResource := mcpResourceResponse.Contents[0]

	resource := Resource{}
	// search for the name and description in the agent resources
	for _, rsrc := range agent.Resources {
		if rsrc.URI == mcpResource.TextResourceContents.Uri {
			resource.Name = rsrc.Name
			resource.Description = rsrc.Description
			break
		}
	}
	//? is there a better way to find the name and description in the list of resources?

	resource.URI = mcpResource.TextResourceContents.Uri
	resource.MimeType = *mcpResource.TextResourceContents.MimeType
	resource.Text = mcpResource.TextResourceContents.Text

	return resource, nil
}

// TODO-> to be implemented
func (agent *Agent) ReadResourceByName(name string) (Resource, error) {
	return Resource{}, errors.New("not implemented yet")
}

func (agent *Agent) GetPrompt(name string, args any) (Prompt, error) {

	mcpPromptResponse, err := agent.mcpClient.GetPrompt(agent.ctx, name, args)
	if err != nil {
		return Prompt{}, fmt.Errorf("failed to get prompt %s: %v", name, err)
	}

	messages := []Message{}
	for _, msg := range mcpPromptResponse.Messages {
			messages = append(messages, Message{
				Role: string(msg.Role),
				Content: Content{
					Type: string(msg.Content.Type),
					Text: msg.Content.TextContent.Text,
				},
			})
	}

	description := ""
	for _, prompt := range agent.Prompts {
		if prompt.Name == name {
			description = prompt.Description
			break
		}
	}
	//? is there a better way to find the description in the list of prompts?

	mcpPrompt := Prompt{
		Name:        name,
		Description: description,
		Messages:    messages,
	}

	return mcpPrompt, nil
}

// --- Helpers ---

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
