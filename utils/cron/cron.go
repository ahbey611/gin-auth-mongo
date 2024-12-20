package cron

import (
	"gin-auth-mongo/utils/jwkmanager"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func StartCron() {
	loc, _ := time.LoadLocation("Asia/Kuala_Lumpur")
	option := cron.WithLocation(loc)
	c := cron.New(option)

	// every sunday
	_, err := c.AddFunc("0 0 * * 0", func() {
		log.Println("Updating keys")
		jwkmanager.UpdateKeys()
		log.Println("Keys updated")
	})

	if err != nil {
		panic(err)
	}

	c.Start()
	log.Println("Cron job started")
	log.Println(time.Now().Format("2024-01-01 00:00:00"))
	log.Println(c.Location())
}
