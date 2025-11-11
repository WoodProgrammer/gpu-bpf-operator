package imagelayer

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/rs/zerolog/log"
)

type ImageLayers interface {
	FetchImageLayers(imageRef, outDir string) error
	FlattenImage(ctx context.Context, imageRef, outDir string) error
}

type ImageLayerHandler struct {
	ImageLayer *ImageLayers
}

func (i *ImageLayerHandler) FetchImageLayers(imageRef, outDir string) error {
	outDir = fmt.Sprintf("/tmp/%s.img.zip", imageRef)
	if err := i.FlattenImage(imageRef, outDir); err != nil {
		log.Err(err).Msg("error while calling FlattenImage()")
		return err
	}
	log.Info().Msg("Image Layers properly downloaded")
	return nil
}

func (i *ImageLayerHandler) FlattenImage(imageRef, output string) error {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		log.Err(err).Msg("error while calling name.ParseReference()")
		return err
	}
	img, err := remote.Image(ref)
	if err != nil {
		log.Err(err).Msg("error while calling remote.Image()")
		return err
	}

	out, err := os.Create(output)
	if err != nil {
		log.Err(err).Msg("error while calling os.Create()")
		return err
	}
	defer out.Close()

	gz := gzip.NewWriter(out)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()

	layers, err := img.Layers()
	if err != nil {
		log.Err(err).Msg("error while calling img.Layers()")
		return err
	}

	for _, layer := range layers {
		rc, err := layer.Uncompressed()
		if err != nil {
			log.Err(err).Msg("error while calling layer.Uncompressed()")
			return err
		}
		defer rc.Close()

		tr := tar.NewReader(rc)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Err(err).Msg("error during the iteration on layers")
				return err
			}
			if err := tw.WriteHeader(hdr); err != nil {
				log.Err(err).Msg("error while calling tw.WriteHeader()")
				return err
			}
			if _, err := io.Copy(tw, tr); err != nil {
				log.Err(err).Msg("error while calling io.Copy()")
				return err
			}
		}
	}
	return nil
}
