package runkit

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"gopkg.in/yaml.v3"

	"github.com/eunomie/docker-runx/internal/registry"
)

type (
	Files struct {
		Files []struct {
			Name    string `yaml:"name"`
			Content string `yaml:"content"`
		} `yaml:"files"`
	}

	Cache interface {
		Get(digest, src string) (*RunKit, error)
		Set(digest string, runxConfig, runxDoc []byte) error
	}
)

func Get(ctx context.Context, cache Cache, src string) (*RunKit, error) {
	var (
		err         error
		desc        *v1.Descriptor
		indexDigest string
		cached      *RunKit
		index       v1.ImageIndex
		manifest    v1.Descriptor
		runxImg     v1.Image
		layers      []v1.Layer
		runxConfig  []byte
		runxDoc     []byte
		rk          = &RunKit{
			Files: make(map[string]string),
		}
		remoteOpts = registry.WithOptions(ctx, nil)
		ref, _     = name.ParseReference(src)
	)

	desc, err = remote.Head(ref, remoteOpts...)
	if err != nil {
		return nil, fmt.Errorf("could not get image %s: %w", src, err)
	}

	if !desc.MediaType.IsIndex() {
		return nil, fmt.Errorf("image %s can't be read by 'docker runx': should be an index", src)
	}

	indexDigest = desc.Digest.String()

	cached, err = cache.Get(indexDigest, src)
	if err == nil && cached != nil {
		return cached, nil
	}

	index, err = remote.Index(ref, remoteOpts...)
	if err != nil {
		return nil, fmt.Errorf("could not get index %s: %w", src, err)
	}

	found := false
	manifests, _ := index.IndexManifest()
	for _, m := range manifests.Manifests {
		if a, ok := m.Annotations[RunxAnnotation]; ok && a == RunxManifestType {
			manifest = m
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("image %s can't be read by 'docker runx': no runx manifest found", src)
	}
	runxImg, err = index.Image(manifest.Digest)
	if err != nil {
		return nil, fmt.Errorf("could not get runx image %s: %w", src, err)
	}

	layers, err = runxImg.Layers()
	if err != nil {
		return nil, fmt.Errorf("could not get runx layers %s: %w", src, err)
	}

	for _, l := range layers {
		mt, err := l.MediaType()
		if err != nil {
			continue
		}
		if mt == RunxConfigType {
			dataReader, err := l.Uncompressed()
			if err != nil {
				return nil, fmt.Errorf("could not read runx config %s: %w", src, err)
			}
			runxConfig, err = io.ReadAll(dataReader)
			if err != nil {
				return nil, fmt.Errorf("could not read runx config %s: %w", src, err)
			}
		} else if mt == RunxDocType {
			dataReader, err := l.Uncompressed()
			if err != nil {
				return nil, fmt.Errorf("could not read runx config %s: %w", src, err)
			}
			runxDoc, err = io.ReadAll(dataReader)
			if err != nil {
				return nil, fmt.Errorf("could not read runx config %s: %w", src, err)
			}
		}
	}

	if len(runxConfig) != 0 {
		err = decodeConfig(rk, src, runxConfig)
		if err != nil {
			return nil, err
		}
	}

	if len(runxDoc) != 0 {
		rk.Readme = string(runxDoc)
	}

	rk.src = src

	err = cache.Set(indexDigest, runxConfig, runxDoc)
	if err != nil {
		// TODO: log error
		return rk, nil
	}

	return rk, nil
}

func b64Decode(content string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(content)
}

func decodeConfig(rk *RunKit, src string, runxConfig []byte) error {
	var (
		config Config
		err    error
		files  Files
	)
	dec := yaml.NewDecoder(bytes.NewReader(runxConfig))
	dec.KnownFields(true)
	// first, the runx config itself
	if err = dec.Decode(&config); err != nil {
		return fmt.Errorf("could not decode runx config %s: %w", src, err)
	}
	// then, the optional files
	if err = dec.Decode(&files); err != nil && err != io.EOF {
		return fmt.Errorf("could not decode runx files %s: %w", src, err)
	} else {
		for _, f := range files.Files {
			c, err := b64Decode(f.Content)
			if err != nil {
				return fmt.Errorf("could not decode runx file %s: %w", f.Name, err)
			}
			rk.Files[f.Name] = string(c)
		}
	}

	if err = yaml.Unmarshal(runxConfig, &config); err != nil {
		return fmt.Errorf("could not unmarshal runx config %s: %w", src, err)
	}
	var actions []Action
	for _, a := range config.Actions {
		// TODO: fix reading of multiline YAML strings
		a.Command = strings.ReplaceAll(a.Command, "\n", " ")

		if a.Dockerfile != "" {
			if c, ok := rk.Files[a.Dockerfile]; ok {
				a.DockerfileContent = c
			}
		}

		if config.Default == a.ID {
			a.isDefault = true
		}

		for i, o := range a.Options {
			if o.Type == OptTypeNotSet {
				if len(o.Values) > 0 {
					o.Type = OptTypeSelect
				} else {
					o.Type = OptTypeInput
				}
				a.Options[i] = o
			}
		}

		actions = append(actions, a)
	}
	config.Actions = actions
	rk.Config = config

	return nil
}
