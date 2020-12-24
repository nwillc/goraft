package model

type RaftLogEntry struct {
	Position uint `json:"position" gorm:"primaryKey"`
	Value    uint `json:"value"`
}
