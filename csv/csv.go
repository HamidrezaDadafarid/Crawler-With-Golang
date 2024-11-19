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
	gormUserAd := repository.NewGormUser_Ad(db)
	gormAd := repository.NewGormAd(db)

	user, err := u.GetByTelegramId(telID)
	if err != nil {
		log.Println("Couldn't find user", err)
		return "", err
	}

	userAd, err := gormUserAd.GetByUserId([]uint{user.ID})
	if err != nil {
		log.Println("Couldn't find user ad", err)
		return "", err
	}

	adID := []uint{}

	for _, v := range userAd {
		adID = append(adID, v.AdId)
	}

	ads, err := gormAd.GetById(adID)
	if err != nil {
		log.Println("Couldn't find user: ", err)
		return "", err
	}

	fileName := fmt.Sprintf("./telegram/%s.csv", uuid.New()) // Unique file name
	// Create a CSV file
	f, err := os.Create(fileName)

	if err != nil {
		log.Println("Couldn't create csv file: ", err)
		return "", err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Write headers
	header := []string{"آیدی", "سایت", "لینک", "طول جغرافیایی", "عرض جغرافیایی", "توضیحات", "قیمت فروش", "قیمت اجاره", "مقدار رهن", "شهر", "محله", "متراژ", "تعداد اتاق", "دسته بندی خرید/اجاره", "سن بنا", "دسته بندی آپارتمان/ویلا", "طبقه", "انباری", "آسانسور", "پارکینگ", "عنوان آگهی", "لینک عکس"}
	if err := w.Write(header); err != nil {
		log.Println("error writing header to csv:", err)
		return "", err
	}

	for _, ad := range ads {
		val := reflect.ValueOf(ad)
		record := make([]string, 0)
		var fieldValue string
		for i := 0; i < val.NumField()-1; i++ {
			if i == 10 || i == 0 { // skip numberOfviews and created_at
				continue
			}

			field := val.Field(i)
			fieldValue = fmt.Sprintf("%v", field.Interface())
			// Append the processed field value to record
			record = append(record, fieldValue)
		}
		if err := w.Write(record); err != nil {
			log.Println("error writing record to csv:", err)
			return "", err
		}
	}
	return fileName, nil
}
