package model

// LogEntry that the RaftServer will persist.
type LogEntry struct {
	ID    int64  `json:"id" gorm:"primaryKey"`
	Term  uint64 `json:"term" gorm:"unique"`
	Value int    `json:"value"`
}
