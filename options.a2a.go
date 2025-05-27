package robby

type A2ASettings []string


func WithA2AServer(settings A2ASettings) AgentOption {
	return func(agent *Agent) {
		//agent.a2aSettings = settings
	}
}