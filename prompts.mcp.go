package robby

import "fmt"

// GetPrompt retrieves a prompt by its name and arguments from the MCP client.
// It constructs a Prompt object with the name, description, and messages.
// The messages are converted from the MCP format to the internal Message format.
// If the prompt is not found or an error occurs, it returns an error.
// If the prompt is found, it returns the Prompt object.
// It requires the MCP server to be running and accessible at the specified address.
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
	//QUESTION: is there a better way to find the description in the list of prompts?

	mcpPrompt := Prompt{
		Name:        name,
		Description: description,
		Messages:    messages,
	}

	return mcpPrompt, nil
}
