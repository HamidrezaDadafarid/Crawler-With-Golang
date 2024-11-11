package models

import (
	"fmt"
	"log"
	"main/models"
	"main/utils"
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
	session := models.GetUserSession(c.Chat().ID)
	session.State = "selecting_filter"

	filterMenu.Reply(
		filterMenu.Row(btnPrice, btnCity, btnNeighborhood),
		filterMenu.Row(btnArea, btnNumberOfRooms, btnAdDate),
		filterMenu.Row(btnAge, btnCategoryAV, btnFloorNumber),
		filterMenu.Row(btnElevator, btnStorage),
		filterMenu.Row(btnSendFilter, btnBackFilterMenu),
	)

	t.Bot.Handle(&btnPrice, func(c telebot.Context) (err error) {
		purchaseOrRentMenu.Reply(
			purchaseOrRentMenu.Row(btnPurchase, btnRent),
		)

		t.Bot.Handle(&btnPurchase, func(c telebot.Context) (err error) {
			min := session.Filters.StartPurchasePrice
			max := session.Filters.EndPurchasePrice
			if min != nil || max != nil {
				err = c.Send("بازه قیمت خرید قبلا مشخص شده است", filterMenu)
				return
			}
			session.State = "setting_purchase_price"
			if session.Filters.CategoryPR == nil {
				session.Filters.CategoryPR = new(uint)
			}
			*session.Filters.CategoryPR = 0 // 0 for purchase
			err = c.Send("لطفا بازه قیمت خرید موردنظر را وارد کنید (مثال: 20000-10000)")
			return
		})

		t.Bot.Handle(&btnRent, func(c telebot.Context) (err error) {
			min := session.Filters.StartRentPrice
			max := session.Filters.EndRentPrice
			if min != nil || max != nil {
				err = c.Send("بازه قیمت اجاره و ودیعه قبلا مشخص شده است", filterMenu)
				return
			}
			session.State = "setting_rent_price"
			if session.Filters.CategoryPR == nil {
				session.Filters.CategoryPR = new(uint)
			}
			*session.Filters.CategoryPR = 1 // 1 for rent
			err = c.Send("لطفا بازه قیمت اجاره و ودیعه موردنظر را به ترتیب وارد کنید (مثال: 2000000-20000 100000-10000)")
			return
		})

		return c.Send("لطفا دسته بندی مورد نظر خود را انتخاب کنید", purchaseOrRentMenu)
	})

	t.Bot.Handle(&btnCity, func(c telebot.Context) (err error) {
		if city := session.Filters.City; city != nil {
			err = c.Send("شهر قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_city"
		err = c.Send("لطفا نام شهر را وارد کنید")
		return
	})

	t.Bot.Handle(&btnNeighborhood, func(c telebot.Context) (err error) {
		if neighborhood := session.Filters.Neighborhood; neighborhood != nil {
			err = c.Send("محله قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_neighborhood"
		err = c.Send("لطفا نام محله را وارد کنید")
		return
	})

	t.Bot.Handle(&btnArea, func(c telebot.Context) (err error) {
		min := session.Filters.StartArea
		max := session.Filters.EndArea
		if min != nil || max != nil {
			err = c.Send("بازه متراژ قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_area"
		err = c.Send("لطفا بازه متراژ را وارد کنید (مثال: 120-100)")
		return
	})

	t.Bot.Handle(&btnNumberOfRooms, func(c telebot.Context) (err error) {
		min := session.Filters.StartNumberOfRooms
		max := session.Filters.EndNumberOfRooms
		if min != nil || max != nil {
			err = c.Send("بازه تعداد اتاق خواب قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_number_of_rooms"
		err = c.Send("لطفا بازه تعداد اتاق خواب را وارد کنید (مثال: 3-1)")
		return
	})

	t.Bot.Handle(&btnAge, func(c telebot.Context) (err error) {
		min := session.Filters.StartAge
		max := session.Filters.EndAge
		if min != nil || max != nil {
			err = c.Send("بازه سن بنا قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_age"
		err = c.Send("لطفا بازه سن بنا را وارد کنید (مثال: 10-5)")
		return
	})

	t.Bot.Handle(&btnCategoryAV, func(c telebot.Context) (err error) {
		if categoryAV := session.Filters.CategoryAV; categoryAV != nil {
			err = c.Send("این دسته بندی قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_category_AV"
		err = c.Send("ویلایی یا آپارتمانی؟")
		return
	})

	t.Bot.Handle(&btnFloorNumber, func(c telebot.Context) (err error) {
		if categoryAV := session.Filters.CategoryAV; categoryAV == nil {
			err = c.Send("لطفا قبل از تعیین بازه طبقه ویلایی یا آپارتمانی بودن آگهی را تعیین کنید")
			return
		}
		if categoryAV := session.Filters.CategoryAV; *categoryAV == 0 {
			err = c.Send("برای آگهی های ویلایی نمی توان بازه طبقه مشخص کرد")
			return
		}
		min := session.Filters.StartFloorNumber
		max := session.Filters.EndFloorNumber
		if min != nil || max != nil {
			err = c.Send("بازه طبقه قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_floor_number"
		err = c.Send("لطفا بازه طبقه را وارد کنید (مثال: 2-1)")
		return
	})

	t.Bot.Handle(&btnStorage, func(c telebot.Context) (err error) {
		if storage := session.Filters.Storage; storage != nil {
			err = c.Send("داشتن یا نداشتن انباری قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_storage"
		err = c.Send("انباری داتشته باشد یا نه؟ (با بله یا خیر پاسخ دهید)")
		return
	})

	t.Bot.Handle(&btnElevator, func(c telebot.Context) (err error) {
		if elevator := session.Filters.Elevator; elevator != nil {
			err = c.Send("داشتن یا نداشتن آسانسور قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_elevator"
		err = c.Send("آسانسور داشته باشد یا نه؟ (با بله یا خیر پاسخ دهید)")
		return
	})

	t.Bot.Handle(&btnAdDate, func(c telebot.Context) (err error) {
		min := session.Filters.StartDate
		max := session.Filters.EndDate
		if min != nil || max != nil {
			err = c.Send("بازه تاریخ ایجاد آگهی قبلا مشخص شده است", filterMenu)
			return
		}
		session.State = "setting_ad_date"
		err = c.Send("لطفا بازه تاریخ ایجاد آگهی را وارد کنید (مثال: 01-01-1403  01-01-1402)")
		return
	})

	// TODO: insert it into database
	t.Bot.Handle(&btnSendFilter, func(c telebot.Context) (err error) {
		// session := models.GetUserSession(c.Chat().ID)
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
	session := models.GetUserSession(c.Chat().ID)
	getOutputFileMenu.Reply(
		getOutputFileMenu.Row(btnGetAsZip, btnGetViaEmail),
		getOutputFileMenu.Row(btnBackFilterMenu),
	)

	// TODO
	t.Bot.Handle(&btnGetOutputFile, func(c telebot.Context) (err error) {
		return
	})
	// TODO
	t.Bot.Handle(&btnGetViaEmail, func(c telebot.Context) (err error) {
		return
	})
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

// -adminMenu handlers
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
	session := models.GetUserSession(c.Chat().ID)
	input := c.Text()

	switch session.State {
	case "setting_city":
		if session.Filters.City == nil {
			session.Filters.City = new(string)
		}
		*session.Filters.City = strings.TrimSpace(strings.ToLower(input))
		session.State = ""
		err = c.Send("شهر مورد نظر ثبت شد", filterMenu)

	case "setting_neighborhood":
		if session.Filters.Neighborhood == nil {
			session.Filters.Neighborhood = new(string)
		}
		*session.Filters.Neighborhood = strings.TrimSpace(strings.ToLower(input))
		session.State = ""
		err = c.Send("محله مورد نظر ثبت شد", filterMenu)

	case "setting_area":
		minArea, maxArea, e := utils.ParseRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			return
		}
		if session.Filters.StartArea == nil || session.Filters.EndArea == nil {
			session.Filters.StartArea = new(uint)
			session.Filters.EndArea = new(uint)
		}
		*session.Filters.StartArea = minArea
		*session.Filters.EndArea = maxArea
		session.State = ""
		err = c.Send("بازه متراژ مورد نظر ثبت شد", filterMenu)

	case "setting_number_of_rooms":
		minNumberOfRooms, maxNumberOfRooms, e := utils.ParseRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			return
		}
		if session.Filters.StartNumberOfRooms == nil || session.Filters.EndNumberOfRooms == nil {
			session.Filters.StartNumberOfRooms = new(uint)
			session.Filters.EndNumberOfRooms = new(uint)
		}
		*session.Filters.StartNumberOfRooms = minNumberOfRooms
		*session.Filters.EndNumberOfRooms = maxNumberOfRooms
		session.State = ""
		err = c.Send("بازه تعداد اتاق خواب مورد نظر ثبت شد", filterMenu)

	case "setting_age":
		minAge, maxAge, e := utils.ParseRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			return
		}
		if session.Filters.StartAge == nil || session.Filters.EndAge == nil {
			session.Filters.StartAge = new(uint)
			session.Filters.EndAge = new(uint)
		}
		*session.Filters.StartAge = minAge
		*session.Filters.EndAge = maxAge
		session.State = ""
		err = c.Send("بازه سن بنا ثبت شد", filterMenu)

	case "setting_category_AV":
		if input != "آپارتمانی" && input != "ویلایی" {
			err = c.Send("دسته بندی نامعتبر است. لطفا دوباره امتحان کنید")
			return
		}

		if session.Filters.CategoryAV == nil {
			session.Filters.CategoryAV = new(uint)
		}
		if input == "ویلایی" {
			*session.Filters.CategoryAV = 0
		} else {
			*session.Filters.CategoryAV = 1
		}
		session.State = ""
		err = c.Send("دسته بندی مورد نظر ثبت شد", filterMenu)

	case "setting_floor_number":
		minFloorNumber, maxFloorNumber, e := utils.ParseRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			return
		}
		if session.Filters.StartFloorNumber == nil || session.Filters.EndFloorNumber == nil {
			session.Filters.StartFloorNumber = new(uint)
			session.Filters.EndFloorNumber = new(uint)
		}
		*session.Filters.StartFloorNumber = minFloorNumber
		*session.Filters.EndFloorNumber = maxFloorNumber
		session.State = ""
		err = c.Send("بازه تعداد طبقه مورد نظر ثبت شد", filterMenu)

	case "setting_storage":
		if input != "بله" && input != "خیر" {
			err = c.Send("دسته بندی نامعتبر است. دوباره امتحان کنید")
			return
		}

		if session.Filters.Storage == nil {
			session.Filters.Storage = new(bool)
		}
		if input == "بله" {
			*session.Filters.Storage = true
		} else {
			*session.Filters.Storage = false
		}
		session.State = ""
		err = c.Send("ثبت شد", filterMenu)

	case "setting_elevator":
		if input != "بله" && input != "خیر" {
			err = c.Send("دسته بندی نامعتبر است. دوباره امتحان کنید")
			return
		}

		if session.Filters.Elevator == nil {
			session.Filters.Elevator = new(bool)
		}
		if input == "بله" {
			*session.Filters.Elevator = true
		} else {
			*session.Filters.Elevator = false
		}
		session.State = ""
		_ = c.Send("ثبت شد", filterMenu)

	case "setting_ad_date":
		minDate, maxDate, e := utils.ParseDateRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			return
		}
		if session.Filters.StartDate == nil || session.Filters.EndDate == nil {
			session.Filters.StartDate = new(time.Time)
			session.Filters.EndDate = new(time.Time)
		}
		*session.Filters.StartDate = minDate
		*session.Filters.EndDate = maxDate
		session.State = ""
		err = c.Send("بازه تاریخ ثبت آگهی مورد نظر ثبت شد", filterMenu)

	case "setting_purchase_price":
		minPurchasePrice, maxPurchasePrice, e := utils.ParseRanges(input)
		if e != nil {
			err = e
			return
		}
		if session.Filters.StartPurchasePrice == nil || session.Filters.EndPurchasePrice == nil {
			session.Filters.StartPurchasePrice = new(uint)
			session.Filters.EndPurchasePrice = new(uint)
		}
		*session.Filters.StartPurchasePrice = minPurchasePrice
		*session.Filters.EndPurchasePrice = maxPurchasePrice
		session.State = ""
		err = c.Send("بازه قیمت خرید مورد نظر ثبت شد", filterMenu)

	case "setting_rent_price":
		ranges := strings.Split(strings.TrimSpace(input), " ")
		minRentPrice, maxRentPrice, e := utils.ParseRanges(ranges[0])
		if e != nil {
			err = e
			return
		}
		if session.Filters.StartRentPrice == nil || session.Filters.EndRentPrice == nil {
			session.Filters.StartRentPrice = new(uint)
			session.Filters.EndRentPrice = new(uint)
		}
		*session.Filters.StartRentPrice = minRentPrice
		*session.Filters.EndRentPrice = maxRentPrice

		minMortgagePrice, maxMortgagePrice, e := utils.ParseRanges(ranges[1])
		if e != nil {
			err = e
			return
		}
		if session.Filters.StartMortgagePrice == nil || session.Filters.EndMortgagePrice == nil {
			session.Filters.StartMortgagePrice = new(uint)
			session.Filters.EndMortgagePrice = new(uint)
		}
		*session.Filters.StartMortgagePrice = minMortgagePrice
		*session.Filters.EndMortgagePrice = maxMortgagePrice
		session.State = ""
		err = c.Send("بازه قیمت مورد نظر ثبت شد", filterMenu)

	default:
		err = c.Send("لطفا از منو آیتم مورد نظر خود را را انتخاب کنید")
	}
	return
}
