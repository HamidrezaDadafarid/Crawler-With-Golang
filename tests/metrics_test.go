package tests

import (
	"fmt"
	"main/database"
	"main/models"
	"main/repository"
	"testing"
)

func TestMetricsRepo(t *testing.T) {
	dbManager := database.GetInstnace()
	repo := repository.NewGormUMetric(dbManager.Db)

	addedMetiric, err := repo.Add(models.Metrics{
		TimeSpent:           10,
		CpuUsage:            10,
		RamUsage:            10,
		RequestCount:        100,
		SucceedRequestCount: 50,
		FailRequestCount:    51,
	})
	if err != nil {
		fmt.Println(err)
		t.Error("failed to add metric")
	}

	a, err := repo.GetTopTen()
	fmt.Println(a[0])
	if err != nil {
		t.Error("Failed to get metric")
	}

	err = repo.Delete(addedMetiric.ID)

	if err != nil {
		t.Error("Failed to delete user")
	}
}
