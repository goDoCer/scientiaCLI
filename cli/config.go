package cli

import (
	"encoding/json"
	"io/ioutil"
)

// Config is used to configure App
type Config struct {
	SaveDir string `json:"saveDir"`
}

// ReadConfig reads the filepath and returns a Config
func ReadConfig(filepath string) (Config, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// SaveConfig writes the Config to the filepath
func (c Config) SaveConfig(filepath string) error {
	file, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
