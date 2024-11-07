package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"main/utils"

	"gopkg.in/telebot.v4"
)

type TelegramConfig struct {
	Token string
}

type Telegram struct {
	Bot    *telebot.Bot
	Config *TelegramConfig
}

var (
	startMenu            = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnLoginAsSuperAdmin = startMenu.Text("Login As Super Admin")
	btnLoginAsAdmin      = startMenu.Text("Login As Admin")
	btnLoginAsUser       = startMenu.Text("Login As User")

	userMenu          = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnSetFilters     = userMenu.Text("Set Filters")
	btnShareBookmarks = userMenu.Text("Share Bookmarks")
	btnGetOutputFile  = userMenu.Text("Get Output File")
	btnDeleteHistory  = userMenu.Text("Detele History")

	filterMenu       = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnPrice         = filterMenu.Text("Price")
	btnCity          = filterMenu.Text("City")
	btnNeighborhood  = filterMenu.Text("Neighborhood")
	btnArea          = filterMenu.Text("Area")
	btnNumberOfRooms = filterMenu.Text("Number Of Rooms")
	btnCategoryPMR   = filterMenu.Text("Purchase, Mortgage, Rent Category")
	btnAge           = filterMenu.Text("Age")
	btnCategoryAV    = filterMenu.Text("Apartment Or Villa Category")
	btnFloorNumber   = filterMenu.Text("Floor Number")
	btnStorage       = filterMenu.Text("Storage")
	btnElevator      = filterMenu.Text("Elevator")
	btnAdDate        = filterMenu.Text("Ad Date")
	btnSendFilter    = filterMenu.Text("Send Filter")

	getOutputFileMenu = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnGetAsZip       = getOutputFileMenu.Text("Get Ads As Zip File")
	btnGetViaEmail    = getOutputFileMenu.Text("Get Ads Via Email")
)

