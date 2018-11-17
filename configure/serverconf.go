package configure

// ServerConfigure : spec for host and port
type ServerConfigure struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
