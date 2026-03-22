package sessiondomain

// AgentService handles agent responses to user messages.
type AgentService interface {
	// Respond processes a message and generates a response using the LLM agent.
	Respond(ctx *MessageContext) error
}
