package models

import "gorm.io/gorm"

type Metrics struct {
	gorm.Model
	ID                  uint    `gorm:"primaryKey;autoIncrement"`
	TimeSpent           float64 //second
	CpuUsage            float64
	RamUsage            float64
	RequestCount        int
	SucceedRequestCount uint
	FailRequestCount    uint
}
