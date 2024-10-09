package runkit

type (
	RunKit struct {
		Actions []Action `yaml:"actions" json:"actions"`
		src     string
	}

	Action struct {
		ID      string     `yaml:"id" json:"id"`
		Type    ActionType `yaml:"type" json:"type"`
		Command string     `yaml:"cmd" json:"cmd,omitempty"`
	}

	ActionType string
)

const (
	ActionTypeRun ActionType = "run"
)
