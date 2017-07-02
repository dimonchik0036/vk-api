package vkapi

import (
	"log"
	"time"
)

func ExampleClient_GetLPUpdatesChan() {
	client, err := NewClientFromToken("<access_token>")
	//client, err := NewClientFromLogin("<username>", "<password>", ScopeMessages)
	if err != nil {
		log.Panic(err)
	}

	if err := client.InitLongPoll(0, 2); err != nil {
		log.Panic(err)
	}

	updates, off, err := client.GetLPUpdatesChan(100, LPConfig{25, 0})
	if err != nil {
		log.Panic(err)
	}

	go func() {
		time.Sleep(time.Minute)
		*off = false
	}()

	for update := range updates {
		log.Print("Code: ", update.Code)
		if update.Message != nil {
			log.Print(update.Message)
		}

		if update.FriendNotification != nil {
			log.Print(update.FriendNotification)
		}
	}
}
