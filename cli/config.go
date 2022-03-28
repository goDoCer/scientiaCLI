package cli

import (
	"encoding/json"
	"io/ioutil"
)

// config is used to configure App
type config struct {
	SaveDir string `json:"saveDir"`
}

// ReadConfig reads the filepath and returns a Config
func readConfig(filepath string) (config, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return config{}, err
	}

	var cfg config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return config{}, err
	}

	return cfg, nil
}

// SaveConfig writes the Config to the filepath
func (cfg config) saveConfig(filepath string) error {
	file, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath, file, 0777)
	if err != nil {
		return err
	}

	return nil
}
