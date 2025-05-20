package robby

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

type Agent struct {
	ctx       context.Context
	dmrClient openai.Client
	mcpClient *mcp_golang.Client
	mcpCmd    *exec.Cmd
	//messages   []openai.ChatCompletionMessageParamUnion
	Params openai.ChatCompletionNewParams
	Tools  []openai.ChatCompletionToolParam

	ToolCalls []openai.ChatCompletionMessageToolCall

	lastError error
	//lastResult any

	//tools      *[]openai.ChatCompletionToolParam
}

type AgentOption func(*Agent)

// WithHost sets the server host
func WithDMRClient(ctx context.Context, baseURL string) AgentOption {
	return func(agent *Agent) {
		agent.ctx = ctx
		agent.dmrClient = openai.NewClient(
			option.WithBaseURL(baseURL),
			option.WithAPIKey(""),
		)
	}
}

// Add default value for the url?
// or WitDMRClientFromContainer

func WithParams(params openai.ChatCompletionNewParams) AgentOption {
	return func(agent *Agent) {
		agent.Params = params
	}
}

// TODO: To be implemented
func WithSTDIOMCPClient() AgentOption {
	return func(agent *Agent) {}
}

type STDIOCommandOption []string

func WithDocker() STDIOCommandOption {
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
func WithSocat() STDIOCommandOption {
	return STDIOCommandOption{
		"socat",
		"STDIO",
		"TCP:host.docker.internal:8811",
	}
}

func WithMCPToolkitClient(command STDIOCommandOption) AgentOption {
	return func(agent *Agent) {

		cmd := exec.Command(
			command[0],
			command[1:]...,
		)

		stdin, err := cmd.StdinPipe()
		if err != nil {
			agent.lastError = fmt.Errorf("😡 failed to get stdin pipe: %v", err)
			return
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			agent.lastError = fmt.Errorf("😡 failed to get stdout pipe: %v", err)
			return
		}

		if err := cmd.Start(); err != nil {
			agent.lastError = fmt.Errorf("😡 failed to start server: %v", err)
			return
		}

		clientTransport := stdio.NewStdioServerTransportWithIO(stdout, stdin)

		mcpClient := mcp_golang.NewClient(clientTransport)

		if _, err := mcpClient.Initialize(agent.ctx); err != nil {
			agent.lastError = fmt.Errorf("😡 failed to initialize client: %v", err)
			return
		}
		agent.mcpClient = mcpClient
		agent.mcpCmd = cmd
	}
}

// TODO-> To be implemented
func WithTools(tools []openai.ChatCompletionToolParam) AgentOption {
	return func(agent *Agent) {
		agent.Tools = tools
	}
}

func WithMCPTools(tools []string) AgentOption {
	return func(agent *Agent) {

		// Get the tools from the MCP client
		mcpTools, err := agent.mcpClient.ListTools(agent.ctx, nil)
		if err != nil {
			agent.lastError = err
			return
		}

		//fmt.Println("✋🛠️ Tools: ", mcpTools.Tools)

		// Convert the tools to OpenAI format
		filteredTools := []mcp_golang.ToolRetType{}
		for _, tool := range mcpTools.Tools {
			for _, t := range tools {
				if tool.Name == t {
					filteredTools = append(filteredTools, tool)
				}
			}
		}
		/*
			fmt.Println("✋🛠️ Filtered Tools:")
			for _, tool := range filteredTools {
				fmt.Println("✋🛠️ Tool: ", tool.Name)
				fmt.Println("✋🛠️ Description: ", *tool.Description)
				fmt.Println("✋🛠️ Schema: ", tool.InputSchema)
			}
		*/
		agent.Tools = convertToOpenAITools(filteredTools)
	}
}

// Apply allows adding options to an existing Agent instance
func (agent *Agent) Apply(options ...AgentOption) error {
	for _, option := range options {
		option(agent)
	}
	if agent.lastError != nil {
		return agent.lastError
	}
	return nil
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
	// Check if the MCP client is initialized
	//if agent.mcpClient == nil {
	//	return nil, errors.New("MCP client is not initialized")
	//}

	agent.Params.Tools = agent.Tools

	completion, err := agent.dmrClient.Chat.Completions.New(agent.ctx, agent.Params)
	if err != nil {
		return nil, err
	}
	// TODO: add a detected tool call property, or return the tool calls
	detectedToolCalls := completion.Choices[0].Message.ToolCalls
	if len(detectedToolCalls) == 0 {
		return nil, errors.New("no tool calls detected")
	}
	agent.ToolCalls = detectedToolCalls

	return detectedToolCalls, nil
}

// TODO-> to be implemented
func (agent *Agent) ToolsLoopCompletion() ([]string, error) {
	return nil, nil
}

func (agent *Agent) ExecuteToolCalls(toolsImpl map[string]func(any) any) ([]string, error) {
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
		response := toolFunc(args)
		responses = append(responses, fmt.Sprintf("%v", response))
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

		// Call the tool with the arguments
		toolResponse, err := agent.mcpClient.CallTool(agent.ctx, toolCall.Function.Name, args)
		if err != nil {
			return nil, err //? should I return the error? == stop here?
		}
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
	if len(responses) == 0 {
		return nil, errors.New("no tool responses found")
	}
	return responses, nil
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

func ToolCallsToJSONString(tools []openai.ChatCompletionMessageToolCall) (string, error) {
	// Create a temporary structure to handle the proper JSON marshaling
	type Function struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}

	type ToolCall struct {
		ID       string   `json:"id"`
		Function Function `json:"function"`
	}

	// Convert the tools to our temporary structure
	tempTools := make([]ToolCall, len(tools))
	for i, tool := range tools {
		// Parse the arguments string as raw JSON
		var rawArgs json.RawMessage
		if err := json.Unmarshal([]byte(tool.Function.Arguments), &rawArgs); err != nil {
			return "", err
		}

		tempTools[i] = ToolCall{
			ID: tool.ID,
			Function: Function{
				Name:      tool.Function.Name,
				Arguments: rawArgs,
			},
		}
	}

	// Marshal with indentation
	jsonString, err := json.MarshalIndent(tempTools, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}
