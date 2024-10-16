package runkit

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/cli/cli/command"
	"github.com/eunomie/docker-runx/internal/constants"
)

const (
	runxConfigFile = "runx.yaml"
	runxDocFile    = "README.md"
	accessFile     = "access"
)

var subCacheDir = filepath.Join(constants.SubCommandName, "cache", "sha256")

type (
	LocalCache struct {
		cacheDir string
	}

	CacheEntry struct {
		LastAccess *time.Time
		Digest     string
		Size       int64
	}
)

func NewLocalCache(cli command.Cli) *LocalCache {
	rootDir := filepath.Dir(cli.ConfigFile().Filename)
	cacheDir := filepath.Join(rootDir, subCacheDir)

	return &LocalCache{
		cacheDir: cacheDir,
	}
}

func (c *LocalCache) Get(digest, src string) (*RunKit, error) {
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
		if err := c.writeAccessFile(digest); err != nil {
			return nil, err
		}

		rk.src = src
		return rk, nil
	}
	return nil, nil
}

func (c *LocalCache) writeAccessFile(digest string) error {
	accessDate := time.Now().Format(time.RFC3339)
	return os.WriteFile(filepath.Join(c.cacheDir, digest, accessFile), []byte(accessDate), 0o644)
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
	if err := c.writeAccessFile(digest); err != nil {
		return err
	}
	return nil
}

func (c *LocalCache) ListCache() (string, []CacheEntry, int64, error) {
	totalSize := int64(0)
	var entries []CacheEntry
	err := filepath.WalkDir(c.cacheDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != c.cacheDir {
			s, e := dirSize(path)
			if e != nil {
				return e
			}
			totalSize += s
			t := c.lastAccess(path)
			entries = append(entries, CacheEntry{
				Digest:     filepath.Base(path),
				Size:       s,
				LastAccess: t,
			})
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", nil, 0, nil
		}
		return "", nil, 0, err
	}
	return c.cacheDir, entries, totalSize, nil
}

func (c *LocalCache) EraseAll() error {
	return os.RemoveAll(c.cacheDir)
}

func (c *LocalCache) EraseNotAccessedInLast30Days() error {
	_, entries, _, err := c.ListCache()
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.LastAccess != nil && time.Since(*e.LastAccess) > 30*24*time.Hour {
			if err := os.RemoveAll(filepath.Join(c.cacheDir, e.Digest)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *LocalCache) lastAccess(path string) *time.Time {
	b, err := os.ReadFile(filepath.Join(path, accessFile))
	if err != nil {
		return nil
	}
	if t, err := time.Parse(time.RFC3339, strings.TrimSpace(string(b))); err != nil {
		return nil
	} else {
		return &t
	}
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}
