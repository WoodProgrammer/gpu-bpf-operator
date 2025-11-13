package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	imageLayer "github.com/WoodProgrammer/generic-gpu-operator/imagelayer"
	"github.com/rs/zerolog/log"
)

func CreateNewImageLayerHandler() imageLayer.ImageLayerHandler {
	return imageLayer.ImageLayerHandler{}
}
func main() {
	imageLayerHandler := CreateNewImageLayerHandler()
	imageRef := os.Getenv("IMAGE_REF")

	outDir := os.Getenv("OUTPUT_DIRECTORY")
	if len(imageRef) == 0 || len(outDir) == 0 {
		err := errors.New(ENVIRONMENT_VARIABLE_ERR)
		log.Fatal().Err(err).Msg("OUTPUT_DIRECTORY IMAGE_REF")
	}

	err := imageLayerHandler.FetchImageLayers(imageRef, outDir)
	if err != nil {
		log.Fatal().Err(err).Msg("error while calling imageLayerHandler.FetchImageLayers()")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bpfFilePath := os.Getenv("BPF_FILE_PATH")
	if len(imageRef) == 0 || len(outDir) == 0 {
		err := errors.New(ENVIRONMENT_VARIABLE_ERR)
		log.Fatal().Err(err).Msgf("BPF_FILE_PATH")
	}
	bpfFilePath = imageLayerHandler.ScriptPathDir + bpfFilePath
	sigChan := setupSignalHandler()
	if err := executeBpftraceScript(ctx, sigChan, cancel, bpfFilePath); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute bpftrace script")
	}
	log.Info().Msg("Application shutdown complete")
}

// setupSignalHandler configures signal handling for graceful shutdown
func setupSignalHandler() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Info().Msg("Signal handler configured (SIGINT, SIGTERM)")
	return sigChan
}

// executeBpftraceScript executes the bpftrace script and streams output
func executeBpftraceScript(ctx context.Context, sigChan chan os.Signal, cancel context.CancelFunc, bpfFilePath string) error {
	var wg sync.WaitGroup
	stdArrMap := map[string]io.ReadCloser{
		"stdout": nil,
		"stderr": nil,
	}

	log.Info().Msg("Starting bpftrace script execution...")
	cmd := exec.CommandContext(ctx, "/usr/bin/bpftrace", bpfFilePath)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdArrMap["stdout"] = stdoutPipe

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdArrMap["stderr"] = stderrPipe
	if err := cmd.Start(); err != nil {
		return err
	}

	log.Info().Msg("bpftrace script started successfully")

	go func() {
		<-sigChan
		log.Info().Msg("Received interrupt signal, shutting down bpftrace...")
		cancel()
	}()
	defer cancel()
	for k, pipe := range stdArrMap {
		wg.Add(1)
		go func(source string, r io.Reader) {
			defer wg.Done()
			streamOutput(ctx, r, source)
		}(k, pipe)
	}
	wg.Wait()
	log.Info().Msg("bpftrace script stopped successfully")
	return nil
}

func streamOutput(ctx context.Context, pipe io.Reader, source string) {
	scanner := bufio.NewScanner(pipe)
	scanner.Buffer(make([]byte, 1024), SCAN_BUFF)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					log.Error().Err(err).Str("source", source).Msg("Error reading output")
				}
				return
			}
			log.Info().Str("source", source).Msg(scanner.Text())
		}
	}
}
