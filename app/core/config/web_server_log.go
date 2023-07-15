package config

type WebServerLog struct {
	Enable          bool     `yaml:"enable"`
	LogIDShowHeader bool     `yaml:"log_id_show_header"`
	LogPath         string   `yaml:"log_path"`
	LogFormat       string   `yaml:"log_format"`
	SkipPaths       []string `yaml:"skip_paths"`
	Output          string   `yaml:"output"`
}
