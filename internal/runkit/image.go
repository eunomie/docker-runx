package runkit

import (
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/partial"
	"github.com/google/go-containerregistry/pkg/v1/static"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

const (
	RunxAnnotation         = "vnd.docker.reference.type"
	RunxManifestType       = "runx-manifest"
	RunxConfigType         = "application/vnd.runx.config+yaml"
	RunxDocType            = "application/vnd.runx.readme+txt"
	RunxFileNameAnnotation = "vnd.runx.filename"
	RunxLayerFile          = "application/vnd.runx.file.gzip"
)

func Image(runxConfig, runxDoc []byte) (v1.Image, *v1.Descriptor, error) {
	var err error

	img := mutate.MediaType(empty.Image, types.OCIManifestSchema1)
	img = mutate.ConfigMediaType(img, types.OCIConfigJSON)

	if len(runxConfig) > 0 {
		runxConfigLayer := static.NewLayer(runxConfig, RunxConfigType)
		img, err = mutate.Append(img, mutate.Addendum{
			Layer: runxConfigLayer,
		})
		if err != nil {
			return nil, nil, err
		}
	}
	if len(runxDoc) > 0 {
		runxDocLayer := static.NewLayer(runxDoc, RunxDocType)
		img, err = mutate.Append(img, mutate.Addendum{
			Layer: runxDocLayer,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	config, _ := img.ConfigFile()
	config.Architecture = "unknown"
	config.OS = "unknown"
	img, _ = mutate.ConfigFile(img, config)

	desc, _ := partial.Descriptor(img)
	desc.Platform = config.Platform()
	desc.Annotations = make(map[string]string)
	desc.Annotations[RunxAnnotation] = RunxManifestType

	return img, desc, nil
}
