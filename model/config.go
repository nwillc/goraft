package model

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Term    uint     `json:"term"`
	Members map[string]Member `json:"members"`
}

func ReadConfig(filename string) (Config, error) {
	configFile, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer configFile.Close()
	bytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
