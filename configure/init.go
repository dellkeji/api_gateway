package configure

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Specification : detail configure
type Specification struct {
	Debug       bool             `yaml:debug`
	Version     string           `yaml:"version"`
	Location    string           `yaml:"location"`
	SvrConf     *ServerConfigure `yaml:"server_configure"`
	DBConf      *DBConfigure     `yaml:"database_configure"`
	ChannelSize int              `yaml:"channel_size"`
	// GraceTimeOut         int64         `envconfig:"GRACE_TIMEOUT"`
	// MaxIdleConnsPerHost  int           `envconfig:"MAX_IDLE_CONNS_PER_HOST"`
	// BackendFlushInterval time.Duration `envconfig:"BACKEND_FLUSH_INTERVAL"`
	// IdleConnTimeout      time.Duration `envconfig:"IDLE_CONN_TIMEOUT"`
	// RequestID            bool          `envconfig:"REQUEST_ID_ENABLED"`
	// Log                  logging.LogConfig
	// Web                  Web
	// Database             Database
	// Stats                Stats
	// Tracing              Tracing
	// TLS                  TLS
	// Cluster              Cluster
	// RespondingTimeouts   RespondingTimeouts
}

// GlobalConfigurations :
var GlobalConfigurations = &Specification{}

// InitConfig : load configure
func InitConfig(configFile string) error {
	GlobalConfigurations.Location = "Local"
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &GlobalConfigurations)
	if err != nil {
		return err
	}
	return nil
}
