package model

type Status struct {
	Name string `json:"name" gorm:"primaryKey"`
	Term uint64 `json:"term" gorm:"primaryKey;autoIncrement:false"`
}
