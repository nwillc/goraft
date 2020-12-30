package model

type RaftLogEntry struct {
	Term  uint `json:"term" gorm:"primaryKey"`
	Value uint `json:"value"`
}
