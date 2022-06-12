package cli

import (
	"encoding/json"
	"io/ioutil"

	"github.com/goDoCer/scientiaCLI/scientia"
)

// config is used to configure App
type config struct {
	SaveDir      string `json:"saveDir"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// loadConfig reads the filepath and returns a Config
func loadConfig(filepath string) (config, error) {
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

// save writes the Config to the filepath
func (cfg config) save(filepath string) error {
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

func (cfg *config) tokens() (scientia.LoginTokens, bool) {
	return scientia.LoginTokens{
		AccessToken:  cfg.AccessToken,
		RefreshToken: cfg.RefreshToken,
	}, cfg.AccessToken != ""
}

func (cfg *config) updateTokens(tokens scientia.LoginTokens) {
	cfg.AccessToken = tokens.AccessToken
	cfg.RefreshToken = tokens.RefreshToken
}
