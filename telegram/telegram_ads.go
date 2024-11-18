package main

import (
	"fmt"
	"main/database"
	"main/models"
	"main/utils"
	"strings"

	telebot "gopkg.in/tucnak/telebot.v2"
)

func (tg *Telegram) handleAddAd(ctx telebot.Context) error {
	input := ctx.Message().Payload
	if input == "" {
		ctx.Send("Please provide the ad information in the format: Title|Description|SellPrice|City|Neighborhood|Meters|NumberOfRooms|PhoneNumber.")
		return nil
	}

	parts := strings.Split(input, "|")
	if len(parts) < 8 {
		ctx.Send("The entered information format is incorrect. Please try again.")
		return nil
	}

	sellPrice, err := utils.ParseUint(parts[2])
	if err != nil {
		ctx.Send("Invalid sell price format. Please enter a numeric value.")
		return nil
	}

	meters, err := utils.ParseUint(parts[5])
	if err != nil {
		ctx.Send("Invalid meters value. Please enter a numeric value.")
		return nil
	}

	numberOfRooms, err := utils.ParseUint(parts[6])
	if err != nil {
		ctx.Send("Invalid number of rooms. Please enter a numeric value.")
		return nil
	}

	ad := models.Ads{
		Title:             parts[0],
		Description:       parts[1],
		SellPrice:         sellPrice,
		City:              parts[3],
		Neighborhood:      parts[4],
		Meters:            meters,
		NumberOfRooms:     numberOfRooms,
		SellerPhoneNumber: parts[7],
	}

	if err := models.AddAd(database.DB, &ad); err != nil {
		ctx.Send("Error adding the ad: " + err.Error())
		return nil
	}

	ctx.Send("The ad was added successfully!")
	return nil
}

func (tg *Telegram) handleEditAd(ctx telebot.Context) error {
	input := ctx.Message().Payload
	if input == "" {
		ctx.Send("Please provide the ad information in the format: AdID|NewTitle|NewDescription|NewSellPrice|City|Neighborhood|Meters|NumberOfRooms|PhoneNumber.")
		return nil
	}

	parts := strings.Split(input, "|")
	if len(parts) < 9 {
		ctx.Send("The entered information format is incorrect. Please try again.")
		return nil
	}

	adID := utils.ParseUint(parts[0])
	sellPrice, err := utils.ParseUint(parts[3])
	if err != nil {
		ctx.Send("Invalid sell price format. Please enter a numeric value.")
		return nil
	}

	meters, err := utils.ParseUint(parts[6])
	if err != nil {
		ctx.Send("Invalid meters value. Please enter a numeric value.")
		return nil
	}

	numberOfRooms, err := utils.ParseUint(parts[7])
	if err != nil {
		ctx.Send("Invalid number of rooms. Please enter a numeric value.")
		return nil
	}

	var ad models.Ads
	if err := database.DB.First(&ad, adID).Error; err != nil {
		ctx.Send("No ad found with this ID.")
		return nil
	}

	ad.Title = parts[1]
	ad.Description = parts[2]
	ad.SellPrice = sellPrice
	ad.City = parts[4]
	ad.Neighborhood = parts[5]
	ad.Meters = meters
	ad.NumberOfRooms = numberOfRooms
	ad.SellerPhoneNumber = parts[8]

	if err := models.EditAd(database.DB, &ad); err != nil {
		ctx.Send("Error editing the ad: " + err.Error())
		return nil
	}

	ctx.Send("The ad was successfully edited!")
	return nil
}

func (tg *Telegram) handleDeleteAd(ctx telebot.Context) error {
	adID := utils.ParseUint(ctx.Message().Payload)
	if adID == 0 {
		ctx.Send("Please enter the ad ID.")
		return nil
	}

	if err := models.DeleteAd(database.DB, adID); err != nil {
		ctx.Send("Error deleting the ad: " + err.Error())
		return nil
	}

	ctx.Send("The ad was successfully deleted!")
	return nil
}

func (tg *Telegram) handleViewAds(ctx telebot.Context) error {
	ads, err := models.GetAds(database.DB, 0)
	if err != nil || len(ads) == 0 {
		ctx.Send("No ads are available to display.")
		return nil
	}

	for _, ad := range ads {
		msg := fmt.Sprintf(
			"Ad ID: %d\nTitle: %s\nDescription: %s\nPrice: %d\nCity: %s\nNeighborhood: %s\nMeters: %d\nRooms: %d\nPhone: %s\nViews: %d",
			ad.AdID, ad.Title, ad.Description, ad.SellPrice, ad.City, ad.Neighborhood, ad.Meters, ad.NumberOfRooms, ad.SellerPhoneNumber, ad.NumberOfViews)
		ctx.Send(msg)
	}

	return nil
}

func (tg *Telegram) RegisterAdHandlers() {
	tg.Bot.Handle("/add_ad", tg.handleAddAd)
	tg.Bot.Handle("/edit_ad", tg.handleEditAd)
	tg.Bot.Handle("/delete_ad", tg.handleDeleteAd)
	tg.Bot.Handle("/view_ads", tg.handleViewAds)
}
