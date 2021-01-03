package model

type LogEntry struct {
	Term  uint64 `json:"term" gorm:"primaryKey"`
	Value int    `json:"value"`
}
