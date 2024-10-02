package warden

type Application struct {
	ID          string   `json:"id" toml:"id"`
	Workdir     string   `json:"workdir" toml:"workdir"`
	Command     []string `json:"command" toml:"command"`
	Environment []string `json:"environment" toml:"environment"`
}

type ApplicationInfo struct {
	Application
	Status AppStatus `json:"status"`
}
