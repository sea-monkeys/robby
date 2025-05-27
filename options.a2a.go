package robby


type A2AServerSettings map[string]any
type A2AServerSkills []map[string]any

//type AgentCard map[string]any

func WithA2AServer(settings A2AServerSettings, skills A2AServerSkills) AgentOption {
	return func(agent *Agent) {
		agent.AgentCard = AgentCard{
			Name:        settings["name"].(string),
			Description: settings["description"].(string),
			URL:         settings["url"].(string),
			Version:     settings["version"].(string),
			Capabilities: map[string]interface{}{
				"streaming": false,
				"pushNotifications": false,
			},
		}
		agent.AgentCard.Skills = skills

	}
}	
/*
			Streaming:         false, // Not implementing subscribe
			PushNotifications: false,
*/