package warden

type WardenApplication struct {
	ID          string   `json:"id" yaml:"id"`
	Workdir     string   `json:"workdir" yaml:"workdir"`
	Command     []string `json:"command" yaml:"command"`
	Environment []string `json:"environment" yaml:"environment"`
}
