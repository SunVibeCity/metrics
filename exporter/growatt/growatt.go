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

type PlantList struct {
	Back struct {
		Data []struct {
			PlantMoneyText string `json:"plantMoneyText"`
			PlantName      string `json:"plantName"`
			PlantID        string `json:"plantId"`
			IsHaveStorage  string `json:"isHaveStorage"`
			TodayEnergy    string `json:"todayEnergy"`
			TotalEnergy    string `json:"totalEnergy"`
			CurrentPower   string `json:"currentPower"`
		} `json:"data"`
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
		fmt.Println(string(bs))
		for k, v := range resp.Header {
			fmt.Print(k)
			fmt.Print(" : ")
			fmt.Println(v)
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