package model

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type Config struct {
	HeartbeatTimeout int      `json:"heartbeat_timeout"`
	ElectionTimeout  int      `json:"election_timeout"`
	MinOffset        int      `json:"min_offset"`
	MaxOffset        int      `json:"max_offset"`
	Members          []Member `json:"members"`
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

func (c *Config) ElectionCountdown() time.Duration {
	return (time.Duration(c.ElectionTimeout) * time.Millisecond) * time.Duration(rand.Intn(c.MaxOffset-c.MinOffset)+c.MinOffset)
}

func (c *Config) HeartbeatCountdown() time.Duration {
	return time.Duration(c.HeartbeatTimeout) * time.Millisecond
}
