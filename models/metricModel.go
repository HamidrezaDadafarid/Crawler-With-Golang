package models

import "gorm.io/gorm"

type Metrics struct {
	gorm.Model
	ID                  uint `gorm:"primaryKey;autoIncrement"`
	TimeSpent           uint //second
	CpuUsage            uint
	RamUsage            uint
	RequestCount        uint
	succeedRequestCount uint
	FailRequestCount    uint
}
