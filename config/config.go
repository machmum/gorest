package config

import (
	"fmt"
	"io/ioutil"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
)

// Load returns Configuration struct
func Load(path string) (*Configuration, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}
	var cfg = new(Configuration)
	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return cfg, nil
}

// Configuration holds data necessery for configuring application
type Configuration struct {
	Debug  bool      `yaml:"debug,omitempty"`
	Server *Server   `yaml:"server,omitempty"`
	DB     *Database `yaml:"database,omitempty"`
	Redis  *RDB      `yaml:"redis,omitempty"`
	// Conn   Connection
}

type Connection struct {
	Mysql *gorm.DB      `json:"mysql,omitempty"`
	Redis *redis.Client `json:"redis,omitempty"`
	Psql  *gorm.DB      `json:"psql,omitempty"`
}

// Database holds data necessery for database configuration
type Database struct {
	PSN        string `yaml:"psn,omitempty"`
	LogQueries bool   `yaml:"log_queries,omitempty"`
	Timeout    int    `yaml:"timeout_seconds,omitempty"`
	Mysql      string `yaml:"mysql,omitempty"`
}

// Server holds data necessery for server configuration
type Server struct {
	Port         string `yaml:"port,omitempty"`
	Debug        bool   `yaml:"debug,omitempty"`
	ReadTimeout  int    `yaml:"read_timeout_seconds,omitempty"`
	WriteTimeout int    `yaml:"write_timeout_seconds,omitempty"`
}

type RDB struct {
	Host     string      `yaml:"host"`
	Port     string      `yaml:"port,omitempty"`
	Password string      `yaml:"password,omitempty"`
	Lifetime RDBLifetime `yaml:"lifetime"`
	Prefix   RDBPrefix   `yaml:"prefix"`
}

type RDBLifetime struct {
	Token uint `yaml:"token"`
	Apps  uint `yaml:"apps"`
}

type RDBPrefix struct {
	Access  string `yaml:"access"`
	Refresh string `yaml:"refresh"`
	Apps    string `yaml:"apps"`
}
