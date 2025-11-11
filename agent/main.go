package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"os/signal"
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
		err := errors.New("IMAGE_REF or OUTPUT_DIRECTORY is empty")
		log.Fatal().Err(err).Msg("Please set your IMAGE_REF and OUTPUT_DIRECTORY which contains image registry address and directory to keep the BPF scripts")
	}

	err := imageLayerHandler.FetchImageLayers(imageRef, outDir)
	if err != nil {
		log.Fatal().Err(err).Msg("error while calling imageLayerHandler.FetchImageLayers()")
	}
	// Step 2: Set up context and signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := setupSignalHandler()

	// Step 3: Execute the bpftrace script
	if err := executeBpftraceScript(ctx, sigChan, cancel); err != nil {
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
func executeBpftraceScript(ctx context.Context, sigChan chan os.Signal, cancel context.CancelFunc) error {
	log.Info().Msg("Starting bpftrace script execution...")
	cmd := exec.CommandContext(ctx, "/usr/bin/bpftrace", BT_FILE_PATH)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	log.Info().Msg("bpftrace script started successfully")

	// Stream outputs concurrently
	go streamOutput(stdoutPipe, "stdout")
	go streamOutput(stderrPipe, "stderr")

	// Wait for interrupt signal
	<-sigChan
	log.Info().Msg("Received interrupt signal, shutting down bpftrace...")
	cancel()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		// Ignore "signal: killed" error as it's expected
		if err.Error() != "signal: killed" {
			log.Error().Err(err).Msg("Error while waiting for bpftrace to exit")
			return err
		}
	}

	log.Info().Msg("bpftrace script stopped successfully")
	return nil
}

// streamOutput reads from a pipe line-by-line and logs with source tag
func streamOutput(pipe io.Reader, source string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		log.Info().Str("source", source).Msg(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Str("source", source).Msg("Error reading output")
	}
}
