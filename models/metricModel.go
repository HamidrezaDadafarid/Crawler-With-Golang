package models

import "gorm.io/gorm"

type Metrics struct {
	gorm.Model
	TimeSpent           float64 //second
	CpuUsage            float64
	RamUsage            float64
	RequestCount        int
	SucceedRequestCount uint
	FailRequestCount    uint
}
