package main

import (
	"html/template"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	TEMPLATE_FILE_PATH = "templates/cuda_events.bt.tmpl"
	BT_FILE_PATH       = "/tmp/cuda_events.bt"
)

type TemplateProbeLib struct {
	ProbeLib []string
	LibPath  string
}

func main() {
	templateProbeLib := TemplateProbeLib{}
	probeLib := []string{
		"cudaEventRecord",
		"cudaEventSynchronize",
	}

	templateProbeLib.ProbeLib = probeLib
	templateProbeLib.LibPath = "/lib/cudaPath"

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

	tmpl, err := template.New("cuda_events.bt.tmpl").Funcs(funcMap).ParseFiles(TEMPLATE_FILE_PATH)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(BT_FILE_PATH)
	if err != nil {
		log.Error().Err(err).Msg("Error while running os.Create()")
		panic(err)
	}
	err = tmpl.Execute(f, templateProbeLib)
	if err != nil {
		panic(err)
	}
	log.Info().Msg("CUDA Event tracer successfully generated.")
}
