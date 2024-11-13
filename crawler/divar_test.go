package tests

import (
	"main/models"
	"sync"
	"testing"

	"github.com/go-rod/rod"
	"github.com/stretchr/testify/suite"
)

// DivarTestSuite creates a TestSuite for the Divar crawler
type DivarTestSuite struct {
	suite.Suite
	divar     *divar
	browser   *rod.Browser
	waitGroup *sync.WaitGroup
}

// SetupTest initializes the TestSuite before each test
func (suite *DivarTestSuite) SetupTest() {
	suite.divar = &divar{}
	suite.browser = rod.New().MustConnect()
	suite.waitGroup = &sync.WaitGroup{}
}

// TearDownTest cleans up after each test
func (suite *DivarTestSuite) TearDownTest() {
	suite.browser.MustClose()
}

// TestGetUniqueID tests the getUniqueID function
func (suite *DivarTestSuite) TestGetUniqueID() {
	testCases := []struct {
		input    string
		expected string
	}{
		{"https://divar.ir/v/%DB%B5%DB%B3%D9%85%D8%AA%D8%B1-%D9%BE%D8%A7%D8%B1%DA%A9%DB%8C%D9%86%DA%AF-%D8%BA%D8%B1%D9%82-%D9%86%D9%88%D8%B1-%D9%88%D8%A7%D9%85-%D8%AF%D8%A7%D8%B1-%D9%81%D8%A7%D8%B2-%DB%B1/wZnI0lho", "/wZnI0lho"},
		{"https://divar.ir/v/advertisment-title/some_uniqueID_321", "/some_uniqueID_321"},
		{"https://divar.ir/v/", ""},
	}

	for _, tc := range testCases {
		suite.Equal(tc.expected, getUniqueID(tc.input))
	}
}

// TestGetTargets tests the GetTargets method
func (suite *DivarTestSuite) TestGetTargets() {
	ads := suite.divar.GetTargets(1, suite.browser)
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
	suite.divar.GetDetails(ad, suite.browser, suite.waitGroup)
	suite.Equal(`۵۳متر*پارکینگ*غرق نور*وام دار/فاز ۱`, ad.Title)
	suite.Equal(`apartment`, ad.CategoryAV)
	suite.Empty(ad.City)
	suite.Empty(ad.Mahale)
	suite.Equal(`sale`, ad.CategoryPR)
	suite.Equal(53, ad.Meters)
	suite.Equal(1393, ad.Age)
	suite.Equal(1, ad.NumberOfRooms)
	suite.Equal(35.677170839626, ad.Latitude)
	suite.Equal(51.022097727083, ad.Longitude)
	suite.Equal(1100000000, ad.SellPrice)
	suite.True(ad.Anbary)
	suite.True(ad.Elevator)
	suite.True(ad.Parking)
	suite.Equal(`https://s100.divarcdn.com/static/photo/neda/post/ZP4D1gWtW_IAfxuC-hKQqQ/a6c503aa-4367-45f3-8eee-71a68347db8f.jpg`, ad.PictureLink)
	suite.Equal(-1, ad.RentPrice)
	suite.Equal(-1, ad.MortgagePrice)

}

// TestChangeFarsiToEng tests the changeFarsiToEng function
func (suite *DivarTestSuite) TestChangeFarsiToEng() {
	testCases := []struct {
		input    string
		expected int
	}{
		{"۱۲۳۴۵", 12345},
		{"۰", 0},
		{"۹۹۹", 999},
		{"invalid", -1},
	}

	for _, tc := range testCases {
		suite.Equal(tc.expected, changeFarsiToEng(tc.input))
	}
}

// TestGetLatitudeAndLongitude tests the getLatitudeAndLongitude function
func (suite *DivarTestSuite) TestGetLatitudeAndLongitude() {
	testCases := []struct {
		input          string
		expectedLat    float64
		expectedLong   float64
		expectedErrMsg string
	}{
		{"https://balad.ir/location?latitude=35.677170839626&longitude=51.022097727083", 35.677170839626, 51.022097727083, ""},
	}

	for _, tc := range testCases {
		lat, long := GetLatitudeAndLongitude(tc.input)
		suite.Equal(tc.expectedLat, lat)
		suite.Equal(tc.expectedLong, long)
	}
}

// TestDivarSuite runs the entire test suite
func TestDivarSuite(t *testing.T) {
	suite.Run(t, new(DivarTestSuite))
}
