package config

type DBItemConf struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Database     string `yaml:"database"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Charset      string `yaml:"charset"`
	TimeOut      int    `yaml:"timeout"`
	WriteTimeOut int    `yaml:"write_time_out"`
	ReadTimeOut  int    `yaml:"read_time_out"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type DBLog struct {
	Enable bool   `yaml:"enable"`
	Level  string `yaml:"level"`
	Path   string `yaml:"path"`
	Type   string `yaml:"type"`
	Format string `yaml:"format"`
}

type dbItem struct {
	Log   DBLog      `yaml:"log"`
	Write DBItemConf `yaml:"write"`
	Read  DBItemConf `yaml:"read"`
}
