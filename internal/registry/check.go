package registry

import (
	"context"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func ImageExist(ctx context.Context, src string) bool {
	var (
		err        error
		remoteOpts = WithOptions(ctx, nil)
		ref, _     = name.ParseReference(src)
	)

	_, err = remote.Head(ref, remoteOpts...)
	return err == nil
}
