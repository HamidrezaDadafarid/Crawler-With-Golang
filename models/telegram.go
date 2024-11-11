package models

import (
	"fmt"
	"log"
	"main/utils"
	"strconv"
	"strings"
	"time"

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
	superAdminMenu       = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnAddAdmin          = superAdminMenu.Text("اضافه کردن ادمین")
	btnManageAdmins      = superAdminMenu.Text("مدیریت کردن ادمین ها")
	btnSetNumberOfAds    = superAdminMenu.Text("تنظیم تعداد آیتم های جستجو شده")
	btnSetCrawlTimeLimit = superAdminMenu.Text("تنظیم محدودیت زمانی فرآیند جستجو")

	adminMenu          = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnSeeCrawlDetails = superAdminMenu.Text("See Crawl Details")

	userMenu          = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnSetFilters     = userMenu.Text("ثبت فیلتر")
	btnShareBookmarks = userMenu.Text("اشتراک گذاری آگهی های مورد علاقه")
	btnGetOutputFile  = userMenu.Text("خروجی گرفتن از آگهی ها")
	btnDeleteHistory  = userMenu.Text("پاک کردن تاریخچه")
	// تنظیم بازه زمانی

	filterMenu        = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnPrice          = filterMenu.Text("قیمت")
	btnCity           = filterMenu.Text("شهر")
	btnNeighborhood   = filterMenu.Text("محله")
	btnArea           = filterMenu.Text("متراژ")
	btnNumberOfRooms  = filterMenu.Text("تعداد اتاق خواب")
	btnCategoryPMR    = filterMenu.Text("خرید یا اجاره؟")
	btnAge            = filterMenu.Text("سن بنا")
	btnCategoryAV     = filterMenu.Text("ویلایی یا آپارتمانی؟")
	btnFloorNumber    = filterMenu.Text("طبقه")
	btnStorage        = filterMenu.Text("انباری")
	btnElevator       = filterMenu.Text("آسانسور")
	btnAdDate         = filterMenu.Text("تاریخ ایجاد آگهی")
	btnSendFilter     = filterMenu.Text("ثبت فیلتر")
	btnBackFilterMenu = filterMenu.Text("بازگشت")

	getOutputFileMenu        = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnGetAsZip              = getOutputFileMenu.Text("دریافت آگهی ها به صورت فایل زیپ")
	btnGetViaEmail           = getOutputFileMenu.Text("دریافت آگهی ها از طریق ایمیل")
	btnBackGetOutputFileMenu = getOutputFileMenu.Text("بازگشت")
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
	t.Bot.Handle("/start", t.handleStart)
	t.Bot.Handle(telebot.OnText, t.handleText)
}

func (t *Telegram) Start() {
	t.registerHandlers()
	log.Println("Starting Telegram bot...")
	t.Bot.Start()
}

func (t *Telegram) handleStart(c telebot.Context) (err error) {
	welcomeMsg := "به ربات خزنده خوش اومدین :)"
	// TODO
	// get user's role with this telegram id from database
	// telegram_ID := c.Sender().ID
	// role := SELECT role FROM users WHERE telegram_ID = telegram_ID
	role := "user"
	if role == "user" {
		userMenu.Reply(
			userMenu.Row(btnSetFilters, btnShareBookmarks),
			userMenu.Row(btnGetOutputFile, btnDeleteHistory),
		)

		t.Bot.Handle(&btnSetFilters, t.handleSetFilters)
		t.Bot.Handle(&btnShareBookmarks, t.handleShareBookmarks)
		t.Bot.Handle(&btnGetOutputFile, t.handleGetOutput)
		t.Bot.Handle(&btnDeleteHistory, t.handleDeleteHistory)

		err = c.Send(welcomeMsg, userMenu)

	} else if role == "admin" {
		adminMenu.Reply(
			adminMenu.Row(btnSeeCrawlDetails),
		)

		t.Bot.Handle(&btnSeeCrawlDetails, t.handleSeeCrawlDetails)

		err = c.Send(welcomeMsg, adminMenu)

	} else {
		superAdminMenu.Reply(
			superAdminMenu.Row(btnAddAdmin, btnManageAdmins),
			superAdminMenu.Row(btnSetCrawlTimeLimit, btnSetNumberOfAds),
		)

		t.Bot.Handle(&btnAddAdmin, t.handleAddAdmin)
		t.Bot.Handle(&btnManageAdmins, t.handleManageAdmins)
		t.Bot.Handle(&btnSetCrawlTimeLimit, t.handleSetCrawlTimeLimit)
		t.Bot.Handle(&btnSetNumberOfAds, t.handleSetNumberOfAds)

		err = c.Send(welcomeMsg, superAdminMenu)
	}
	return
}

