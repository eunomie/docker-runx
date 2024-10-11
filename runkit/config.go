package runkit

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	lc := LocalConfig{}

	wd, err := os.Getwd()
	if err != nil {
		return lc, err
	}

	p := wd

	for {
		c, err := read(filepath.Join(p, DefaultLocalConfigFile))
		if err != nil {
			return lc, err
		}

		lc = merge(c, lc)

		if p == "/" {
			break
		}

		p = filepath.Clean(filepath.Join(p, ".."))
	}

	return lc, nil
}

func merge(a, b LocalConfig) LocalConfig {
	a.Ref = cmp.Or(b.Ref, a.Ref)
	if a.Images == nil {
		a.Images = b.Images
		return a
	}
	for imgName, img := range b.Images {
		i, ok := a.Images[imgName]
		if !ok {
			a.Images[imgName] = img
			continue
		}
		i.Default = cmp.Or(img.Default, i.Default)
		i.AllActions.Opts = mergeOpts(i.AllActions.Opts, img.AllActions.Opts)
		i.Actions = mergeActions(i.Actions, img.Actions)
		a.Images[imgName] = i
	}

	return a
}

func mergeOpts(a, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func mergeActions(a, b map[string]ConfigAction) map[string]ConfigAction {
	for k, v := range b {
		actA, ok := a[k]
		if !ok {
			a[k] = v
			continue
		}
		actA.Opts = mergeOpts(actA.Opts, v.Opts)
		a[k] = actA
	}
	return a
}

func read(filePath string) (LocalConfig, error) {
	data, err := os.ReadFile(filePath)
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