func NewTelegramBot(config *TelegramConfig) (*Telegram, error) {
	pref := telebot.Settings{
		Token:  config.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	telegram := &Telegram{
		Bot:    bot,
		Config: config,
	}

	return telegram, nil
}

func (t *Telegram) registerHandlers() {
	t.Bot.Handle(telebot.OnText, t.handleText)
	t.Bot.Handle("/start", t.handleStart)
}

func (t *Telegram) Start() {
	t.registerHandlers()
	log.Println("Starting Telegram bot...")
	t.Bot.Start()
}

func (t *Telegram) handleStart(c telebot.Context) error {
	welcomeMsg := "Welcome To The Real Estate Bot."

	startMenu.Reply(
		startMenu.Row(btnLoginAsSuperAdmin, btnLoginAsAdmin, btnLoginAsUser),
	)

	t.Bot.Handle(&btnLoginAsSuperAdmin, t.handleLoginSuperAdmin)
	t.Bot.Handle(&btnLoginAsAdmin, t.handleLoginAdmin)
	t.Bot.Handle(&btnLoginAsUser, t.handleLoginUser)

	return c.Send(welcomeMsg, startMenu)
}

// TODO
func (t *Telegram) handleLoginSuperAdmin(c telebot.Context) error {
	return nil
}

// TODO
func (t *Telegram) handleLoginAdmin(c telebot.Context) error {
	return nil
}

func (t *Telegram) handleLoginUser(c telebot.Context) error {
	userMenu.Reply(
		userMenu.Row(btnSetFilters, btnShareBookmarks),
		userMenu.Row(btnGetOutputFile, btnDeleteHistory),
	)

	t.Bot.Handle(&btnSetFilters, t.handleFilters)
	t.Bot.Handle(&btnShareBookmarks, t.handleShareBookmarks)
	t.Bot.Handle(&btnGetOutputFile, t.handleGetOutput)
	t.Bot.Handle(&btnDeleteHistory, t.handleDeleteHistory)

	return c.Send("Please select an item", userMenu)
}

func (t *Telegram) handleFilters(c telebot.Context) error {
	session := GetUserSession(c.Chat().ID)
	session.State = "selecting_filter"

	filterMenu.Reply(
		filterMenu.Row(btnPrice, btnCity, btnNeighborhood),
		filterMenu.Row(btnArea, btnNumberOfRooms, btnCategoryPMR),
		filterMenu.Row(btnAge, btnCategoryAV, btnFloorNumber),
		filterMenu.Row(btnStorage, btnElevator, btnAdDate),
		filterMenu.Row(btnSendFilter),
	)

	t.Bot.Handle(&btnPrice, func(c telebot.Context) error {
		_, existsMin := session.Filters["min_price"]
		_, existsMax := session.Filters["max_price"]
		if existsMin || existsMax {
			return c.Send("Price range is already set.", filterMenu)
		}
		session.State = "setting_price_range"
		return c.Send("Please enter the price range (e.g., 100000-500000)")
	})

	t.Bot.Handle(&btnCity, func(c telebot.Context) error {
		if _, existsCity := session.Filters["city"]; existsCity {
			return c.Send("City is already set.", filterMenu)
		}
		session.State = "setting_city"
		return c.Send("Please enter the city name:")
	})

	t.Bot.Handle(&btnNeighborhood, func(c telebot.Context) error {
		if _, existsNeighborhood := session.Filters["neighborhood"]; existsNeighborhood {
			return c.Send("Neighborhood is already set.", filterMenu)
		}
		session.State = "setting_neighborhood"
		return c.Send("Please enter the neighborhood name")
	})

	t.Bot.Handle(&btnArea, func(c telebot.Context) error {
		_, existsMin := session.Filters["min_area"]
		_, existsMax := session.Filters["max_area"]
		if existsMin || existsMax {
			return c.Send("Area range is already set.", filterMenu)
		}
		session.State = "setting_area"
		return c.Send("Please enter the area range (e.g 100-150)")
	})

	t.Bot.Handle(&btnNumberOfRooms, func(c telebot.Context) error {
		_, existsMin := session.Filters["min_number_of_rooms"]
		_, existsMax := session.Filters["max_number_of_rooms"]
		if existsMin || existsMax {
			return c.Send("Number of rooms range is already set.", filterMenu)
		}
		session.State = "setting_number_of_rooms"
		return c.Send("Please enter the number of rooms range (e.g 1-3)")
	})

	t.Bot.Handle(&btnCategoryPMR, func(c telebot.Context) error {
		if _, existsCategory := session.Filters["category_PMR"]; existsCategory {
			return c.Send("Category is already set.", filterMenu)
		}
		session.State = "setting_category_PMR"
		return c.Send("Please enter the category (purchase, mortgage or rent)")
	})

	t.Bot.Handle(&btnAge, func(c telebot.Context) error {
		_, existsMin := session.Filters["min_age"]
		_, existsMax := session.Filters["max_age"]
		if existsMin || existsMax {
			return c.Send("Age range is already set.", filterMenu)
		}
		session.State = "setting_age"
		return c.Send("Please enter the age range (e.g 3-8)")
	})

	t.Bot.Handle(&btnCategoryAV, func(c telebot.Context) error {
		if _, existsCategory := session.Filters["category_AV"]; existsCategory {
			return c.Send("Category is already set.", filterMenu)
		}
		session.State = "setting_category_AV"
		return c.Send("Please enter the category (apartment or villa)")
	})

	t.Bot.Handle(&btnFloorNumber, func(c telebot.Context) error {
		_, existsMin := session.Filters["min_floor_number"]
		_, existsMax := session.Filters["max_floor_number"]
		if existsMin || existsMax {
			return c.Send("Floor number range is already set.", filterMenu)
		}
		session.State = "setting_floor_number"
		return c.Send("Please enter the floor number range (e.g 1-4)")
	})

	t.Bot.Handle(&btnStorage, func(c telebot.Context) error {
		if _, existsStorage := session.Filters["storage"]; existsStorage {
			return c.Send("Storage is already set.", filterMenu)
		}
		session.State = "setting_storage"
		return c.Send("Please enter whether you want the house to have storage or not (yes or no)")
	})

	t.Bot.Handle(&btnElevator, func(c telebot.Context) error {
		if _, existsElevator := session.Filters["elevator"]; existsElevator {
			return c.Send("Elevator is already set.", filterMenu)
		}
		session.State = "setting_elevator"
		return c.Send("Please enter whether you want the house to have elevator or not (yes or no)")
	})

	t.Bot.Handle(&btnAdDate, func(c telebot.Context) error {
		_, existsMin := session.Filters["min_date"]
		_, existsMax := session.Filters["max_date"]
		if existsMin || existsMax {
			return c.Send("Date range is already set.", filterMenu)
		}
		session.State = "setting_ad_date"
		return c.Send("Please enter the ad date range (e.g 1403-8-10 1403-10-8)")
	})

	// TODO
	t.Bot.Handle(&btnSendFilter, func(c telebot.Context) error {
		return nil
	})

	return c.Send("Select a filter to set", filterMenu)
}

// TODO
func (t *Telegram) handleShareBookmarks(c telebot.Context) error {
	return nil
}

// TODO
func (t *Telegram) handleGetOutput(c telebot.Context) error {
	getOutputFileMenu.Reply(
		getOutputFileMenu.Row(btnGetAsZip, btnGetViaEmail),
	)

	// TODO
	// t.Bot.Handle(&btnGetOutputFile, t.handleGetOutputFile)
	// t.Bot.Handle(&btnGetViaEmail, t.GetOutputViaEmail)
	return nil
}

// TODO
func (t *Telegram) handleDeleteHistory(c telebot.Context) error {
	return nil
}

func (t *Telegram) handleText(c telebot.Context) error {
	session := GetUserSession(c.Chat().ID)
	input := c.Text()

	switch session.State {
	case "setting_price_range":
		minPrice, maxPrice, err := utils.ParseRanges(input)
		if err != nil {
			return c.Send(fmt.Sprintf("%s. Please try again.", err.Error()))
		}
		session.Filters["min_price"] = strconv.Itoa(minPrice)
		session.Filters["max_price"] = strconv.Itoa(maxPrice)
		session.State = ""
		return c.Send("Price range set successfully.", filterMenu)

	case "setting_city":
		session.Filters["city"] = strings.TrimSpace(strings.ToLower(input))
		session.State = ""
		return c.Send("City set successfully.", filterMenu)

	case "setting_neighborhood":
		session.Filters["neighborhood"] = strings.TrimSpace(strings.ToLower(input))
		session.State = ""
		return c.Send("Neighborhood set seccessfully.", filterMenu)

	case "setting_area":
		minArea, maxArea, err := utils.ParseRanges(input)
		if err != nil {
			return c.Send(fmt.Sprintf("%s. Please try again.", err.Error()))
		}
		session.Filters["min_area"] = strconv.Itoa(minArea)
		session.Filters["max_area"] = strconv.Itoa(maxArea)
		session.State = ""
		return c.Send("Area range set seccessfully.", filterMenu)

	case "setting_number_of_rooms":
		minNumberOfRooms, maxNumberOfRooms, err := utils.ParseRanges(input)
		if err != nil {
			return c.Send(fmt.Sprintf("%s. Please try again.", err.Error()))
		}
		session.Filters["min_number_of_rooms"] = strconv.Itoa(minNumberOfRooms)
		session.Filters["max_number_of_rooms"] = strconv.Itoa(maxNumberOfRooms)
		session.State = ""
		return c.Send("Number of rooms range set seccessfully.", filterMenu)

	case "setting_category_PMR":
		lowerInput := strings.TrimSpace(strings.ToLower(input))
		if lowerInput != "purchase" && lowerInput != "mortgage" && lowerInput != "rent" {
			return c.Send("Invalid category. Please try again.")
		}
		session.Filters["category_PMR"] = lowerInput
		session.State = ""
		return c.Send("category set seccussfully", filterMenu)

	case "setting_age":
		minAge, maxAge, err := utils.ParseRanges(input)
		if err != nil {
			return c.Send(fmt.Sprintf("%s. Please try again.", err.Error()))
		}
		session.Filters["min_age"] = strconv.Itoa(minAge)
		session.Filters["max_age"] = strconv.Itoa(maxAge)
		session.State = ""
		return c.Send("Age range set seccessfully.", filterMenu)

	case "setting_category_AV":
		lowerInput := strings.TrimSpace(strings.ToLower(input))
		if lowerInput != "apartment" && lowerInput != "villa" {
			return c.Send("Invalid category. Please try again.")
		}
		session.Filters["category_AV"] = lowerInput
		session.State = ""
		return c.Send("category set successfully.", filterMenu)

	case "setting_floor_number":
		minFloorNumber, maxFloorNumber, err := utils.ParseRanges(input)
		if err != nil {
			return c.Send(fmt.Sprintf("%s. Please try again.", err.Error()))
		}
		session.Filters["min_floor_number"] = strconv.Itoa(minFloorNumber)
		session.Filters["max_floor_number"] = strconv.Itoa(maxFloorNumber)
		session.State = ""
		return c.Send("Number of floor range set seccessfully.", filterMenu)

	case "setting_storage":
		lowerInput := strings.TrimSpace(strings.ToLower(input))
		if lowerInput != "yes" && lowerInput != "no" {
			return c.Send("Invalid input format. Please try again.")
		}
		session.Filters["storage"] = lowerInput
		session.State = ""
		return c.Send("Storage set successfully.", filterMenu)

	case "setting_elevator":
		lowerInput := strings.TrimSpace(strings.ToLower(input))
		if lowerInput != "yes" && lowerInput != "no" {
			return c.Send("Invalid input format. Please try again.")
		}
		session.Filters["elevator"] = lowerInput
		session.State = ""
		return c.Send("Elevator set successfully.", filterMenu)

	case "setting_ad_date":
		minDate, maxDate, err := utils.ParseDateRanges(input)
		if err != nil {
			return c.Send(fmt.Sprintf("%s. Please try again.", err.Error()))
		}
		session.Filters["min_date"] = minDate
		session.Filters["max_date"] = maxDate
		session.State = ""
		return c.Send("Date range set successfully.", filterMenu)

	default:
		return c.Send("Please use the menu to select options.")
	}
}

// func (t *Telegram) handleSearch(c telebot.Context) error {
// 	session := GetUserSession(c.Chat().ID)
// 	if len(session.Filters) == 0 {
// 		return c.Send("No filters set. Please set at least one filter before searching using /filters.")
// 	}

// 	// Call the crawler's search function
// 	ads, err := SearchAds(session.Filters) // Assume this function exists in your crawler package
// 	if err != nil {
// 		return c.Send("Error during search: " + err.Error())
// 	}

// 	if len(ads) == 0 {
// 		return c.Send("No ads found matching your criteria.")
// 	}

// 	// Send the ads to the user
// 	for _, ad := range ads {
// 		if err := t.sendAd(c, ad); err != nil {
// 			log.Println("Error sending ad:", err)
// 		}
// 	}

// 	return nil
// }

// func (t *Telegram) sendAd(c telebot.Context, ad Ad) error {
// 	// Send the ad photo if available
// 	if ad.PhotoURL != "" {
// 		photo := &telebot.Photo{File: telebot.FromURL(ad.PhotoURL)}
// 		if _, err := t.Bot.Send(c.Chat(), photo); err != nil {
// 			return err
// 		}
// 	}

// 	// Prepare ad details message
// 	adText := fmt.Sprintf("*%s*\nPrice: %s\n\n%s\n[View Ad](%s)", ad.Title, ad.Price, ad.Description, ad.URL)

// 	// Send the ad details
// 	return c.Send(adText, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
// }
