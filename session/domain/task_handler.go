package sessiondomain

// TaskHandler handles automatic task activations (cron, timer)
type TaskHandler interface {
	OnActivation(ctx *ActivationContext)
}
