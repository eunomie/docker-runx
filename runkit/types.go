package runkit

type (
	RunKit struct {
		Config Config
		Readme string

		src string
	}

	Config struct {
		Actions []Action `yaml:"actions" json:"actions"`
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
