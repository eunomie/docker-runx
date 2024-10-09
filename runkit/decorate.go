package runkit

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/match"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/types"

	"github.com/eunomie/docker-runx/internal/registry"
	"github.com/eunomie/docker-runx/internal/runkit"
)

func Decorate(ctx context.Context, src, dest string, runxConfig, runxDoc []byte) error {
	var (
		index                    v1.ImageIndex
		desc                     *remote.Descriptor
		remoteOpts               = registry.WithOptions(ctx, nil)
		ref, _                   = name.ParseReference(src)
		destRef, _               = name.ParseReference(dest)
		runxImage, runxDesc, err = runkit.Image(runxConfig, runxDoc)
	)

	if err != nil {
		return fmt.Errorf("could not create runx image: %w", err)
	}

	desc, err = remote.Get(ref, remoteOpts...)
	if err != nil {
		return fmt.Errorf("could not get image %s: %w", src, err)
	}

	if desc.MediaType.IsImage() {
		img, err := remote.Image(ref, remoteOpts...)
		if err != nil {
			return fmt.Errorf("could not get image %s: %w", src, err)
		}
		configFile, _ := img.ConfigFile()
		imgDesc := desc.Descriptor
		imgDesc.Platform = configFile.Platform()

		index = // create a manifest
		mutate.AppendManifests(
			// as an index
			mutate.IndexMediaType(empty.Index, types.OCIImageIndex),
			// with the referenced image
			mutate.IndexAddendum{
				Add:        img,
				Descriptor: imgDesc,
			},
			// and the new runx image
			mutate.IndexAddendum{
				Add:        runxImage,
				Descriptor: *runxDesc,
			})
	} else if desc.MediaType.IsIndex() {
		index, err = remote.Index(ref, remoteOpts...)
		if err != nil {
			return fmt.Errorf("could not get image index %s: %w", src, err)
		}

		// remove existing runx manifest
		manifests, _ := index.IndexManifest()
		for _, manifest := range manifests.Manifests {
			if _, ok := manifest.Annotations[runkit.RunxManifestType]; ok {
				index = mutate.RemoveManifests(index, match.Digests(manifest.Digest))
			}
		}

		// add the new runx manifest
		index = mutate.AppendManifests(
			index,
			mutate.IndexAddendum{
				Add:        runxImage,
				Descriptor: *runxDesc,
			})
	} else {
		return fmt.Errorf("unsupported media type %s", desc.MediaType)
	}

	err = remote.WriteIndex(destRef, index, remoteOpts...)
	if err != nil {
		return fmt.Errorf("could not write index %s: %w", dest, err)
	}

	return nil
}
