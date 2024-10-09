package runkit

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"gopkg.in/yaml.v3"

	"github.com/eunomie/docker-runx/internal/registry"
	"github.com/eunomie/docker-runx/internal/runkit"
)

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
		config     Config
		rk         RunKit
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
		if a, ok := m.Annotations[runkit.RunxAnnotation]; ok && a == runkit.RunxManifestType {
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
		if mt == runkit.RunxConfigType {
			dataReader, err := l.Uncompressed()
			if err != nil {
				return nil, fmt.Errorf("could not read runx config %s: %w", src, err)
			}
			runxConfig, err = io.ReadAll(dataReader)
			if err != nil {
				return nil, fmt.Errorf("could not read runx config %s: %w", src, err)
			}
		} else if mt == runkit.RunxDocType {
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
		if err = yaml.Unmarshal(runxConfig, &config); err != nil {
			return nil, fmt.Errorf("could not unmarshal runx config %s: %w", src, err)
		}
		var actions []Action
		// TODO: fix reading of multiline YAML strings
		for _, a := range config.Actions {
			// a := a
			a.Command = strings.ReplaceAll(a.Command, "\n", " ")
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
