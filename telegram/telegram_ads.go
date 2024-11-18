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

	ctx.Send(fmt.Sprintf("Ad added successfully. Ad ID: %d", ad.AdID))
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

func (tg *Telegram) handleListAds(ctx telebot.Context) error {
	var ads []models.Ads
	if err := database.DB.Find(&ads).Error; err != nil {
		ctx.Send("Failed to fetch ads.")
		return err
	}

	if len(ads) == 0 {
		ctx.Send("No ads found.")
		return nil
	}

	response := "List of Ads:\n"
	for _, ad := range ads {
		response += fmt.Sprintf("ID: %d, Description: %s, Price: %d, City: %s\n", ad.AdID, ad.Description, ad.SellPrice, ad.City)
	}
	ctx.Send(response)
	return nil
}

func (tg *Telegram) handleUpdateAd(ctx telebot.Context) error {
	adDetails := strings.Split(ctx.Text(), "|")
	if len(adDetails) < 6 {
		ctx.Send("Invalid format. Use: /update_ad ad_id|description|price|city|neighborhood|category")
		return nil
	}

	adID, err := strconv.Atoi(adDetails[0])
	if err != nil {
		ctx.Send("Invalid ad ID format.")
		return err
	}

	sellPrice, err := strconv.Atoi(adDetails[2])
	if err != nil {
		ctx.Send("Invalid price format.")
		return err
	}

	var ad models.Ads
	if err := database.DB.First(&ad, adID).Error; err != nil {
		ctx.Send("Ad not found.")
		return err
	}

	ad.Description = adDetails[1]
	ad.SellPrice = uint(sellPrice)
	ad.City = adDetails[3]
	ad.Neighborhood = adDetails[4]

	if err := database.DB.Save(&ad).Error; err != nil {
		ctx.Send("Failed to update ad.")
		return err
	}

	ctx.Send("Ad updated successfully.")
	return nil
}
