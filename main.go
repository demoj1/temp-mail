package main

import (
	"net/http"

	"fmt"

	"temp_mail/mail"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/available_domains", func(c *gin.Context) {
		domains, err := mail.AvailableDomains()
		if checkError(err, c) {
			return
		}
		c.JSON(http.StatusOK, domains)
	})

	r.POST("/create_mail", func(c *gin.Context) {
		email := c.DefaultPostForm("email", "")
		if checkEmptyField(email, "email", c) {
			return
		}
		domain := c.DefaultPostForm("domain", "")
		if checkEmptyField(domain, "domain", c) {
			return
		}

		tm, err := mail.GetEmail(email, domain)
		if checkError(err, c) {
			return
		}

		c.JSON(http.StatusOK, tm)
	})

	r.POST("messages", func(c *gin.Context) {
		email_hash := c.DefaultPostForm("email_hash", "")
		if checkEmptyField(email_hash, "email_hash", c) {
			return
		}

		msgs, _ := mail.GetMessages(email_hash)

		c.JSON(http.StatusOK, msgs)
	})
	r.Run()
}

func checkError(err error, c *gin.Context) bool {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	return err != nil
}

func checkEmptyField(value, name string, c *gin.Context) bool {
	if value == "" {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Поле %s обязательно и не может быть пустым.", name),
			},
		)
	}

	return value == ""
}
