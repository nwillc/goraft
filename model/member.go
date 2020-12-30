package model

import (
	"fmt"
)

type Member struct {
	Name string   `json:"name"`
	Port uint32   `json:"port"`
}

var _ fmt.Stringer = (* Member)(nil)

func (m *Member) String() string {
	return fmt.Sprintf("{ name: %s, port: %d }", m.Name, m.Port)
}

func (m *Member) Address() string {
	return fmt.Sprintf(":%d", m.Port)
}

