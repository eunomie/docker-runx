package runkit

type (
	RunKit struct {
		Config Config
		Readme string
		Files  map[string]string

		src string
	}

	Config struct {
		Default string   `yaml:"default,omitempty" json:"default,omitempty"`
		Actions []Action `yaml:"actions" json:"actions"`
	}

	Action struct {
		isDefault         bool
		ID                string            `yaml:"id" json:"id"`
		Desc              string            `yaml:"desc,omitempty" json:"desc,omitempty"`
		Type              ActionType        `yaml:"type" json:"type"`
		Command           string            `yaml:"cmd" json:"cmd,omitempty"`
		Env               []string          `yaml:"env,omitempty" json:"env,omitempty"`
		Options           []Opt             `yaml:"opts,omitempty" json:"opts,omitempty"`
		Shell             map[string]string `yaml:"shell,omitempty" json:"shell,omitempty"`
		Dockerfile        string            `yaml:"dockerfile,omitempty" json:"dockerfile,omitempty"`
		DockerfileContent string
	}

	Opt struct {
		Name        string   `yaml:"name" json:"name"`
		Type        OptType  `yaml:"type,omitempty" json:"type,omitempty"`
		Description string   `yaml:"desc" json:"desc,omitempty"`
		NoPrompt    bool     `yaml:"no-prompt,omitempty" json:"no-prompt,omitempty"`
		Prompt      string   `yaml:"prompt,omitempty" json:"prompt,omitempty"`
		Required    bool     `yaml:"required,omitempty" json:"required,omitempty"`
		Values      []string `yaml:"values,omitempty" json:"values,omitempty"`
		Default     string   `yaml:"default,omitempty" json:"default,omitempty"`
	}

	ActionType string

	OptType string

	LocalConfig struct {
		Ref    string                 `yaml:"ref,omitempty" json:"ref,omitempty"`
		Images map[string]ConfigImage `yaml:"images,omitempty" json:"images,omitempty"`
	}

	ConfigImage struct {
		Default    string                  `yaml:"default,omitempty" json:"default,omitempty"`
		AllActions ConfigAction            `yaml:"all-actions,omitempty" json:"all-actions,omitempty"`
		Actions    map[string]ConfigAction `yaml:"actions,omitempty" json:"actions,omitempty"`
	}

	ConfigAction struct {
		Opts map[string]string `yaml:"opts,omitempty" json:"opts,omitempty"`
	}
)

const (
	ActionTypeRun   ActionType = "run"
	ActionTypeBuild ActionType = "build"

	OptTypeNotSet  OptType = ""
	OptTypeInput   OptType = "input"
	OptTypeSelect  OptType = "select"
	OptTypeConfirm OptType = "confirm"
)

func (action *Action) IsDefault() bool {
	return action.isDefault
}
