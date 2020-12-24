package model

type RaftLogEntry struct {
	ID    uint `json:"id", gorm:"primaryKey"`
	Value uint `json:"value"`
}
