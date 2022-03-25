package cli

import (
	"encoding/json"
	"os"
	"scientia-cli/scientia"
)

func saveDetails(tokens scientia.LoginTokens) error {
	b, err := json.Marshal(tokens)
	if err != nil {
		return err
	}
	return os.WriteFile(tokenPath, b, 0777)
}

func loadDetails() (*scientia.LoginTokens, error) {
	b, err := os.ReadFile(tokenPath)
	if err != nil {
		return nil, err
	}
	var tokens scientia.LoginTokens
	err = json.Unmarshal(b, &tokens)
	return &tokens, err
}
