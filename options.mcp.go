package robby

import (
	"encoding/json"
	"fmt"
	"os/exec"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

type STDIOCommandOption []string

// WithDockerMCPToolkit returns a STDIOCommandOption that runs the MCP toolkit using Docker.
// It uses the Alpine image with Socat to connect to the MCP server running on host.docker.internal:8811.
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

// WithSocatMCPToolkit returns a STDIOCommandOption that runs the MCP toolkit using Socat.
// It connects to the MCP server running on host.docker.internal:8811.
func WithSocatMCPToolkit() STDIOCommandOption {
	return STDIOCommandOption{
		"socat",
		"STDIO",
		"TCP:host.docker.internal:8811",
	}
}

// WithMCPClient initializes the Agent with an MCP client using the provided command.
// It runs the command to connect to the MCP server and sets up the client transport.
// The command should be a valid command that can be executed in the environment where the agent runs.
// It returns an AgentOption that can be used to configure the agent.
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

// WithMCPResources fetches the resources from the MCP server and sets them in the agent.
// It filters the resources based on the provided names and converts them to a Resource format.
// It requires the MCP server to be running and accessible at the specified address.
// The resources are expected to be in the format defined by the MCP server.
// It returns an AgentOption that can be used to configure the agent.
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


// WithMCPPrompts fetches the prompts from the MCP server and sets them in the agent.
// It filters the prompts based on the provided names and converts them to a Prompt format.
// It requires the MCP server to be running and accessible at the specified address.
// The prompts are expected to be in the format defined by the MCP server.
// It returns an AgentOption that can be used to configure the agent.
func WithMCPPrompts(prompts []string) AgentOption {
	// TODO: -> factorize the arguments conversion
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
