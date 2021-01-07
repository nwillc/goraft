package model

// LogEntry that the RaftServer will persist.
type LogEntry struct {
	Term  uint64 `json:"term" gorm:"primaryKey;autoIncrement:false"`
	Value int    `json:"value"`
}
