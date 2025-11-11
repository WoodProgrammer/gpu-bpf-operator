package main

const (
	TEMPLATE_FILE_PATH       = "templates/nvidia_events.bt.tmpl"
	BT_FILE_PATH             = "/tmp/nvidia_events.bt"
	ENVIRONMENT_VARIABLE_ERR = "Please fille the environment for "
	SCAN_BUFF                = 1024 * 1024
)

type TemplateProbeLib struct {
	ProbeLib []string
	LibPath  string
}

type Probe struct {
	Kind string `json:"Kind"`
	Name string `json:"Name"`
}
