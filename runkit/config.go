package runkit

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

const (
	DefaultLocalConfigFile = ".docker/runx.yaml"
)

var (
	localConfig LocalConfig
	readOnce    = sync.Once{}
)

func GetLocalConfig() LocalConfig {
	readOnce.Do(func() {
		lc, err := getLocalConfig()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error reading local config: %v\n", err)
			os.Exit(1)
		}
		localConfig = lc
	})
	return localConfig
}

func getLocalConfig() (LocalConfig, error) {
	data, err := os.ReadFile(DefaultLocalConfigFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// in this case we just don't have a config file, so use everything as default
			return LocalConfig{}, nil
		}
		return LocalConfig{}, err
	}

	var config LocalConfig
	if err = yaml.Unmarshal(data, &config); err != nil {
		return LocalConfig{}, err
	}

	return config, nil
}
