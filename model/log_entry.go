package model

type LogEntry struct {
	Term  uint `json:"term" gorm:"primaryKey"`
	Value uint `json:"value"`
}
