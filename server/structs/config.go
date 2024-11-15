package structs

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
)

// ServiceConfig is the configuration for a specific service
type ServiceConfig struct {
	FriendlyURL    string                 `json:"friendly_url"`
	Protocol       string                 `json:"protocol"`
	ProtocolConfig map[string]interface{} `json:"protocol_config,omitempty"`
}

// Config mirrors the config.json file which holds
// a dictionary of services with their ServiceConfigs
type Config struct {
	Host             string                     `json:"host,omitempty"`
	DBPath           string                     `json:"db_path,omitempty"`
	CheckInterval    int                        `json:"check_interval,omitempty"`
	NoFrontend       bool                       `json:"no_frontend,omitempty"`
	FrontendPath     string                     `json:"frontend_path,omitempty"`
	FrontendTitle    string                     `json:"frontend_title,omitempty"`
	FrontendLogoPath string                     `json:"frontend_logo_path,omitempty"`
	AllowedAPIKeys   []string                   `json:"allowed_api_keys,omitempty"`
	Notifiers        map[string]json.RawMessage `json:"notifiers,omitempty"`
	Services         map[string]ServiceConfig   `json:"services"`
}

const (
	hostDefault             = "0.0.0.0:3002"
	dbPathDefault           = "./db.sqlite"
	checkIntervalDefault    = 120
	frontendPathDefault     = "../frontend/build"
	frontendTitleDefault    = "Service Status"
	frontendLogoPathDefault = "logo.png"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (c *Config) setDefaults() error {
	if c.Host == "" {
		log.Printf("\"host\" not defined in config -> using default (%s)", hostDefault)
		c.Host = hostDefault
	}
	if c.DBPath == "" {
		log.Printf("\"db_path\" not defined in config -> using default (%s)", dbPathDefault)
		c.DBPath = dbPathDefault
	}
	if c.CheckInterval == 0 {
		log.Printf("\"check_interval\" not defined in config -> using default (%d)", checkIntervalDefault)
		c.CheckInterval = checkIntervalDefault
	}
	if c.FrontendPath == "" {
		log.Printf("\"frontend_path\" not defined in config -> using default (%s)", frontendPathDefault)
		c.FrontendPath = frontendPathDefault
	}
	if c.FrontendTitle == "" {
		log.Printf("\"frontend_title\" not defined in config -> using default (%s)", frontendTitleDefault)
		c.FrontendTitle = frontendTitleDefault
	}
	if c.FrontendLogoPath == "" {
		log.Printf("\"frontend_logo_path\" not defined in config -> using default (%s)", frontendLogoPathDefault)
		c.FrontendLogoPath = frontendLogoPathDefault
	}
	if len(c.AllowedAPIKeys) == 0 {
		log.Println("\"allowed_api_keys\" not defined in config -> genereting a temporary one")

		b := make([]rune, 30)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		newKey := string(b)
		c.AllowedAPIKeys = []string{newKey}

		log.Printf("Your temporary API key is '%s' (keep in mind that it will be regenerated after a restart if you do not generate one by yourself)", newKey)
	}

	return nil
}

// ForService returns the ServiceConfig for a specific
// Service by name
func (c *Config) ForService(name string) *ServiceConfig {
	if sc, ok := c.Services[name]; ok {
		return &sc
	}
	return nil
}

// ParseConfig unmarshalls the file provided to the
// Config struct and returns a pointer to it
func ParseConfig(filename string) (*Config, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	c := Config{}

	if err := json.Unmarshal(jsonBytes, &c); err != nil {
		return nil, err
	}

	if err := c.setDefaults(); err != nil {
		return nil, err
	}

	return &c, nil
}
