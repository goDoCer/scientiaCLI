package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/goDoCer/scientiaCLI/scientia"
)

// Config is used to configure App
type Config struct {
	SaveDir      string `json:"saveDir"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// LoadConfig reads the filepath and returns a Config
func LoadConfig(filepath string) (Config, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Save writes the Config to the filepath
func (cfg Config) Save(filepath string) error {
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

func (cfg *Config) Tokens() scientia.LoginTokens {
	return scientia.LoginTokens{
		AccessToken:  cfg.AccessToken,
		RefreshToken: cfg.RefreshToken,
	}
}

func (cfg *Config) UpdateTokens(tokens scientia.LoginTokens) {
	cfg.AccessToken = tokens.AccessToken
	cfg.RefreshToken = tokens.RefreshToken
}
