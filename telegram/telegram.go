package telegram

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"main/csv"
	"main/database"
	"main/email"
	logg "main/log"
	"main/middlewares"
	"main/models"
	"main/repository"
	"main/utils"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"gopkg.in/telebot.v4"
	"gorm.io/gorm"
)

type TelegramConfig struct {
	Token string
}

type Telegram struct {
	Bot     *telebot.Bot
	Mutex   sync.Mutex
	Config  *TelegramConfig
	Loggers logg.TelegramLogger
}

var (
	superAdminMenu       = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnAddAdmin          = superAdminMenu.Text("اضافه کردن ادمین")
	btnManageAdmins      = superAdminMenu.Text("مدیریت کردن ادمین ها")
	btnSetNumberOfAds    = superAdminMenu.Text("تنظیم تعداد آیتم های جستجو شده")
	btnSetCrawlTimeLimit = superAdminMenu.Text("تنظیم محدودیت زمانی فرآیند جستجو")
	// اطلاعاتت تمامی کاربران به همراه اطلاعات کرال های انجام شده هر کاربر
	// CRUD --> Ehsan

	adminMenu          = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnSeeCrawlDetails = superAdminMenu.Text("دیدن اطلاعات کرال های انجام شده")

	userMenu          = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnBookmarkAd     = userMenu.Text("اضافه کردن آگهی به لیست علاقه مندی ها")
	btnSetFilters     = userMenu.Text("ثبت فیلتر")
	btnShareBookmarks = userMenu.Text("اشتراک گذاری آگهی های مورد علاقه")
	btnGetOutputFile  = userMenu.Text("خروجی گرفتن از آگهی ها")
	btnDeleteHistory  = userMenu.Text("پاک کردن تاریخچه")
	btnAddWatchList   = userMenu.Text("تنظیم کردن watch-list")
	// watchList --> Hirad

	watchListMenu = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnFilterID   = watchListMenu.Text("آیدی فیلتر مورد نظر")
	btnTime       = watchListMenu.Text("بازه زمانی بروزرسانی")
	btnSubmitWL   = watchListMenu.Text("ثبت watchlist")
	btnBackWL     = watchListMenu.Text("بازگشت")

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

	logFile, err := os.Create("log/telegram.log")
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram log file: %w", err)
	}

	infoLogger := log.New(logFile, "INFO: ", log.LstdFlags)
	errorLogger := log.New(logFile, "ERROR: ", log.LstdFlags)

	telegramLogger := logg.TelegramLogger{
		InfoLogger:  infoLogger,
		ErrorLogger: errorLogger,
	}

	telegram := &Telegram{
		Bot:     bot,
		Config:  config,
		Loggers: telegramLogger,
	}

	return telegram, nil
}

type TelegramRciver struct {
	Id string
}

func (reciver *TelegramRciver) Recipient() string {
	return reciver.Id
}

func (t *Telegram) SendMessageToUser(idReciver string, message interface{}) error {
	_, err := t.Bot.Send(&TelegramRciver{
		Id: idReciver,
	}, message)
	return err
}

func (t *Telegram) registerHandlers() {
	t.Bot.Handle("/start", t.handleStart)
	t.Bot.Handle(telebot.OnText, t.handleText)
}

func (t *Telegram) Start() {
	t.registerHandlers()
	t.Loggers.InfoLogger.Println("Starting Telegram bot...")
	t.Bot.Start()
}

