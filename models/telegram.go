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
	btnAge            = filterMenu.Text("سن بنا")
	btnCategoryAV     = filterMenu.Text("ویلایی یا آپارتمانی؟")
	btnFloorNumber    = filterMenu.Text("طبقه")
	btnStorage        = filterMenu.Text("انباری")
	btnElevator       = filterMenu.Text("آسانسور")
	btnAdDate         = filterMenu.Text("تاریخ ایجاد آگهی")
	btnSendFilter     = filterMenu.Text("ثبت فیلتر")
	btnBackFilterMenu = filterMenu.Text("بازگشت")

	purchaseOrRentMenu = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnPurchase        = purchaseOrRentMenu.Text("خرید")
	btnRent            = purchaseOrRentMenu.Text("اجاره")

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

// -userMenu handlers
func (t *Telegram) handleSetFilters(c telebot.Context) (err error) {
	session := GetUserSession(c.Chat().ID)
	session.State = "selecting_filter"

	filterMenu.Reply(
		filterMenu.Row(btnPrice, btnCity, btnNeighborhood),
		filterMenu.Row(btnArea, btnNumberOfRooms, btnAdDate),
		filterMenu.Row(btnAge, btnCategoryAV, btnFloorNumber),
		filterMenu.Row(btnElevator, btnStorage),
		filterMenu.Row(btnBackFilterMenu),
	)

	t.Bot.Handle(&btnPrice, func(c telebot.Context) (err error) {
		purchaseOrRentMenu.Reply(
			purchaseOrRentMenu.Row(btnPurchase, btnRent),
		)

		t.Bot.Handle(&btnPurchase, func(c telebot.Context) (err error) {
			_, existsMin := session.Filters["min_purchase_price"]
			_, existsMax := session.Filters["max_purchase_price"]
			if existsMin || existsMax {
				err = c.Send("بازه قیمت خرید قبلا مشخص شده است", filterMenu)
				return
			}
			session.State = "setting_purchase_price"
			session.Filters["category_PR"] = "purchase"
			err = c.Send("لطفا بازه قیمت خرید موردنظر را وارد کنید (مثال: 20000-10000)")
			return
		})

		t.Bot.Handle(&btnRent, func(c telebot.Context) (err error) {
			_, existsMin := session.Filters["min_rent_price"]
			_, existsMax := session.Filters["max_rent_price"]
			if existsMin || existsMax {
				err = c.Send("بازه قیمت خرید قبلا مشخص شده است", filterMenu)
				return
			}
			session.State = "setting_rent_price"
			session.Filters["category_PR"] = "rent"
			err = c.Send("لطفا بازه قیمت اجاره و ودیعه موردنظر را به ترتیب وارد کنید (مثال: 2000000-20000 100000-10000)")
			return
		})

		return c.Send("لطفا دسته بندی مورد نظر خود را انتخاب کنید", purchaseOrRentMenu)
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
		err = c.Send("لطفا نام محله را وارد کنید")
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
		err = c.Send("لطفا بازه متراژ را وارد کنید (مثال: 120-100)")
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
		err = c.Send("لطفا بازه تعداد اتاق خواب را وارد کنید (مثال: 3-1)")
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
		err = c.Send("لطفا بازه سن بنا را وارد کنید (مثال: 10-5)")
		return
	})

	t.Bot.Handle(&btnCategoryAV, func(c telebot.Context) (err error) {
		if _, existsCategory := session.Filters["category_AV"]; existsCategory {
			err = c.Send("این دسته بندی قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_category_AV"
		err = c.Send("ویلایی یا آپارتمانی؟")
		return
	})

	t.Bot.Handle(&btnFloorNumber, func(c telebot.Context) (err error) {
		if _, exists := session.Filters["category_AV"]; !exists {
			err = c.Send("لطفا قبل از تعیین بازه طبقه ویلایی یا آپارتمانی بودن آگهی را تعیین کنید")
			return
		}
		if category := session.Filters["category_AV"]; category == "villa" {
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
		err = c.Send("لطفا بازه طبقه را وارد کنید (مثال: 2-1)")
		return
	})

	t.Bot.Handle(&btnStorage, func(c telebot.Context) (err error) {
		if _, existsStorage := session.Filters["storage"]; existsStorage {
			err = c.Send("داشتن یا نداشتن انباری قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_storage"
		err = c.Send("انباری داتشته باشد یا نه؟ (با بله یا خیر پاسخ دهید)")
		return
	})

	t.Bot.Handle(&btnElevator, func(c telebot.Context) (err error) {
		if _, existsElevator := session.Filters["elevator"]; existsElevator {
			err = c.Send("داشتن یا نداشتن آسانسور قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_elevator"
		err = c.Send("آسانسور داشته باشد یا نه؟ (با بله یا خیر پاسخ دهید)")
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
		err = c.Send("لطفا بازه تاریخ ایجاد آگهی را وارد کنید (مثال: 01-01-1403  01-01-1402)")
		return
	})

	// TODO: Mohammad --> insert it into database
	t.Bot.Handle(&btnSendFilter, func(c telebot.Context) (err error) {
		return
	})

	// TODO: Hirad --> back button
	t.Bot.Handle(&btnBackFilterMenu, func(c telebot.Context) (err error) {
		return
	})

	err = c.Send("لطفا یک فیلتر را انتخاب کنید", filterMenu)
	return
}

// TODO
func (t *Telegram) handleShareBookmarks(c telebot.Context) (err error) {
	return nil
}

// TODO
func (t *Telegram) handleGetOutput(c telebot.Context) (err error) {
	getOutputFileMenu.Reply(
		getOutputFileMenu.Row(btnGetAsZip, btnGetViaEmail),
		getOutputFileMenu.Row(btnBackGetOutputFileMenu),
	)

	// TODO
	t.Bot.Handle(&btnGetAsZip, func(c telebot.Context) (err error) {
		return
	})

	// TODO
	t.Bot.Handle(&btnGetViaEmail, func(c telebot.Context) (err error) {
		return
	})

	// TODO: Hirad --> back button
	t.Bot.Handle(&btnBackGetOutputFileMenu, func(c telebot.Context) (err error) {
		return
	})
	return nil
}

// TODO
func (t *Telegram) handleDeleteHistory(c telebot.Context) (err error) {
	return nil
}

// -adminMenu handlers
// TODO
func (t *Telegram) handleSeeCrawlDetails(c telebot.Context) (err error) {
	return nil
}

// -superAdmin menu handlers
// TODO
func (t *Telegram) handleAddAdmin(c telebot.Context) (err error) {
	return nil
}

// TODO
func (t *Telegram) handleManageAdmins(c telebot.Context) (err error) {
	return nil
}

// TODO
func (t *Telegram) handleSetCrawlTimeLimit(c telebot.Context) (err error) {
	return nil
}

// TODO
func (t *Telegram) handleSetNumberOfAds(c telebot.Context) (err error) {
	return nil
}

func (t *Telegram) handleText(c telebot.Context) (err error) {
	session := GetUserSession(c.Chat().ID)
	input := c.Text()

	switch session.State {

	case "setting_city":
		session.Filters["city"] = strings.TrimSpace(strings.ToLower(input))
		session.State = ""
		err = c.Send("شهر مورد نظر ثبت شد", filterMenu)

	case "setting_neighborhood":
		session.Filters["neighborhood"] = strings.TrimSpace(strings.ToLower(input))
		session.State = ""
		err = c.Send("محله مورد نظر ثبت شد", filterMenu)

	case "setting_area":
		minArea, maxArea, e := utils.ParseRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			return
		}
		session.Filters["min_area"] = strconv.Itoa(minArea)
		session.Filters["max_area"] = strconv.Itoa(maxArea)
		session.State = ""
		err = c.Send("بازه متراژ مورد نظر ثبت شد", filterMenu)

	case "setting_number_of_rooms":
		minNumberOfRooms, maxNumberOfRooms, e := utils.ParseRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			return
		}
		session.Filters["min_number_of_rooms"] = strconv.Itoa(minNumberOfRooms)
		session.Filters["max_number_of_rooms"] = strconv.Itoa(maxNumberOfRooms)
		session.State = ""
		err = c.Send("بازه تعداد اتاق خواب مورد نظر ثبت شد", filterMenu)

	case "setting_age":
		minAge, maxAge, e := utils.ParseRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			return
		}
		session.Filters["min_age"] = strconv.Itoa(minAge)
		session.Filters["max_age"] = strconv.Itoa(maxAge)
		session.State = ""
		err = c.Send("بازه سن بنا ثبت شد", filterMenu)

	case "setting_category_AV":
		if input != "آپارتمانی" && input != "ویلایی" {
			err = c.Send("دسته بندی نامعتبر است. لطفا دوباره امتحان کنید")
			return
		}
		if input == "ویلایی" {
			session.Filters["category_AV"] = "villa"
		} else {
			session.Filters["category_AV"] = "apartment"
		}
		session.State = ""
		err = c.Send("دسته بندی مورد نظر ثبت شد", filterMenu)

	case "setting_floor_number":
		minFloorNumber, maxFloorNumber, e := utils.ParseRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			return
		}
		session.Filters["min_floor_number"] = strconv.Itoa(minFloorNumber)
		session.Filters["max_floor_number"] = strconv.Itoa(maxFloorNumber)
		session.State = ""
		err = c.Send("بازه تعداد طبقه مورد نظر ثبت شد", filterMenu)

	case "setting_storage":
		lowerInput := strings.TrimSpace(strings.ToLower(input))
		if lowerInput != "بله" && lowerInput != "خیر" {
			err = c.Send("دسته بندی نامعتبر است. دوباره امتحان کنید")
			return
		}
		session.Filters["storage"] = lowerInput
		session.State = ""
		err = c.Send("ثبت شد", filterMenu)

	case "setting_elevator":
		lowerInput := strings.TrimSpace(strings.ToLower(input))
		if lowerInput != "بله" && lowerInput != "خیر" {
			err = c.Send("دسته بندی نامعتبر است. دوباره امتحان کنید")
			return
		}
		session.Filters["elevator"] = lowerInput
		session.State = ""
		_ = c.Send("ثبت شد", filterMenu)
		for key, value := range session.Filters {
			_ = c.Send(fmt.Sprintf("key={%s}, value={%s}", key, value))
		}

	case "setting_ad_date":
		minDate, maxDate, e := utils.ParseDateRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			return
		}
		session.Filters["min_date"] = minDate
		session.Filters["max_date"] = maxDate
		session.State = ""
		err = c.Send("بازه تاریخ ثبت آگهی مورد نظر ثبت شد", filterMenu)

	case "setting_purchase_price":
		minPurchasePrice, maxPurchasePrice, e := utils.ParseRanges(input)
		if e != nil {
			err = e
			return
		}
		session.Filters["min_purchase_price"] = strconv.Itoa(minPurchasePrice)
		session.Filters["max_purchase_price"] = strconv.Itoa(maxPurchasePrice)
		session.State = ""
		err = c.Send("بازه قیمت خرید مورد نظر ثبت شد", filterMenu)

	case "setting_rent_price":
		ranges := strings.Split(strings.TrimSpace(input), " ")
		minRentPrice, maxRentPrice, e := utils.ParseRanges(ranges[0])
		if e != nil {
			err = e
			return
		}
		session.Filters["min_rent_price"] = strconv.Itoa(minRentPrice)
		session.Filters["max_rent_price"] = strconv.Itoa(maxRentPrice)

		minMortgagePrice, maxMortgagePrice, e := utils.ParseRanges(ranges[1])
		if e != nil {
			err = e
			return
		}
		session.Filters["min_mortgage_price"] = strconv.Itoa(minMortgagePrice)
		session.Filters["max_mortgage_price"] = strconv.Itoa(maxMortgagePrice)
		session.State = ""
		err = c.Send("بازه قیمت مورد نظر ثبت شد", filterMenu)

	default:
		err = c.Send("لطفا از منو آیتم مورد نظر خود را را انتخاب کنید")
	}
	return
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
