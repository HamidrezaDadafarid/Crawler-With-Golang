package tests

import (
	"fmt"
	"main/database"
	"main/models"
	"main/repository"
	"testing"
)

func TestUsersRepo(t *testing.T) {
	dbManager := database.GetInstnace()
	repo := repository.NewGormUser(dbManager.Db)

	addedUser, err := repo.Add(models.Users{
		TelegramId: "telegram id",
		Role:       "User",
	})
	if err != nil {
		fmt.Println(err)
		t.Error("failed to add user")
	}

	addedUser.Role = "Admin"
	err = repo.Update(addedUser)
	if err != nil {
		t.Error("failed to update user")
	}

	updatedUser, err := repo.Get([]uint{addedUser.ID})

	if err != nil {
		t.Error("Failed to get user")
	}

	if updatedUser[0].Role != "Admin" {
		t.Error("failed to update user")
	}

	err = repo.Delete(addedUser.ID)

	if err != nil {
		t.Error("failed to delete user")
	}
}