func (t *Telegram) handleStart(c telebot.Context) (err error) {
	welcomeMsg := "به ربات خزنده خوش اومدین :)"

	session := models.GetUserSession(c.Chat().ID)
	telegram_ID := c.Sender().ID
	gormUser := repository.NewGormUser(database.GetInstnace().Db)
	user, e := gormUser.GetByTelegramId(strconv.Itoa(int(telegram_ID)))
	if e != nil && errors.Is(e, gorm.ErrRecordNotFound) {
		user, _ = gormUser.Add(models.Users{TelegramId: strconv.Itoa(int(telegram_ID)), Role: "user"})
	}
	if user.Role == "user" {
		userMenu.Reply(
			userMenu.Row(btnSetFilters, btnShareBookmarks, btnBookmarkAd),
			userMenu.Row(btnGetOutputFile, btnDeleteHistory, btnAddWatchList),
		)

		t.Bot.Handle(&btnSetFilters, t.handleSetFilters)
		t.Bot.Handle(&btnShareBookmarks, t.handleShareBookmarks)
		t.Bot.Handle(&btnGetOutputFile, t.handleGetOutput)
		t.Bot.Handle(&btnDeleteHistory, t.handleDeleteHistory)
		t.Bot.Handle(&btnBookmarkAd, t.handleBookmarkAd)
		t.Bot.Handle(&btnAddWatchList, t.handleAddWatchList)

		err = c.Send(welcomeMsg, userMenu)

	} else if user.Role == "admin" {
		adminMenu.Reply(
			adminMenu.Row(btnSeeCrawlDetails),
		)

		t.Bot.Handle(&btnSeeCrawlDetails, t.handleSeeCrawlDetails)

		err = c.Send(welcomeMsg, adminMenu)

	} else if user.Role == "super_admin" {
		if !session.IsAuthenticated {
			err = c.Send("لطقا رمز خود را وارد کنید")
			session.State = "awaiting_password"
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
	}
	return
}

func (t *Telegram) handleAddWatchList(c telebot.Context) (err error) {
	session := models.GetUserSession(c.Chat().ID)
	session.WatchList = models.WatchList{}
	session.State = "adding_watchlist"

	watchListMenu.Reply(
		watchListMenu.Row(btnFilterID, btnTime),
		watchListMenu.Row(btnSubmitWL, btnBackWL),
	)

	t.Bot.Handle(&btnFilterID, func(ctx telebot.Context) (err error) {
		if filterID := session.WatchList.FilterId; filterID != 0 {
			err = c.Send("آیدی فیلتر قبلا مشخص شده", watchListMenu)
			return
		}
		session.State = "watchlist_filterID"
		err = c.Send("لطفا آیدی فیلتر مورد نظر را وارد کنید")
		return
	})

	t.Bot.Handle(&btnTime, func(ctx telebot.Context) (err error) {
		if wlTime := session.WatchList.Time; wlTime != 0 {
			err = c.Send("بازه زمانی قبلا مشخص شده است", watchListMenu)
			return

		}
		session.State = "watchlist_time"
		err = c.Send(" لطفا بازه زمانی بروزرسانی را برحسب دقیقه وارد کنید (e.g: 10)")
		return
	})

	t.Bot.Handle(&btnSubmitWL, func(ctx telebot.Context) (err error) {
		session := models.GetUserSession(c.Chat().ID)
		if session.WatchList.FilterId == 0 || session.WatchList.Time == 0 {
			err = c.Send("لطفا ابتدا تمامی آیتم های خواسته شده را وارد کنید", watchListMenu)
			t.Loggers.InfoLogger.Println("Empty fields")
			return
		}
		gormUser := repository.NewGormUser(database.GetInstnace().Db)
		gormWL := repository.NewWatchList(database.GetInstnace().Db)

		user, _ := gormUser.GetByTelegramId(strconv.Itoa(int(session.ChatID)))
		session.WatchList.UserID = user.ID

		gormWL.Add(session.WatchList)

		return c.Send("شما با موفقیت ثبت شد  watch-list", userMenu)
	})

	t.Bot.Handle(&btnBackWL, func(c telebot.Context) (err error) {
		session.State = ""
		t.Bot.Handle(&btnSetFilters, t.handleSetFilters)
		t.Bot.Handle(&btnShareBookmarks, t.handleShareBookmarks)
		t.Bot.Handle(&btnBookmarkAd, t.handleBookmarkAd)
		t.Bot.Handle(&btnGetOutputFile, t.handleGetOutput)
		t.Bot.Handle(&btnDeleteHistory, t.handleDeleteHistory)
		t.Bot.Handle(&btnAddWatchList, t.handleAddWatchList)

		err = c.Send("به صفحه اصلی بازگشتید", userMenu)
		return
	})

	err = c.Send("لطفا یک گزینه را انتخاب کنید ", watchListMenu)
	return
}

// -userMenu handlers
func (t *Telegram) handleSetFilters(c telebot.Context) (err error) {
	session := models.GetUserSession(c.Chat().ID)
	session.Filters = models.Filters{}
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

	t.Bot.Handle(&btnSendFilter, func(c telebot.Context) (err error) {
		session := models.GetUserSession(c.Chat().ID)
		gormFilter := repository.NewGormFilter(database.GetInstnace().Db)
		gormUser := repository.NewGormUser(database.GetInstnace().Db)
		u, _ := gormUser.GetByTelegramId(strconv.Itoa(int(session.ChatID)))
		gormUserAd := repository.NewGormUser_Ad(database.GetInstnace().Db)

		gormFilter.Add(session.Filters)
		gormAd := repository.NewGormAd(database.GetInstnace().Db)
		ads, e := gormAd.Get(session.Filters)
		if e == nil {
			for _, ad := range ads {
				gormUserAd.Add(models.Users_Ads{UserId: u.ID,
					AdId: ad.ID})
				t.Loggers.ErrorLogger.Println(ad)
			}
		}

		return c.Send(fmt.Sprintf("فیلتر شما با موفیقت ثبت شد. آیدی فیلتر شما %d می باشد", session.Filters.ID), userMenu)

	})

	t.Bot.Handle(&btnBackFilterMenu, func(c telebot.Context) (err error) {
		session.State = ""
		t.Bot.Handle(&btnSetFilters, t.handleSetFilters)
		t.Bot.Handle(&btnShareBookmarks, t.handleShareBookmarks)
		t.Bot.Handle(&btnBookmarkAd, t.handleBookmarkAd)
		t.Bot.Handle(&btnGetOutputFile, t.handleGetOutput)
		t.Bot.Handle(&btnDeleteHistory, t.handleDeleteHistory)
		t.Bot.Handle(&btnAddWatchList, t.handleAddWatchList)

		err = c.Send("به صفحه اصلی بازگشتید", userMenu)
		return
	})

	err = c.Send("لطفا یک فیلتر را انتخاب کنید", filterMenu)
	return
}

func (t *Telegram) handleShareBookmarks(c telebot.Context) (err error) {
	session := models.GetUserSession(c.Chat().ID)
	session.State = "sharing_bookmarks"
	err = c.Send("لطفا آیدی کاربر مورد نظر را وارد کنید")
	return
}

func (t *Telegram) handleGetOutput(c telebot.Context) error {
	session := models.GetUserSession(c.Chat().ID)
	getOutputFileMenu.Reply(
		getOutputFileMenu.Row(btnGetAsZip, btnGetViaEmail),
		getOutputFileMenu.Row(btnBackFilterMenu),
	)

	t.Bot.Handle(&btnGetAsZip, func(c telebot.Context) (err error) {
		session.State = "sending_zip_file"
		filename, e := csv.ExportCsv(strconv.Itoa(int(session.ChatID)), database.GetInstnace().Db)

		if e != nil {
			t.Loggers.ErrorLogger.Println("exporting csv failed")
		}
		uniqueZip := fmt.Sprintf("./telegram/%s.zip", uuid.New())
		adsZip, e := os.Create(uniqueZip)

		if e != nil {
			t.Loggers.InfoLogger.Println("Creating zip file failed")
		}

		zipWriter := zip.NewWriter(adsZip)

		f, e := os.Open(filename)
		if e != nil {
			t.Loggers.InfoLogger.Println("Opening csv file failed")
		}

		z, e := zipWriter.Create("Advertisments/ads.csv")

		if e != nil {
			t.Loggers.InfoLogger.Println("Creating directory in zip file failed")
		}
		if _, e := io.Copy(z, f); e != nil {
			t.Loggers.InfoLogger.Println("Couldn't copy .csv file into zip file")
		}
		zipWriter.Close()

		e = t.SendMessageToUser(strconv.Itoa(int(session.ChatID)), &telebot.Document{File: telebot.FromDisk(uniqueZip),
			FileName:             uniqueZip,
			DisableTypeDetection: true})

		if e != nil {
			t.Loggers.ErrorLogger.Println("sending file failed", e)
			c.Send("مشکلی پیش آمد! مجددا تلاش کنید")
			os.Remove(uniqueZip)
			os.Remove(filename)
			adsZip.Close()
			f.Close()
			return
		}

		adsZip.Close()
		f.Close()
		os.Remove(filename)
		os.Remove(uniqueZip)

		session.State = ""
		t.Loggers.InfoLogger.Println("Zip output sent")
		return
	})

	t.Bot.Handle(&btnGetViaEmail, func(c telebot.Context) (err error) {
		session.State = "setting_user_email"
		err = c.Send("لطفا ایمیل خود را وارد کنید")
		return
	})

	t.Bot.Handle(&btnBackGetOutputFileMenu, func(c telebot.Context) (err error) {
		session.State = ""
		t.Bot.Handle(&btnSetFilters, t.handleSetFilters)
		t.Bot.Handle(&btnShareBookmarks, t.handleShareBookmarks)
		t.Bot.Handle(&btnBookmarkAd, t.handleBookmarkAd)
		t.Bot.Handle(&btnGetOutputFile, t.handleGetOutput)
		t.Bot.Handle(&btnDeleteHistory, t.handleDeleteHistory)
		t.Bot.Handle(&btnAddWatchList, t.handleAddWatchList)

		err = c.Send("به صفحه اصلی بازگشتید", userMenu)
		return
	})
	return c.Send("نحوه خروجی گرفتن را انتخاب کنید", getOutputFileMenu)
}

func (t *Telegram) handleDeleteHistory(c telebot.Context) (err error) {
	telegram_ID := c.Sender().ID
	userRepo := repository.NewGormUser(database.GetInstnace().Db)
	userAdsRepo := repository.NewGormUser_Ad(database.GetInstnace().Db)
	user, e := userRepo.GetByTelegramId(strconv.Itoa(int(telegram_ID)))
	if e != nil {
		err = e
		return
	}
	userAds, e := userAdsRepo.GetByUserId([]uint{user.ID})
	if e != nil {
		err = e
		return
	}
	for _, item := range userAds {
		e = userAdsRepo.Delete(item.User.ID, item.AdId)
		if e != nil {
			err = e
			return
		}
	}

	err = c.Send("تاریخچه شما حذف شد", userMenu)
	t.Loggers.InfoLogger.Println("user history deleted")
	return
}

func (t *Telegram) handleBookmarkAd(c telebot.Context) (err error) {
	session := models.GetUserSession(c.Chat().ID)
	session.State = "adding_bookmark"
	err = c.Send("لطفا آیدی آگهی مورد علاقه را وارد کنید")
	return
}

// -adminMenu handlers
// TODO
func (t *Telegram) handleSeeCrawlDetails(c telebot.Context) (err error) {
	return nil
}

// superAdminMenu handlers
func (t *Telegram) handleAddAdmin(c telebot.Context) (err error) {
	session := models.GetUserSession(c.Chat().ID)
	session.State = "adding_admin"
	err = c.Send("لطفا آیدی تلگرام کاربر مورد نظر خود را وارد کنید")
	return
}

func (t *Telegram) handleManageAdmins(c telebot.Context) (err error) {
	session := models.GetUserSession(c.Chat().ID)
	session.State = "managing_admin"
	err = c.Send("لطفا آیدی تلگرام ادمین مورد نظر را وارد کنید")
	return
}

func (t *Telegram) handleSetCrawlTimeLimit(c telebot.Context) (err error) {
	session := models.GetUserSession(c.Chat().ID)
	session.State = "setting_crawl_time_limit"
	err = c.Send("لطفا یک عدد به دقیقه برای تنظیم محدودیت زمانی فرایند جستجو کرالر وارد کنید")
	return
}

func (t *Telegram) handleSetNumberOfAds(c telebot.Context) (err error) {
	session := models.GetUserSession(c.Chat().ID)
	session.State = "setting_max_searched_items"
	err = c.Send("لطفا حداکثر تعداد آیتم های جستجو شده در هر کرال را وارد کنید")
	return
}

func (t *Telegram) handleText(c telebot.Context) (err error) {
	session := models.GetUserSession(c.Chat().ID)
	input := c.Text()
	lowerInput := strings.ToLower(strings.TrimSpace(input))

	switch session.State {
	case "awaiting_password":
		e, _ := middlewares.Authentication(input)
		if e != nil {
			err = c.Send("رمز ورود اشتباه است دوباره امتحان کنید")
			t.Loggers.ErrorLogger.Println("error in authenticating super admin")
			return
		} else {
			session.IsAuthenticated = true
			session.State = ""
			t.Loggers.InfoLogger.Println("super admin logged in successfully")

			superAdminMenu.Reply(
				superAdminMenu.Row(btnAddAdmin, btnManageAdmins),
				superAdminMenu.Row(btnSetCrawlTimeLimit, btnSetNumberOfAds),
			)

			t.Bot.Handle(&btnAddAdmin, t.handleAddAdmin)
			t.Bot.Handle(&btnManageAdmins, t.handleManageAdmins)
			t.Bot.Handle(&btnSetCrawlTimeLimit, t.handleSetCrawlTimeLimit)
			t.Bot.Handle(&btnSetNumberOfAds, t.handleSetNumberOfAds)

			err = c.Send("به ربات خزنده خوش آمدید :)", superAdminMenu)
		}

	case "watchlist_filterID":
		filterID, e := strconv.Atoi(input)
		if e != nil || filterID <= 0 {
			c.Send("آیدی وارد شده نامعتبر است")
			t.Loggers.InfoLogger.Println("Invalid filterID")
			return
		}
		session.WatchList.FilterId = uint(filterID)
		session.State = ""
		err = c.Send("آیدی مورد نظر ثبت شد", watchListMenu)
		t.Loggers.InfoLogger.Println("Set watchlist filterID")

	case "watchlist_time":
		wlTime, e := strconv.Atoi(input)
		if e != nil || wlTime <= 5 {
			c.Send("زمان وارد شده نامعتبر است! حداقل 5 دقیقه")
			t.Loggers.InfoLogger.Println("Invalid time duration")
			return
		}
		session.WatchList.Time = wlTime
		session.State = ""
		err = c.Send("بازه زمانی مورد نظر ثبت شد", watchListMenu)
		t.Loggers.InfoLogger.Println("Set time duration for watch-list")

	case "setting_city":
		if session.Filters.City == nil {
			session.Filters.City = new(string)
		}
		*session.Filters.City = strings.TrimSpace(strings.ToLower(input))
		session.State = ""
		err = c.Send("شهر مورد نظر ثبت شد", filterMenu)
		t.Loggers.InfoLogger.Println("set city")

	case "setting_neighborhood":
		if session.Filters.Neighborhood == nil {
			session.Filters.Neighborhood = new(string)
		}
		*session.Filters.Neighborhood = strings.TrimSpace(strings.ToLower(input))
		session.State = ""
		err = c.Send("محله مورد نظر ثبت شد", filterMenu)
		t.Loggers.InfoLogger.Println("set neighborhood")

	case "setting_area":
		minArea, maxArea, e := utils.ParseRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			t.Loggers.ErrorLogger.Println("parsing area range")
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
		t.Loggers.InfoLogger.Println("area range set")

	case "setting_number_of_rooms":
		minNumberOfRooms, maxNumberOfRooms, e := utils.ParseRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			t.Loggers.ErrorLogger.Println("parsing number of rooms range")
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
		t.Loggers.InfoLogger.Println("number of rooms range set")

	case "setting_age":
		minAge, maxAge, e := utils.ParseRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			t.Loggers.ErrorLogger.Println("parsing building age range")
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
		t.Loggers.InfoLogger.Println("building age range set")

	case "setting_category_AV":
		if input != "آپارتمانی" && input != "ویلایی" {
			err = c.Send("دسته بندی نامعتبر است. لطفا دوباره امتحان کنید")
			t.Loggers.ErrorLogger.Println("invalid input for AV-category")
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
		t.Loggers.InfoLogger.Println("AV-category set")

	case "setting_floor_number":
		minFloorNumber, maxFloorNumber, e := utils.ParseRanges(input)
		if e != nil {
			err = c.Send(e.Error())
			t.Loggers.ErrorLogger.Println("parsing floor number range")
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
		t.Loggers.InfoLogger.Println("floor number range set")

	case "setting_storage":
		if input != "بله" && input != "خیر" {
			err = c.Send("دسته بندی نامعتبر است. دوباره امتحان کنید")
			t.Loggers.InfoLogger.Println("invalid input for storage exsistence")
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
		t.Loggers.InfoLogger.Println("storage exsistence set")

	case "setting_elevator":
		if input != "بله" && input != "خیر" {
			err = c.Send("دسته بندی نامعتبر است. دوباره امتحان کنید")
			t.Loggers.ErrorLogger.Println("invalid input for elevator exsistence")
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
		err = c.Send("ثبت شد", filterMenu)
		t.Loggers.InfoLogger.Println("elevator exsistence set")

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
		t.Loggers.InfoLogger.Println("ad date range set")

	case "setting_purchase_price":
		minPurchasePrice, maxPurchasePrice, e := utils.ParseRanges(input)
		if e != nil {
			err = e
			t.Loggers.ErrorLogger.Println("parsing purchase price range")
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
		t.Loggers.InfoLogger.Println("purchase price range set")

	case "setting_rent_price":
		ranges := strings.Split(strings.TrimSpace(input), " ")
		minRentPrice, maxRentPrice, e := utils.ParseRanges(ranges[0])
		if e != nil {
			err = e
			t.Loggers.ErrorLogger.Println("parsing rent price range")
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
			t.Loggers.ErrorLogger.Println("parsing mortgage price range")
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
		t.Loggers.InfoLogger.Println("rent and mortgage price range set")

	case "setting_user_email":
		emailRegex := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
		if !emailRegex.Match([]byte(input)) {
			err = c.Send("ایمیل نامعتبر است دوباره امتحان کنید")
			t.Loggers.InfoLogger.Println("invalid email input")
			return
		}
		session.State = ""
		session.Email = input

		filename, e := csv.ExportCsv(strconv.Itoa(int(session.ChatID)), database.GetInstnace().Db)

		if e != nil {
			t.Loggers.ErrorLogger.Println("exporting csv failed")
		}

		err = email.SendEmail(session.Email, filename)
		if err != nil {
			t.Loggers.ErrorLogger.Println("sending email failed")
		}
		os.Remove(filename)

		session.State = ""
		err = c.Send("ایمیل شما ثبت شد", userMenu)
		t.Loggers.InfoLogger.Println("user's email set")
		return

	case "adding_bookmark":
		userAd := repository.NewGormUser_Ad(database.GetInstnace().Db)
		adId, e := strconv.Atoi(input)
		if e != nil {
			err = c.Send("آیدی آگهی نامعتبر است دوباره امتحان کنید")
			t.Loggers.InfoLogger.Println("Invalid Ad ID")
			return
		}
		e = userAd.Update(models.Users_Ads{
			UserId:     uint(session.ChatID),
			AdId:       uint(adId),
			IsBookmark: true,
		})

		if e != nil {
			t.Loggers.ErrorLogger.Println("adding bookmark failed: ", e)
			err = c.Send("خطایی در افزودن آگهی به علاقه مندی ها رخ داد. مجددا امتحان کنید")
			return
		}

		session.State = ""
		err = c.Send("آگهی به لیست علاقه مندی ها اضافه شد")
		return

	case "sharing_bookmarks":
		re := regexp.MustCompile(`^[0-9]+$`)
		if !re.MatchString(input) {
			t.Loggers.InfoLogger.Println("Invalid user ID!")
			err = c.Send("آیدی کاربر نامعتبر است! مجددا امتحان کنید")
			return
		}

		userID := input

		userAd := repository.NewGormUser_Ad(database.GetInstnace().Db)
		ads, e := userAd.GetByUserId([]uint{uint(session.ChatID)})
		if e != nil {
			t.Loggers.ErrorLogger.Println("getting user's ads failed: ", e)
			return
		}
		links := ""
		for _, ad := range ads {
			if ad.IsBookmark {
				links += fmt.Sprintf("%s\n", ad.Ad.Link)
			}
		}
		e = t.SendMessageToUser(userID, links)
		if e != nil {
			t.Loggers.ErrorLogger.Println("sharing bookmarks failed: ", e)
			c.Send("آیدی کاربر نامعتبر است! مجددا تلاش کنید")
			return
		}

		session.State = ""
		err = c.Send("آگهی ها به کاربر ارسال شد")
		return

	case "adding_admin":
		gormUser := repository.NewGormUser(database.GetInstnace().Db)
		user, e := gormUser.GetByTelegramId(lowerInput)
		if e != nil && errors.Is(e, gorm.ErrRecordNotFound) {
			user, _ = gormUser.Add(models.Users{TelegramId: lowerInput, Role: "admin"})
			err = c.Send("نقش کاربر مورد نظر به ادمین تغییر یافت", superAdminMenu)
			t.Loggers.InfoLogger.Println("user's role changed to admin")
			session.State = ""
			return
		}
		user.Role = "admin"
		e = gormUser.Update(user)
		if e != nil {
			err = e
			t.Loggers.ErrorLogger.Println("error updating user's role")
			return
		}
		err = c.Send("نقش کاربر مورد نظر به ادمین تغییر یافت", superAdminMenu)
		t.Loggers.InfoLogger.Println("user's role changed to admin")
		session.State = ""
		return

	case "managing_admin":
		gormUser := repository.NewGormUser(database.GetInstnace().Db)
		user, e := gormUser.GetByTelegramId(lowerInput)
		if e != nil && errors.Is(e, gorm.ErrRecordNotFound) {
			err = c.Send("کاربری با آیدی مورد نظر یافت نشد لطفا دوباره تلاش کنید")
			t.Loggers.InfoLogger.Println("user not found")
			return
		}
		user.Role = "user"
		e = gormUser.Update(user)
		if e != nil {
			err = e
			t.Loggers.ErrorLogger.Println("error updating user's role")
			return
		}
		err = c.Send("نقش ادمین مورد نظر به کاربر ساده تغییر یافت", superAdminMenu)
		t.Loggers.InfoLogger.Println("admin's role changed to user")
		session.State = ""
		return

	case "setting_crawl_time_limit":
		re := regexp.MustCompile(`^[0-9]+$`)
		if !re.MatchString(input) {
			err = c.Send("لطفا زمان مورد نظر خود را به دقیقه وارد کنید")
			t.Loggers.InfoLogger.Println("invalid user input for crawl time limit")
			return
		}
		e := os.Setenv("TIMEOUT", input)
		if e != nil {
			err = e
			t.Loggers.ErrorLogger.Println("error in setting crawl time limit")
			return
		}
		err = c.Send("محدودیت زمان فرآیند جستجو آپدیت شد", superAdminMenu)
		t.Loggers.InfoLogger.Println("crawl time limit updated")
		session.State = ""
		return

	case "setting_max_searched_items":
		re := regexp.MustCompile(`^[0-9]+$`)
		if !re.MatchString(input) {
			err = c.Send("لطفا حداکثر تعداد آیتم های جستجو شده در هر کرال را به صورت یک عدد وارد کنید")
			t.Loggers.InfoLogger.Println("invalid user input for max searched items in each crawl")
			return
		}
		e := os.Setenv("MAX_SEARCHED_ITEMS", input)
		if e != nil {
			err = e
			t.Loggers.ErrorLogger.Println("error in setting max searched items in each crawlW")
			return
		}
		err = c.Send("تعداد آیتم های جستجو شده آپدیت شد", superAdminMenu)
		t.Loggers.InfoLogger.Println("crawl max searched items updated")
		session.State = ""
		return

	default:
		err = c.Send("لطفا از منو آیتم مورد نظر خود را را انتخاب کنید")
	}
	return
}
