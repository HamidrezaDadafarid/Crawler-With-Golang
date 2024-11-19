package tests

import (
	"main/crawler"
	"main/models"
	"sync"
	"testing"

	"github.com/go-rod/rod"
	"github.com/stretchr/testify/suite"
)

// DivarTestSuite creates a TestSuite for the Divar crawler
type DivarTestSuite struct {
	suite.Suite
	divar     *crawler.CrawlerAbstract
	browser   *rod.Browser
	waitGroup *sync.WaitGroup
}

// SetupTest initializes the TestSuite before each test
func (suite *DivarTestSuite) SetupTest() {
	suite.browser = rod.New().MustConnect()
	suite.waitGroup = &sync.WaitGroup{}
	suite.divar = crawler.NewDivarCrawler(suite.waitGroup, suite.browser, &crawler.Settings{Page: 10, Items: 30, Timeout: 200}, &models.Metrics{})
}

// TearDownTest cleans up after each test
func (suite *DivarTestSuite) TearDownTest() {
	suite.browser.MustClose()
}

// TestGetTargets tests the GetTargets method
func (suite *DivarTestSuite) TestGetTargets() {
	ads := suite.divar.Crawler.GetTargets(suite.divar.Settings.Page, suite.browser)
	suite.NotEmpty(ads)
	for _, ad := range ads {
		suite.NotEmpty(ad.UniqueId)
		suite.Equal("divar", ad.Link)
	}
}

// TestGetDetails tests the GetDetails method
func (suite *DivarTestSuite) TestGetDetails() {
	ad := &models.Ads{
		UniqueId: "/wZnI0lho",
		Link:     "divar",
	}
	suite.waitGroup.Add(1)
	suite.divar.Crawler.GetDetails(ad, suite.browser, suite.waitGroup)
	suite.Equal(`۵۳متر*پارکینگ*غرق نور*وام دار/فاز ۱`, ad.Title)
	suite.Equal(uint(1), ad.CategoryAV)
	suite.Empty(ad.City)
	suite.Empty(ad.Neighborhood)
	suite.Equal(uint(0), ad.CategoryPR)
	suite.Equal(uint(53), ad.Meters)
	suite.Equal(uint(1393), ad.Age)
	suite.Equal(uint(1), ad.NumberOfRooms)
	suite.Equal(35.677170839626, ad.Latitude)
	suite.Equal(51.022097727083, ad.Longitude)
	suite.Equal(uint(1100000000), ad.SellPrice)
	suite.True(ad.Storage)
	suite.True(ad.Elevator)
	suite.True(ad.Parking)
	suite.Equal(`https://s100.divarcdn.com/static/photo/neda/post/ZP4D1gWtW_IAfxuC-hKQqQ/a6c503aa-4367-45f3-8eee-71a68347db8f.jpg`, ad.PictureLink)
	suite.Equal(uint(0), ad.RentPrice)
	suite.Equal(uint(0), ad.MortgagePrice)

}

// TestDivarSuite runs the entire test suite
func TestDivarSuite(t *testing.T) {
	suite.Run(t, new(DivarTestSuite))
}
