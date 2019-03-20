package growatt

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type httpClient struct {
	sess http.Client
	url string
	user User
	username string
	token string
}

type LoginResponse struct {
	User User `json:"back"`
}
type User struct {
	UserID    int  `json:"userId"`
	UserLevel int  `json:"userLevel"`
	Success   bool `json:"success"`
}

type PlantListResponse struct {
	Back struct {
		PlantList [] Plant `json:"data"`
		TotalData struct {
			CurrentPowerSum string `json:"currentPowerSum"`
			CO2Sum          string `json:"CO2Sum"`
			IsHaveStorage   string `json:"isHaveStorage"`
			ETotalMoneyText string `json:"eTotalMoneyText"`
			TodayEnergySum  string `json:"todayEnergySum"`
			TotalEnergySum  string `json:"totalEnergySum"`
		} `json:"totalData"`
		Success bool `json:"success"`
	} `json:"back"`
}

type Plant struct {
	PlantMoneyText string `json:"plantMoneyText"`
	PlantName      string `json:"plantName"`
	PlantID        string `json:"plantId"`
	IsHaveStorage  string `json:"isHaveStorage"`
	TodayEnergy    string `json:"todayEnergy"`
	TotalEnergy    string `json:"totalEnergy"`
	CurrentPower   string `json:"currentPower"`
}

type Timespan int
const (
	Day 	Timespan = 1
	Month 	Timespan = 2
	Year 	Timespan = 3
	Total 	Timespan = 4
)

func New(username string, password string, url string) *httpClient {
	jar, _ := cookiejar.New(nil)
	sess := http.Client{
		Jar:jar,
	}
	return &httpClient{
		sess:sess,
		url:url,
		username:username,
		token:hashPassword(password),
		user:User{},
	}
}

func hashPassword(p string) string {
	h := md5.New()
	_, err := io.WriteString(h, p)
	if err != nil {
		panic(err)
	}
	ohbs := []byte(h.Sum(nil))
	for i, b := range ohbs {
		if b < 0x10 {
			ohbs[i] = b + 0xc0
		}
	}
	return fmt.Sprintf("%x", ohbs)
}

func (c *httpClient) Login() bool {
	if (*c).user.Success != true {
		resp, err := c.sess.PostForm(c.url + "LoginAPI.do", url.Values{"userName": {c.username}, "password": {c.token}})
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err.Error())
		}
		lr := LoginResponse{}
		err = json.Unmarshal(bs, &lr)
		if err != nil {
			panic(err.Error())
		}
		(*c).user = lr.User
	}
	return (*c).user.Success == true
}

func (c *httpClient) Logout() bool {
	if (*c).user.Success == true {
		resp, err := c.sess.Get(c.url + "logout.do")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		(*c).user.Success = false
	}
	return (*c).user.Success != true
}

func (c *httpClient) PlantList() []Plant {
	if (*c).Login() != true {
		panic("User can't log in.")
	}
	resp, err := c.sess.Get(c.url + "PlantListAPI.do")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	plr := PlantListResponse{}
	err = json.Unmarshal(bs, &plr)
	if err != nil {
		panic(err.Error())
	}
	return plr.Back.PlantList
}

// Return amount of power generated for the given timespan.
// timeSpan=1, date=%Y-%m-%d (2019-03-20): power on each five minutes of the day.
// timeSpan=2, date=%Y-%m : power on each day of the month.
// timeSpan=3, date=%Y : power on each month of the year.
// timeSpan=4, date= : power on each year. `date` parameter is ignored.

func (c *httpClient) PlantDetail(plantId int64, timeSpan Timespan, date string) string {
	if (*c).Login() != true {
		panic("User can't log in.")
	}
	resp, err := c.sess.Get(fmt.Sprintf("%sPlantDetailAPI.do?plantId=%d&type=%d&date=%s", c.url, plantId, timeSpan, date))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return string(bs)
}