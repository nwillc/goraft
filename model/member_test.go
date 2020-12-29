package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMember(t *testing.T) {
	type fields struct {
		Name string
		Port uint32
	}
	tests := []struct {
		name     string
		fields   fields
		stringer string
		addresss string
	}{
		{
			name: "One1",
			fields: fields{
				Name: "One",
				Port: uint32(1),
			},
			stringer: "{ name: One, port: 1 }",
			addresss: ":1",
		},
		{
			name: "Two1",
			fields: fields{
				Name: "Two",
				Port: uint32(1),
			},
			stringer: "{ name: Two, port: 1 }",
			addresss: ":1",
		},
		{
			name: "Three4057",
			fields: fields{
				Name: "Three",
				Port: uint32(4057),
			},
			stringer: "{ name: Three, port: 4057 }",
			addresss: ":4057",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Member{
				Name: tt.fields.Name,
				Port: tt.fields.Port,
			}
			assert.Equal(t, tt.stringer, m.String())
			assert.Equal(t, tt.addresss, m.Address())
		})
	}
}
