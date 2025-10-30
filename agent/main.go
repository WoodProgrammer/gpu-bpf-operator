package main

import (
	"bufio"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	LIB_PATH := os.Getenv("LIB_PATH")
	if len(LIB_PATH) == 0 {
		err := errors.New("Missing environment variable")
		log.Fatal().Err(err).Msg("Please environment variable LIB_PATH")
	}
	if err := generateBpftraceScript(LIB_PATH); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate bpftrace script")
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

// generateBpftraceScript generates a bpftrace script from a template
func generateBpftraceScript(libPath string) error {
	var probes []Probe
	templateData := TemplateProbeLib{
		ProbeLib: []string{},
		LibPath:  libPath,
	}

	log.Info().Msg("Generating bpftrace script from template...")
	probeEnv := os.Getenv("PROBE_CALLS")

	if len(probeEnv) == 0 {
		err := errors.New("missing environment variable PROBE_CALLS")
		log.Err(err).Msg("Please set PROBE_CALLS environment variable")
		return err
	}
	sDec, _ := b64.StdEncoding.DecodeString(probeEnv)
	if err := json.Unmarshal([]byte(sDec), &probes); err != nil {
		log.Err(err).Msg("ERror ")
	}
	// Access the parsed data
	for _, p := range probes {
		templateData.ProbeLib = append(templateData.ProbeLib, p.Name)
	}

	// Create template function map
	funcMap := template.FuncMap{
		"contains": func(needle string, haystack []string) bool {
			for _, item := range haystack {
				if strings.EqualFold(item, needle) {
					return true
				}
			}
			return false
		},
	}

	// Parse template
	tmpl, err := template.New("cuda_events.bt.tmpl").Funcs(funcMap).ParseFiles(TEMPLATE_FILE_PATH)
	if err != nil {
		return err
	}

	// Create output file
	f, err := os.Create(BT_FILE_PATH)
	if err != nil {
		return err
	}
	defer f.Close()

	// Execute template and write to file
	if err := tmpl.Execute(f, templateData); err != nil {
		return err
	}

	log.Info().Str("path", BT_FILE_PATH).Msg("CUDA Event tracer successfully generated")
	return nil
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
