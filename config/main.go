package config

import (
	"encoding/json"
	"io/ioutil"
)

var (
	// Conf is the loaded Config
	Conf Config
)

// Config holds the Config File
// The maps keys are the URLs to check and the values are int slices that define the status codes,
// for whom the health check will succeed
type Config map[string][]int

// LoadConfigFromFile returns a Config type initialized from the json file
func LoadConfigFromFile(path string) (Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
			return nil, err
	}
	var c Config
	if err = json.Unmarshal(file, &c); err != nil {
			return nil, err
	}
	return c, nil
}