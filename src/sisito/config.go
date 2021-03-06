package sisito

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	User     []UserConfig
	Filter   []FilterConfig
	Authz    AuthzConfig
}

type ServerConfig struct {
	Log    string
	Gzip   bool
	Prefix string
}

type DatabaseConfig struct {
	Host     string
	Port     int64
	Database string
	Username string
	Password string
	Timezone string
}

type UserConfig struct {
	Userid   string
	Password string
}

type FilterConfig struct {
	Key      string
	Operator string
	Value    string
	Values   []string
	Join     string
	Sql      string
}

type AuthzConfig struct {
	Recent    bool
	Listed    bool
	Blacklist bool
}

func LoadConfig(flags *Flags) (config *Config, err error) {
	config = &Config{}
	_, err = toml.DecodeFile(flags.Config, config)

	if err != nil {
		return
	}

	database := config.Database

	if database.Host == "" {
		database.Host = "localhost"
	}

	if database.Port == 0 {
		database.Port = 3306
	}

	if database.Username == "" {
		database.Username = "root"
	}

	for i := 0; i < len(config.Filter); i++ {
		filter := &config.Filter[i]

		if filter.Sql == "" && filter.Operator == "" {
			if len(filter.Values) > 0 {
				filter.Operator = "IN"
			} else {
				filter.Operator = "="
			}
		}

		if filter.Join == "" {
			filter.Join = "AND"
		}
	}

	return
}
