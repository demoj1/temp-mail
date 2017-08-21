package main

import (
	"log"
	"temp_mail/mail"
	"time"
)

func main() {
	tm := mail.NewTempMail("dmitryd.prog", "")
	email, err := tm.GetEmailAddress()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(email)

	time.Sleep(1 * time.Second)

	log.Println(tm.GetMailBox("", ""))
}
