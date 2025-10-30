package main

const (
	TEMPLATE_FILE_PATH = "templates/cuda_events.bt.tmpl"
	BT_FILE_PATH       = "/tmp/cuda_events.bt"
)

type TemplateProbeLib struct {
	ProbeLib []string
	LibPath  string
}

type Probe struct {
	Kind string `json:"Kind"`
	Name string `json:"Name"`
}
