package model

// Status of the RaftServer to be persisted.
type Status struct {
	Name string `json:"name" gorm:"primaryKey"`
	Term uint64 `json:"term"`
}
