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

type Files struct {
	Files []struct {
		Name    string `yaml:"name"`
		Content string `yaml:"content"`
	} `yaml:"files"`
}

func Get(ctx context.Context, src string) (*RunKit, error) {
	var (
		err        error
		index      v1.ImageIndex
		desc       *remote.Descriptor
		manifest   v1.Descriptor
		runxImg    v1.Image
		layers     []v1.Layer
		runxConfig []byte
		runxDoc    []byte
		files      Files
		config     Config
		rk         = RunKit{
			Files: make(map[string]string),
		}
		remoteOpts = registry.WithOptions(ctx, nil)
		ref, _     = name.ParseReference(src)
	)

	desc, err = remote.Get(ref, remoteOpts...)
	if err != nil {
		return nil, fmt.Errorf("could not get image %s: %w", src, err)
	}

	if !desc.MediaType.IsIndex() {
		return nil, fmt.Errorf("image %s can't be read by 'docker runx': should be an index", src)
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
		dec := yaml.NewDecoder(bytes.NewReader(runxConfig))
		dec.KnownFields(true)
		// first, the runx config itself
		if err = dec.Decode(&config); err != nil {
			return nil, fmt.Errorf("could not decode runx config %s: %w", src, err)
		}
		// then, the optional files
		if err = dec.Decode(&files); err != nil && err != io.EOF {
			return nil, fmt.Errorf("could not decode runx files %s: %w", src, err)
		} else {
			for _, f := range files.Files {
				c, err := b64Decode(f.Content)
				if err != nil {
					return nil, fmt.Errorf("could not decode runx file %s: %w", f.Name, err)
				}
				rk.Files[f.Name] = string(c)
			}
		}

		if err = yaml.Unmarshal(runxConfig, &config); err != nil {
			return nil, fmt.Errorf("could not unmarshal runx config %s: %w", src, err)
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
			actions = append(actions, a)
		}
		config.Actions = actions
		rk.Config = config
	}

	if len(runxDoc) != 0 {
		rk.Readme = string(runxDoc)
	}

	rk.src = src

	return &rk, nil
}

func b64Decode(content string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(content)
}
