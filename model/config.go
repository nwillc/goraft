package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Member struct {
	Port uint   `json:"port"`
}

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

func (m *Member) Address() string {
	return fmt.Sprintf(":%d", m.Port)
}
