
package telegram

import (
	"fmt"
	"main/database"
	"main/models"
	"strconv"
	"strings"

	"gopkg.in/telebot.v4"
)

func (tg *Telegram) handleAddAd(ctx telebot.Context) error {
	adDetails := strings.Split(ctx.Text(), "|")
	if len(adDetails) < 5 {
		ctx.Send("Invalid format. Use: /add_ad title|description|price|city|category")
		return nil
	}

	sellPrice, err := strconv.Atoi(adDetails[2])
	if err != nil {
		ctx.Send("Invalid price format.")
		return err
	}

	ad := models.Ads{
		Description: adDetails[1],
		SellPrice:   uint(sellPrice),
		City:        adDetails[3],
		Neighborhood: adDetails[4],
	}

	if err := database.DB.Create(&ad).Error; err != nil {
		ctx.Send("Failed to add ad.")
		return err
	}

	ctx.Send(fmt.Sprintf("Ad added successfully. Ad ID: %d", ad.ID))
	return nil
}

func (tg *Telegram) handleDeleteAd(ctx telebot.Context) error {
	adID, err := strconv.Atoi(ctx.Text())
	if err != nil {
		ctx.Send("Invalid ad ID format.")
		return err
	}

	if err := database.DB.Delete(&models.Ads{}, adID).Error; err != nil {
		ctx.Send("Failed to delete ad.")
		return err
	}

	ctx.Send("Ad deleted successfully.")
	return nil
}

func (tg *Telegram) handleEditAd(ctx telebot.Context) error {
	adDetails := strings.Split(ctx.Text(), "|")
	if len(adDetails) < 6 {
		ctx.Send("Invalid format. Use: /edit_ad AdID|title|description|price|city|category")
		return nil
	}

	adID, err := strconv.Atoi(adDetails[0])
	if err != nil {
		ctx.Send("Invalid AdID format.")
		return err
	}

	price, err := strconv.Atoi(adDetails[3])
	if err != nil {
		ctx.Send("Invalid price format.")
		return err
	}

	var ad models.Ads
	if err := database.DB.First(&ad, adID).Error; err != nil {
		ctx.Send("Ad not found.")
		return err
	}

	ad.Title = adDetails[1]
	ad.Description = adDetails[2]
	ad.SellPrice = uint(price)
	ad.City = adDetails[4]
	ad.Neighborhood = adDetails[5]

	if err := database.DB.Save(&ad).Error; err != nil {
		ctx.Send("Failed to update ad.")
		return err
	}

	ctx.Send(fmt.Sprintf("Ad updated successfully. Ad ID: %d", ad.ID))
	return nil
}

func (tg *Telegram) handleGetAds(ctx telebot.Context) error {
	adIDText := ctx.Text()
	var ads []models.Ads

	if adIDText != "" {
		adID, err := strconv.Atoi(adIDText)
		if err != nil {
			ctx.Send("Invalid AdID format.")
			return err
		}
		if err := database.DB.Where("id = ?", adID).Find(&ads).Error; err != nil {
			ctx.Send("Ad not found.")
			return err
		}
	} else {
		if err := database.DB.Find(&ads).Error; err != nil {
			ctx.Send("Failed to retrieve ads.")
			return err
		}
	}

	for _, ad := range ads {
		ctx.Send(fmt.Sprintf("Ad ID: %d
Title: %s
Description: %s
Price: %d
City: %s
Category: %s",
			ad.ID, ad.Title, ad.Description, ad.SellPrice, ad.City, ad.Neighborhood))
	}
	return nil
}
