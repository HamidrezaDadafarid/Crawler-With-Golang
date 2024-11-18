package main

import (
	"fmt"
	"strings"
)


func (tg *Telegram) handleAddAd(ctx telebot.Context) error {
	input := ctx.Message().Payload
	if input == "" {
		ctx.Send("Please provide the ad information in the format: Title|Description|Price.")
		return nil
	}

	parts := strings.Split(input, "|")
	if len(parts) < 3 {
		ctx.Send("The entered information format is incorrect. Please try again.")
		return nil
	}

	ad := models.Ads{
		Title:       parts[0],
		Description: parts[1],
		SellPrice:   utils.ParseUint(parts[2]),
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
		ctx.Send("Please provide the ad information in the format: ID|New Title|New Description|New Price.")
		return nil
	}

	parts := strings.Split(input, "|")
	if len(parts) < 4 {
		ctx.Send("The entered information format is incorrect. Please try again.")
		return nil
	}

	adID := utils.ParseUint(parts[0])
	var ad models.Ads
	if err := database.DB.First(&ad, adID).Error; err != nil {
		ctx.Send("No ad found with this ID.")
		return nil
	}

	ad.Title = parts[1]
	ad.Description = parts[2]
	ad.SellPrice = utils.ParseUint(parts[3])

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
		msg := fmt.Sprintf("شناسه: %d\nعنوان: %s\nتوضیحات: %s\nقیمت: %d", ad.ID, ad.Title, ad.Description, ad.SellPrice)
		ctx.Send(msg)
	}

	return nil
}


func (tg *Telegram) RegisterAdHandlers() {
	tg.Bot.Handle(&btnAddAd, tg.handleAddAd)
	tg.Bot.Handle(&btnEditAd, tg.handleEditAd)
	tg.Bot.Handle(&btnDeleteAd, tg.handleDeleteAd)
	tg.Bot.Handle(&btnViewAds, tg.handleViewAds)
}


func (tg *Telegram) handleViewAds(ctx telebot.Context) error {
	ads, err := models.GetAds(database.DB, 0)
	if err != nil || len(ads) == 0 {
		ctx.Send("No ads are available to display.")
		return nil
	}

	inlineMenu := &telebot.ReplyMarkup{}
	buttons := []telebot.InlineButton{}

	for _, ad := range ads {
		btn := telebot.InlineButton{
			Unique: fmt.Sprintf("ad_%d", ad.ID),
			Text:   fmt.Sprintf("آگهی: %s | قیمت: %d", ad.Title, ad.SellPrice),
		}
		tg.Bot.Handle(&btn, func(ctx telebot.Context) error {
			msg := fmt.Sprintf(
				"شناسه: %d\nعنوان: %s\nتوضیحات: %s\nقیمت: %d\nمحله: %s",
				ad.ID, ad.Title, ad.Description, ad.SellPrice, ad.Neighborhood,
			)
			ctx.Send(msg)
			return nil
		})
		buttons = append(buttons, btn)
	}

	inlineMenu.Inline(buttons...)
	ctx.Send("Please select one of the ads:", inlineMenu)
	return nil
}


func (tg *Telegram) handleDeleteAd(ctx telebot.Context) error {
	adID := utils.ParseUint(ctx.Message().Payload)
	if adID == 0 {
		ctx.Send("Please enter the ad ID.")
		return nil
	}

	
	inlineMenu := &telebot.ReplyMarkup{}
	btnConfirm := telebot.InlineButton{Unique: "confirm_delete", Text: "Confirm deletion"}
	btnCancel := telebot.InlineButton{Unique: "cancel_delete", Text: "Cancel operation"}

	tg.Bot.Handle(&btnConfirm, func(ctx telebot.Context) error {
		if err := models.DeleteAd(database.DB, adID); err != nil {
			ctx.Send("Error deleting the ad: " + err.Error())
			return nil
		}
		ctx.Send("The ad was successfully deleted!")
		return nil
	})

	tg.Bot.Handle(&btnCancel, func(ctx telebot.Context) error {
		ctx.Send("The deletion operation was canceled.")
		return nil
	})

	inlineMenu.Inline(inlineMenu.Row(btnConfirm, btnCancel))
	ctx.Send("Are you sure you want to delete this ad?", inlineMenu)
	return nil
}


func (tg *Telegram) handleEditAd(ctx telebot.Context) error {
	adID := utils.ParseUint(ctx.Message().Payload)
	if adID == 0 {
		ctx.Send("Please enter the ad ID.")
		return nil
	}

	var ad models.Ads
	if err := database.DB.First(&ad, adID).Error; err != nil {
		ctx.Send("No ad found with this ID.")
		return nil
	}

	ctx.Send(fmt.Sprintf("Current Ad:\nID: %d\nTitle: %s\nPrice: %d", ad.ID, ad.Title, ad.SellPrice))
	ctx.Send("Please provide the new information in the format: Title|Description|Price.")
	return nil
}


func (tg *Telegram) handleViewAds(ctx telebot.Context) error {
    input := ctx.Message().Payload
    var ads []models.Ads
    var err error

    if input != "" {
        parts := strings.Split(input, "|")
        filters := make(map[string]interface{})

        for _, part := range parts {
            filter := strings.Split(part, "=")
            if len(filter) == 2 {
                filters[filter[0]] = filter[1]
            }
        }
        ads, err = models.GetFilteredAds(database.DB, filters)
    } else {
        ads, err = models.GetAds(database.DB, 0)
    }

    if err != nil || len(ads) == 0 {
        ctx.Send("No ads available for display with the specified filters.")
        return nil
    }

    for _, ad := range ads {
        msg := fmt.Sprintf("ID: %d\nTitle: %s\nDescription: %s\nPrice: %d\nLocation: %s",
            ad.ID, ad.Title, ad.Description, ad.SellPrice, ad.Neighborhood)
        ctx.Send(msg)
    }

    return nil
}


func (tg *Telegram) handleAddAd(ctx telebot.Context) error {
    input := ctx.Message().Payload
    if input == "" {
        ctx.Send("Please provide the ad information in the format: Title|Description|Price.")
        return nil
    }

    parts := strings.Split(input, "|")
    if len(parts) < 3 {
        ctx.Send("The entered information format is incorrect. Please try again.")
        return nil
    }

    price, err := utils.ParseUint(parts[2])
    if err != nil {
        ctx.Send("Invalid price format. Please enter a numeric value.")
        return nil
    }

    ad := models.Ads{
        Title:       parts[0],
        Description: parts[1],
        SellPrice:   price,
    }

    if err := models.AddAd(database.DB, &ad); err != nil {
        ctx.Send("Error adding the ad: " + err.Error())
        return nil
    }

    ctx.Send("The ad was added successfully!")
    return nil
}
