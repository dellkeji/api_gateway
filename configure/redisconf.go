package configure

// RedisConfigure : spec for redis configure
type RedisConfigure struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Password  string `yaml:"password"`
	DefaultDB int    `yaml:"default_db"`

	MaxPoolSize    int    `yaml:"max_pool_size"`
	MaxConnTimeout int    `yaml:"max_conn_timeout"`
	IdleTimeout    int    `yaml:"idle_timeout"`
	ReadTimeout    int    `yaml:"read_timeout"`
	WriteTimeout   int    `yaml:"write_timeout"`
	QueueName      string `yaml:"queue_name"`
}
