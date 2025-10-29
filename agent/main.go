package main

import (
	"html/template"
	"os"
	"strings"
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
	var tmplFile = "templates/cuda_events.bt.tmpl"

	// Define custom template functions
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

	tmpl, err := template.New("cuda_events.bt.tmpl").Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, templateProbeLib)
	if err != nil {
		panic(err)
	}
}
