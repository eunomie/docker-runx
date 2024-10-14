package runkit

import (
	"os"
	"path/filepath"

	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/internal/constants"
)

const (
	runxConfigFile = "runx.yaml"
	runxDocFile    = "README.md"
)

var subCacheDir = filepath.Join(constants.SubCommandName, "cache", "sha256")

type LocalCache struct {
	cacheDir string
}

func NewLocalCache(cli command.Cli) *LocalCache {
	rootDir := filepath.Dir(cli.ConfigFile().Filename)
	cacheDir := filepath.Join(rootDir, subCacheDir)

	return &LocalCache{
		cacheDir: cacheDir,
	}
}

func (c *LocalCache) Get(digest string) (*RunKit, error) {
	rk := &RunKit{
		Files: make(map[string]string),
	}
	found := false

	configFile := filepath.Join(c.cacheDir, digest, runxConfigFile)
	if runxConfig, err := os.ReadFile(configFile); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		if err = decodeConfig(rk, digest, runxConfig); err != nil {
			return nil, err
		}
		found = true
	}

	readmeFile := filepath.Join(c.cacheDir, digest, runxDocFile)
	if runxDoc, err := os.ReadFile(readmeFile); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		rk.Readme = string(runxDoc)
		found = true
	}

	if found {
		return rk, nil
	}
	return nil, nil
}

func (c *LocalCache) Set(digest string, runxConfig, runxDoc []byte) error {
	digestDir := filepath.Join(c.cacheDir, digest)
	if err := os.MkdirAll(digestDir, 0o755); err != nil {
		return err
	}
	if len(runxConfig) > 0 {
		configFile := filepath.Join(c.cacheDir, digest, runxConfigFile)
		if err := os.WriteFile(configFile, runxConfig, 0o644); err != nil {
			return err
		}
	}
	if len(runxDoc) > 0 {
		readmeFile := filepath.Join(c.cacheDir, digest, runxDocFile)
		if err := os.WriteFile(readmeFile, runxDoc, 0o644); err != nil {
			return err
		}
	}
	return nil
}
