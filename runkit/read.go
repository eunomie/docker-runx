package runkit

import (
	"context"
	"fmt"
	"io"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"gopkg.in/yaml.v2"

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
		actions    RunKit
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
		if mt, err := l.MediaType(); err == nil && mt == runkit.RunxConfigType {
			dataReader, err := l.Uncompressed()
			if err != nil {
				return nil, fmt.Errorf("could not read runx actions %s: %w", src, err)
			}
			runxConfig, err = io.ReadAll(dataReader)
			if err != nil {
				return nil, fmt.Errorf("could not read runx actions %s: %w", src, err)
			}
			break
		}
	}

	if len(runxConfig) == 0 {
		return nil, fmt.Errorf("image %s can't be read by 'docker runx': no runx actions found", src)
	}

	if err = yaml.Unmarshal(runxConfig, &actions); err != nil {
		return nil, fmt.Errorf("could not unmarshal runx actions %s: %w", src, err)
	}

	actions.src = src

	return &actions, nil
}
