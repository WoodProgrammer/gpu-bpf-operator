package main

const (
	TEMPLATE_FILE_PATH = "templates/nvidia_events.bt.tmpl"
	BT_FILE_PATH       = "/tmp/nvidia_events.bt"
)

type TemplateProbeLib struct {
	ProbeLib []string
	LibPath  string
}

type Probe struct {
	Kind string `json:"Kind"`
	Name string `json:"Name"`
}
