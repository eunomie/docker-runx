package runkit

import (
	"context"
	"io"
	"log"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-containerregistry/pkg/registry"
)

func TestDecorateReferences(t *testing.T) {
	t.Run("wrong input source", func(t *testing.T) {
		t.Parallel()
		reg := newTestRegistry(t)
		imgDest := reg + "/test-references:" + strings.ReplaceAll(t.Name(), "/", "_")
		err := Decorate(context.TODO(), "not/a:valid/reference", imgDest, []byte("default: run"), nil)
		if err == nil {
			t.Error("should raise an error about invalid reference")
		} else if !strings.Contains(err.Error(), "could not parse reference") {
			t.Error("should raise an error about invalid reference")
		}
	})

	t.Run("wrong input source", func(t *testing.T) {
		t.Parallel()
		err := Decorate(context.TODO(), "scratch", "not/a:valid/reference", nil, nil)
		if err == nil {
			t.Error("should raise an error about invalid reference")
		} else if !strings.Contains(err.Error(), "could not parse reference") {
			t.Error("should raise an error about invalid reference")
		}
	})
}

func TestDecorateFromScratch(t *testing.T) {
	t.Run("readme only", func(t *testing.T) {
		t.Parallel()
		reg := newTestRegistry(t)
		imgDest := reg + "/test-scratch:" + strings.ReplaceAll(t.Name(), "/", "_")
		err := Decorate(context.TODO(), "scratch", imgDest, nil, []byte("# Docker Runx\n\n**Hello!**"))
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("config only", func(t *testing.T) {
		t.Parallel()
		reg := newTestRegistry(t)
		imgDest := reg + "/test-scratch:" + strings.ReplaceAll(t.Name(), "/", "_")
		err := Decorate(context.TODO(), "scratch", imgDest, []byte("default: run"), nil)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("both", func(t *testing.T) {
		t.Parallel()
		reg := newTestRegistry(t)
		imgDest := reg + "/test-scratch:" + strings.ReplaceAll(t.Name(), "/", "_")
		err := Decorate(context.TODO(), "scratch", imgDest, []byte("default: run"), []byte("# Docker Runx\n\n**Hello!**"))
		if err != nil {
			t.Error(err)
		}
	})
}

func newTestRegistry(t *testing.T) string {
	s := httptest.NewServer(
		registry.New(
			registry.WithBlobHandler(registry.NewInMemoryBlobHandler()),
			registry.Logger(log.New(io.Discard, "", log.Ldate))))
	t.Cleanup(s.Close)
	u, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	return u.Host
}
