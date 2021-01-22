package model

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

// Config is the configuration of a Raft cluster.
type Config struct {
	HeartbeatTimeout int      `json:"heartbeat_timeout"`
	ElectionTimeout  int      `json:"election_timeout"`
	MinOffset        int      `json:"min_offset"`
	MaxOffset        int      `json:"max_offset"`
	Members          []Member `json:"members"`
}

// ReadConfig reads a Config stored as JSON in filename.
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

// ElectionCountdown calculate the election countdown based on the Config.
func (c *Config) ElectionCountdown() time.Duration {
	return (time.Duration(c.ElectionTimeout) * time.Millisecond) * time.Duration(rand.Intn(c.MaxOffset-c.MinOffset)+c.MinOffset)
}

// HeartbeatCountDown calculates the heartbeat countdown based on the Config.
func (c *Config) HeartbeatCountDown() time.Duration {
	return time.Duration(c.HeartbeatTimeout) * time.Millisecond
}

// Peers returns the slice of Member's that are peers to memberName in the configs Members
func (c *Config) Peers(memberName string) []Member {
	var peers []Member
	for _, peer := range c.Members {
		if peer.Name == memberName {
			continue
		}
		peers = append(peers, peer)
	}
	return peers
}
