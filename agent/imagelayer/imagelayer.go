package imagelayer

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/rs/zerolog/log"
)

type ImageLayers interface {
	FetchImageLayers(imageRef, outDir string) error
	DownloadImageLayers(ctx context.Context, imageRef, outDir string) error
}

type ImageLayerHandler struct {
	ImageLayer *ImageLayers
}

func (i *ImageLayerHandler) FetchImageLayers(imageRef, outDir string) error {
	if err := i.DownloadImageLayers(context.Background(), imageRef, outDir); err != nil {
		log.Err(err).Msg("error while calling DownloadImageLayers()")
		return err
	}
	log.Info().Msg("Image Layers properly downloaded")
	return nil
}

func (i *ImageLayerHandler) DownloadImageLayers(ctx context.Context, imageRef, outDir string) error {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		log.Err(err).Msg("parse ref error:")
		return err
	}
	img, err := remote.Image(ref)
	if err != nil {
		log.Err(err).Msg("pull image")
		return err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Err(err).Msg("error while calling os.MkdirAll()")
		return err
	}

	layers, err := img.Layers()
	if err != nil {
		log.Err(err).Msg("error while calling img.Layers()")
		return err
	}

	for i, layer := range layers {
		digest, err := layer.Digest()
		if err != nil {
			log.Err(err).Msg("error while calling layer.Digest()")
			return err
		}
		rc, err := layer.Compressed()
		if err != nil {
			log.Err(err).Msg("error while calling layer.Compressed()")
			return err
		}
		defer rc.Close()

		filename := fmt.Sprintf(
			"layer-%02d-%s.tar.gz",
			i,
			strings.TrimPrefix(digest.String(), "sha256:"),
		)

		destPath := filepath.Join(outDir, filename)
		f, err := os.Create(destPath)
		if err != nil {
			log.Err(err).Msg("error while calling os.Create()")
			return err
		}

		if _, err := io.Copy(f, rc); err != nil {
			f.Close()
			log.Err(err).Msg("error while calling io.Copy()")
			return err
		}
		f.Close()
		log.Info().Msgf("saved %s\n", destPath)
	}
	return nil
}
