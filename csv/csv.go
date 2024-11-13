package csv

import (
	"encoding/csv"
	"fmt"
	"log"
	"main/repository"
	"os"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ExportCsv(telID string, db *gorm.DB) (string, error) {
	u := repository.NewGormUser(db)

	// Get user by telegramID
	user, err := u.GetByTelegramId(telID)
	if err != nil {
		log.Println("Couldn't find user: ", err)
		return "", err
	}

	fileName := fmt.Sprintf("%s.csv", uuid.New()) // Unique file name
	// Create a CSV file
	f, err := os.Create("../telegram/" + fileName)

	if err != nil {
		log.Println("Couldn't create csv file: ", err)
		return "", err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Write headers
	header := []string{"لینک", "طول جغرافیایی", "عرض جغرافیایی", "توضیحات", "قیمت فروش", "قیمت اجاره", "مقدار رهن", "شهر", "محله", "متراژ", "تعداد اتاق", "دسته بندی", "سن بنا", "دسته بندی آپارتمان/ویلا", "طبقه", "انباری", "آسانسور", "عنوان آگهی", "لینک عکس"}
	if err := w.Write(header); err != nil {
		log.Println("error writing header to csv:", err)
		return "", err
	}

	for _, v := range user.Ads {
		val := reflect.ValueOf(v).Elem() // Get the value of the pointer
		record := make([]string, 0)

		for i := 0; i < (val.NumField() - 1); i++ {
			field := val.Field(i)

			//skip: CreatedAt, UpdatedAt, DeletedAt,  ID, UniqueId, NumberOfViews, Users(last field)
			if i != 0 && i != 1 && i != 2 && i != 3 && i != 5 && i != 9 {
				var fieldValue string
				// Check for CategoryPR to insert specific string values
				if i == 17 { // Assuming CategoryPR is at index 17
					if field.Interface().(uint) == 0 {
						fieldValue = "فروشی"
					} else if field.Interface().(uint) == 1 {
						fieldValue = "اجاره"
					} else {
						fieldValue = fmt.Sprintf("%v", field.Interface())
					}
				} else if i == 19 { // Assuming CategoryAV is at index 19
					if field.Interface().(uint) == 0 {
						fieldValue = "آپارتمان"
					} else if field.Interface().(uint) == 1 {
						fieldValue = "ویلایی"
					} else {
						fieldValue = fmt.Sprintf("%v", field.Interface())
					}
				} else if i == 21 || i == 22 { // Check for boolean fields Storage and Elevator
					if field.Interface().(bool) {
						fieldValue = "دارد"
					} else {
						fieldValue = "ندارد"
					}
				} else {
					// Convert other fields to string
					fieldValue = fmt.Sprintf("%v", field.Interface())
				}
				// Append the processed field value to record
				record = append(record, fieldValue)
			}
		}

		if err := w.Write(record); err != nil {
			log.Println("error writing record to csv:", err)
			return "", err
		}
	}
	return fileName, nil
}
