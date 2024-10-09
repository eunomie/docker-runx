package runkit

type (
	RunKit struct {
		Config Config
		Readme string

		src string
	}

	Config struct {
		Default string   `yaml:"default,omitempty" json:"default,omitempty"`
		Actions []Action `yaml:"actions" json:"actions"`
	}

	Action struct {
		ID      string     `yaml:"id" json:"id"`
		Desc    string     `yaml:"desc,omitempty" json:"desc,omitempty"`
		Type    ActionType `yaml:"type" json:"type"`
		Command string     `yaml:"cmd" json:"cmd,omitempty"`
		Env     []string   `yaml:"env,omitempty" json:"env,omitempty"`
		Options []Opt      `yaml:"opts,omitempty" json:"opts,omitempty"`
	}

	Opt struct {
		Name        string   `yaml:"name" json:"name"`
		Description string   `yaml:"desc" json:"desc,omitempty"`
		Prompt      string   `yaml:"prompt,omitempty" json:"prompt,omitempty"`
		Required    bool     `yaml:"required,omitempty" json:"required,omitempty"`
		Values      []string `yaml:"values,omitempty" json:"values,omitempty"`
	}

	ActionType string

	LocalConfig struct {
		Ref    string                 `yaml:"ref,omitempty" json:"ref,omitempty"`
		Images map[string]ConfigImage `yaml:"images,omitempty" json:"images,omitempty"`
	}

	ConfigImage struct {
		Default string `yaml:"default,omitempty" json:"default,omitempty"`
	}
)

const (
	ActionTypeRun ActionType = "run"
)
