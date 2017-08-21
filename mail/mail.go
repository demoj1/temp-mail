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
)

var Alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

type TempMail struct {
	login      string
	domain     string
	api_domain string
}

// NewTempMail - ...
// Если поле login == "" будет сгенерирован случайный логин длинной 10 символов.
// Если поле domain == "" при получение нового почтового адреса домен будет выбран случайно.
func NewTempMail(login, domain string) *TempMail {
	if login == "" {
		login = GenerateLogin(10)
	}

	rand.Seed(time.Now().UnixNano())

	return &TempMail{
		login:      login,
		domain:     domain,
		api_domain: "api.temp-mail.ru"}
}

func (t *TempMail) AvailableDomains() ([]string, error) {
	url := fmt.Sprintf("http://%v/request/domains/format/json/", t.api_domain)
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

func GenerateLogin(n int) string {
	login := make([]rune, n)
	for i := range login {
		login[i] = choiceRune(Alphabet)
	}

	return string(login)
}

func (t *TempMail) GetEmailAddress() (string, error) {
	available_domains, err := t.AvailableDomains()
	if err != nil {
		return "", err
	}

	if t.domain == "" {
		t.domain = choiceString(available_domains)
	}

	if !inString(t.domain, available_domains) {
		return "", errors.New(
			fmt.Sprintf(
				"Domain %v not in available domains.\nPlease choices one from %v",
				t.domain, available_domains))
	}

	return t.login + t.domain, nil
}

func GetMD5Hash(email string) string {
	bytes := md5.Sum([]byte(email))
	return hex.EncodeToString(bytes[:])
}

func (t *TempMail) GetMailBox(email, emailHash string) (map[string]string, error) {
	if email == "" {
		var err error
		email, err = t.GetEmailAddress()
		if err != nil {
			return nil, err
		}
	}

	if emailHash == "" {
		emailHash = GetMD5Hash(email)
	}

	url := fmt.Sprintf("http://%v/request/mail/id/%v/format/json/", t.api_domain, emailHash)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var messages map[string]string
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &messages)
	if err != nil {
		return nil, err
	}

	if error, ok := messages["error"]; ok {
		return nil, errors.New(error)
	}

	return messages, nil
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
