package model

// LogEntry that the RaftServer will persist.
type LogEntry struct {
	EntryNo int64  `json:"entryNo" gorm:"primaryKey;autoIncrement:false"`
	Term    uint64 `json:"term"`
	Value   int64  `json:"value"`
}
