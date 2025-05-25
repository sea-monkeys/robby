package robby

import (
	"errors"
	"fmt"
)

// ReadResource retrieves a resource by its URI from the MCP client.
// It constructs a Resource object with the URI, MIME type, and text content.
// The resource name and description are searched in the agent's resources.
// If the resource is not found or an error occurs, it returns an error.
// If the resource is found, it returns the Resource object.
// It requires the MCP server to be running and accessible at the specified address.
// The resources are expected to be in the format defined by the MCP server.
func (agent *Agent) ReadResource(uri string) (Resource, error) {
	// TODO: pagination, righ now, only resource text is returned
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

func (agent *Agent) ReadResourceByName(name string) (Resource, error) {
	// TODO: to be implemented
	return Resource{}, errors.New("not implemented yet")
}