// User menu handlers
func (t *Telegram) handleSetFilters(c telebot.Context) (err error) {
	session := GetUserSession(c.Chat().ID)
	session.State = "selecting_filter"

	filterMenu.Reply(
		filterMenu.Row(btnPrice, btnCity, btnNeighborhood),
		filterMenu.Row(btnArea, btnNumberOfRooms, btnCategoryPMR),
		filterMenu.Row(btnAge, btnCategoryAV, btnFloorNumber),
		filterMenu.Row(btnStorage, btnElevator, btnAdDate),
		filterMenu.Row(btnSendFilter, btnBackFilterMenu),
	)

	t.Bot.Handle(&btnPrice, func(c telebot.Context) (err error) {
		if _, exists := session.Filters["Category_PR"]; !exists {
			err = c.Send("لطفا قبل از تعیین بازه قیمت، فیلتر خرید یا اجاره را تعیین کنید")
			return
		}

		if category := session.Filters["category_PR"]; category == "purchase" {
			_, existsMin := session.Filters["purchase_min_price"]
			_, existsMax := session.Filters["purchase_max_price"]
			if existsMin || existsMax {
				err = c.Send("بازه قیمت خرید قبلا مشخص شده است", filterMenu)
				return
			}

			session.State = "setting_purchase_price_range"
			err = c.Send("لطفا بازه قیمت خود را وارد کنید: (مثال: 1000-100)")
			return
		} else {
			_, existsMin := session.Filters["rent_min_price"]
			_, existsMax := session.Filters["rent_max_price"]
			if existsMin || existsMax {
				err = c.Send("بازه قیمت اجاره قبلا مشخص شده است", filterMenu)
				return
			}
			session.State = "setting_rent_price_range"
			err = c.Send("لطفا بازه قیمت خود را وارد کنید: (مثال: 1000-100)")
			return
		}

	})

	t.Bot.Handle(&btnCity, func(c telebot.Context) (err error) {
		if _, existsCity := session.Filters["city"]; existsCity {
			err = c.Send("شهر قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_city"
		err = c.Send("لطفا نام شهر را وارد کنید")
		return
	})

	t.Bot.Handle(&btnNeighborhood, func(c telebot.Context) (err error) {
		if _, existsNeighborhood := session.Filters["neighborhood"]; existsNeighborhood {
			err = c.Send("محله قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_neighborhood"
		err = c.Send("لطفا نام محله را وارد کنید", filterMenu)
		return
	})

	t.Bot.Handle(&btnArea, func(c telebot.Context) (err error) {
		_, existsMin := session.Filters["min_area"]
		_, existsMax := session.Filters["max_area"]
		if existsMin || existsMax {
			err = c.Send("بازه متراژ قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_area"
		err = c.Send("لطفا بازه متراژ را وارد کنید (مثال: 120-100)", filterMenu)
		return
	})

	t.Bot.Handle(&btnNumberOfRooms, func(c telebot.Context) (err error) {
		_, existsMin := session.Filters["min_number_of_rooms"]
		_, existsMax := session.Filters["max_number_of_rooms"]
		if existsMin || existsMax {
			err = c.Send("بازه تعداد اتاق خواب قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_number_of_rooms"
		err = c.Send("لطفا بازه تعداد اتاق خواب را وارد کنید (مثال: 3-1)", filterMenu)
		return
	})

	t.Bot.Handle(&btnCategoryPMR, func(c telebot.Context) (err error) {
		if _, existsCategory := session.Filters["category_PMR"]; existsCategory {
			err = c.Send("این دسته بندی قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_category_PMR"
		err = c.Send("خرید یا اجاره؟", filterMenu)
		return
	})

	t.Bot.Handle(&btnAge, func(c telebot.Context) (err error) {
		_, existsMin := session.Filters["min_age"]
		_, existsMax := session.Filters["max_age"]
		if existsMin || existsMax {
			err = c.Send("بازه سن بنا قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_age"
		err = c.Send("لطفا بازه سن بنا را وارد کنید (مثال: 10-5)", filterMenu)
		return
	})

	t.Bot.Handle(&btnCategoryAV, func(c telebot.Context) (err error) {
		if _, existsCategory := session.Filters["category_AV"]; existsCategory {
			err = c.Send("این دسته بندی قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_category_AV"
		err = c.Send("ویلایی یا آپارتمانی؟", filterMenu)
		return
	})

	t.Bot.Handle(&btnFloorNumber, func(c telebot.Context) (err error) {
		if _, exists := session.Filters["Category_AV"]; !exists {
			err = c.Send("لطفا قبل از تعیین بازه طبقه ویلایی یا آپارتمانی بودن آگهی را تعیین کنید")
			return
		}
		if category := session.Filters["Category_AV"]; category == "villa" {
			err = c.Send("برای آگهی های ویلایی نمی توان بازه طبقه مشخص کرد")
			return
		}
		_, existsMin := session.Filters["min_floor_number"]
		_, existsMax := session.Filters["max_floor_number"]
		if existsMin || existsMax {
			err = c.Send("بازه طبقه قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_floor_number"
		err = c.Send("لطفا بازه طبقه را وارد کنید (مثال: 2-1)", filterMenu)
		return
	})

	t.Bot.Handle(&btnStorage, func(c telebot.Context) (err error) {
		if _, existsStorage := session.Filters["storage"]; existsStorage {
			err = c.Send("داشتن یا نداشتن انباری قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_storage"
		err = c.Send("انباری داتشته باشد یا نه؟ (با بله یا خیر پاسخ دهید)", filterMenu)
		return
	})

	t.Bot.Handle(&btnElevator, func(c telebot.Context) (err error) {
		if _, existsElevator := session.Filters["elevator"]; existsElevator {
			err = c.Send("داشتن یا نداشتن آسانسور قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_elevator"
		err = c.Send("آسانسور داشته باشد یا نه؟ (با بله یا خیر پاسخ دهید)", filterMenu)
		return
	})

	t.Bot.Handle(&btnAdDate, func(c telebot.Context) (err error) {
		_, existsMin := session.Filters["min_date"]
		_, existsMax := session.Filters["max_date"]
		if existsMin || existsMax {
			err = c.Send("بازه تاریخ ایجاد آگهی قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_ad_date"
		err = c.Send("لطفا بازه تاریخ ایجاد آگهی را وارد کنید (مثال: 01-01-1403  01-01-1402)", filterMenu)
		return
	})

	// TODO: insert it into database
	t.Bot.Handle(&btnSendFilter, func(c telebot.Context) (err error) {
		return
	})

	// TODO: Hirad
	t.Bot.Handle(&btnBackFilterMenu, func(c telebot.Context) (err error) {
		session.State = ""
		t.Bot.Handle(&btnSetFilters, t.handleSetFilters)
		t.Bot.Handle(&btnShareBookmarks, t.handleShareBookmarks)
		t.Bot.Handle(&btnGetOutputFile, t.handleGetOutput)
		t.Bot.Handle(&btnDeleteHistory, t.handleDeleteHistory)

		err = c.Send("به صفحه اصلی بازگشتید", userMenu)
		return
	})

	err = c.Send("لطفا یک فیلتر را انتخاب کنید", filterMenu)
	return
}

// TODO
func (t *Telegram) handleShareBookmarks(c telebot.Context) error {
	return nil
}

// TODO
func (t *Telegram) handleGetOutput(c telebot.Context) error {
	session := GetUserSession(c.Chat().ID)
	getOutputFileMenu.Reply(
		getOutputFileMenu.Row(btnGetAsZip, btnGetViaEmail),
		getOutputFileMenu.Row(btnBackFilterMenu),
	)

	// TODO
	t.Bot.Handle(&btnGetOutputFile, t.handleGetOutputFile)
	t.Bot.Handle(&btnGetViaEmail, t.handleGetOutputViaEmail)
	// TODO: Hirad
	t.Bot.Handle(&btnBackGetOutputFileMenu, func(c telebot.Context) (err error) {

		session.State = ""
		t.Bot.Handle(&btnSetFilters, t.handleSetFilters)
		t.Bot.Handle(&btnShareBookmarks, t.handleShareBookmarks)
		t.Bot.Handle(&btnGetOutputFile, t.handleGetOutput)
		t.Bot.Handle(&btnDeleteHistory, t.handleDeleteHistory)

		err = c.Send("به صفحه اصلی بازگشتید", userMenu)
		return
	})
	return c.Send("نحوه خروجی گرفتن را انتخاب کنید", getOutputFileMenu) //TODO: change text message
}

// TODO
func (t *Telegram) handleDeleteHistory(c telebot.Context) error {
	return nil
}

// TODO
func (t *Telegram) handleGetOutputFile(c telebot.Context) error {
	return nil
}

// TODO
func (t *Telegram) handleGetOutputViaEmail(c telebot.Context) error {
	return nil
}

// Admin menu handlers
// TODO
func (t *Telegram) handleSeeCrawlDetails(c telebot.Context) error {
	return nil
}

// Super admin menu handlers
// TODO
func (t *Telegram) handleAddAdmin(c telebot.Context) error {
	return nil
}

// TODO
func (t *Telegram) handleManageAdmins(c telebot.Context) error {
	return nil
}

// TODO
func (t *Telegram) handleSetCrawlTimeLimit(c telebot.Context) error {
	return nil
}

// TODO
func (t *Telegram) handleSetNumberOfAds(c telebot.Context) error {
	return nil
}

func (t *Telegram) handleText(c telebot.Context) (err error) {
	session := GetUserSession(c.Chat().ID)
	input := c.Text()

	switch session.State {
	case "setting_price_range":
		if filter := session.Filters["category_PR"]; filter == "purchase" {
			// TODO
		} else {
			// TODO
		}
		minPrice, maxPrice, e := utils.ParseRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			return
		}

		session.Filters["min_price"] = strconv.Itoa(minPrice)
		session.Filters["max_price"] = strconv.Itoa(maxPrice)
		session.State = ""
		return c.Send("بازه قیمت مورد نظر ثبت شد", filterMenu)

	case "setting_city":
		session.Filters["city"] = strings.TrimSpace(strings.ToLower(input))
		session.State = ""
		return c.Send("شهر مورد نظر ثبت شد", filterMenu)

	case "setting_neighborhood":
		session.Filters["neighborhood"] = strings.TrimSpace(strings.ToLower(input))
		session.State = ""
		return c.Send("محله مورد نظر ثبت شد", filterMenu)

	case "setting_area":
		minArea, maxArea, err := utils.ParseRanges(input)
		if err != nil {
			return c.Send(err.Error())
		}
		session.Filters["min_area"] = strconv.Itoa(minArea)
		session.Filters["max_area"] = strconv.Itoa(maxArea)
		session.State = ""
		return c.Send("بازه متراژ مورد نظر ثبت شد", filterMenu)

	case "setting_number_of_rooms":
		minNumberOfRooms, maxNumberOfRooms, err := utils.ParseRanges(input)
		if err != nil {
			return c.Send(err.Error())
		}
		session.Filters["min_number_of_rooms"] = strconv.Itoa(minNumberOfRooms)
		session.Filters["max_number_of_rooms"] = strconv.Itoa(maxNumberOfRooms)
		session.State = ""
		return c.Send("بازه تعداد اتاق خواب مورد نظر ثبت شد", filterMenu)

	case "setting_category_PMR":
		lowerInput := strings.TrimSpace(strings.ToLower(input))
		if lowerInput != "خرید" && lowerInput != "اجاره" {
			return c.Send("دسته بندی نامعتبر است. لطفا دوباره امتحان کنید")
		}
		session.Filters["category_PMR"] = lowerInput
		session.State = ""
		return c.Send("دسته بندی مورد نظر ثبت شد", filterMenu)

	case "setting_age":
		minAge, maxAge, err := utils.ParseRanges(input)
		if err != nil {
			return c.Send(err.Error())
		}
		session.Filters["min_age"] = strconv.Itoa(minAge)
		session.Filters["max_age"] = strconv.Itoa(maxAge)
		session.State = ""
		return c.Send("بازه سن بنا ثبت شد", filterMenu)

	case "setting_category_AV":
		lowerInput := strings.TrimSpace(strings.ToLower(input))
		if lowerInput != "آپارتمانی" && lowerInput != "ویلایی" {
			return c.Send("دسته بندی نامعتبر است. لطفا دوباره امتحان کنید")
		}
		session.Filters["category_AV"] = lowerInput
		session.State = ""
		return c.Send("دسته بندی مورد نظر ثبت شد", filterMenu)

	case "setting_floor_number":
		minFloorNumber, maxFloorNumber, err := utils.ParseRanges(input)
		if err != nil {
			return c.Send(err.Error())
		}
		session.Filters["min_floor_number"] = strconv.Itoa(minFloorNumber)
		session.Filters["max_floor_number"] = strconv.Itoa(maxFloorNumber)
		session.State = ""
		return c.Send("بازه تعداد طبقه مورد نظر ثبت شد", filterMenu)

	case "setting_storage":
		lowerInput := strings.TrimSpace(strings.ToLower(input))
		if lowerInput != "بله" && lowerInput != "خیر" {
			return c.Send("دسته بندی نامعتبر است. دوباره امتحان کنید")
		}
		session.Filters["storage"] = lowerInput
		session.State = ""
		return c.Send("ثبت شد", filterMenu)

	case "setting_elevator":
		lowerInput := strings.TrimSpace(strings.ToLower(input))
		if lowerInput != "بله" && lowerInput != "خیر" {
			return c.Send("دسته بندی نامعتبر است. دوباره امتحان کنید")
		}
		session.Filters["elevator"] = lowerInput
		session.State = ""
		return c.Send("ثبت شد", filterMenu)

	case "setting_ad_date":
		minDate, maxDate, err := utils.ParseDateRanges(input)
		if err != nil {
			return c.Send(err.Error())
		}
		session.Filters["min_date"] = minDate
		session.Filters["max_date"] = maxDate
		session.State = ""
		return c.Send("بازه تاریخ ثبت آگهی مورد نظر ثبت شد", filterMenu)

	default:
		return c.Send("لطفا از منو آیتم مورد نظر خود را را انتخاب کنید")
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
