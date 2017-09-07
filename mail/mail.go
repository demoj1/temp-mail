package mail

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"net/http/cookiejar"

	"net/url"
)

var Alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

const (
	API_DOMAIN  = "api.temp-mail.ru"
	MAIN_PAGE   = "https://temp-mail.ru/"
	CHANGE_MAIL = "https://temp-mail.ru/option/change"
)

type TempMail struct {
	Login     string
	Domain    string
	EmailHash string
}

type Message struct {
	MailFrom      string  `json:"mail_from"`
	MailSubject   string  `json:"mail_subject"`
	MailText      string  `json:"mail_text"`
	MailTimestamp float64 `json:"mail_timestamp"`
}

func AvailableDomains() ([]string, error) {
	url := fmt.Sprintf("http://%v/request/domains/format/json/", API_DOMAIN)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	domainsJson, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var domains []string
	err = json.Unmarshal(domainsJson, &domains)
	if err != nil {
		return nil, err
	}

	return domains, nil
}

func GetRandomEmail() (*TempMail, error) {
	available_domains, err := AvailableDomains()
	if err != nil {
		return nil, err
	}

	domain := choiceString(available_domains)
	login := generateLogin(10)

	err = createEmail(login, domain)
	if err != nil {
		return nil, err
	}

	return newTempMail(login, domain), nil
}

func GetEmail(login, domain string) (*TempMail, error) {
	available_domains, err := AvailableDomains()
	if err != nil {
		return nil, err
	}

	domain = "@p33.org"

	if !inString(domain, available_domains) {
		return nil, errors.New(
			fmt.Sprintf(
				"домен %v не может быть выбран.\nПожалуйста выберите один из следующих %v",
				domain, available_domains))
	}

	err = createEmail(login, domain)
	if err != nil {
		return nil, err
	}

	return newTempMail(login, domain), nil
}

func GetMessages(emailHash string) ([]Message, error) {
	url := fmt.Sprintf("http://%s/request/mail/id/%s/format/json/", API_DOMAIN, emailHash)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	var messages []Message
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func newTempMail(login, domain string) *TempMail {
	rand.Seed(time.Now().UnixNano())

	return &TempMail{
		Login:     login,
		Domain:    domain,
		EmailHash: getMD5Hash(login + domain),
	}
}

func generateLogin(n int) string {
	login := make([]rune, n)
	for i := range login {
		login[i] = choiceRune(Alphabet)
	}

	return string(login)
}

func createEmail(login, domain string) error {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	_, err := client.Get(MAIN_PAGE)
	if err != nil {
		return err
	}

	urlMainPage, _ := url.Parse(MAIN_PAGE)
	cookies := jar.Cookies(urlMainPage)
	csrf := getCookieByName(cookies, "csrf")

	_, err = client.PostForm(CHANGE_MAIL, map[string][]string{
		"csrf":   {csrf},
		"mail":   {login},
		"domain": {domain},
	})

	if err != nil {
		return err
	}

	return nil
}

func getCookieByName(cookie []*http.Cookie, name string) string {
	cookieLen := len(cookie)
	result := ""
	for i := 0; i < cookieLen; i++ {
		if cookie[i].Name == name {
			result = cookie[i].Value
		}
	}
	return result
}

func getMD5Hash(email string) string {
	bytes := md5.Sum([]byte(email))
	return hex.EncodeToString(bytes[:])
}

func choiceRune(array []rune) rune {
	return array[rand.Intn(len(array))]
}

func choiceString(array []string) string {
	return array[rand.Intn(len(array))]
}

func inString(item string, array []string) bool {
	for _, e := range array {
		if item == e {
			return true
		}
	}

	return false
}
