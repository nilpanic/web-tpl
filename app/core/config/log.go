package config

type Log struct {
	Level     string `yaml:"level"`
	Output    string `yaml:"output"`
	LogFormat string `yaml:"log_format"`
	Name      string `yaml:"name"`
}
