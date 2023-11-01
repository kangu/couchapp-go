package tests

import (
	"encoding/json"
	"fmt"
	"os"
)

type TestConfig struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

func LoadTestConfig() (TestConfig, error) {
	var config TestConfig
	// tests are always executed from one folder up
	// so the full /tests path needs to be provided
	configFile, err := os.Open("tests/config.json")
	if err != nil {
		fmt.Printf("Cannot load test config file, make sure it's setup: %s\n", err)
		os.Exit(1)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
