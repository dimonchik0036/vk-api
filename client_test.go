package vkapi

import "log"

func ExampleNewClientFromToken() {
	client, _ := NewClientFromToken("<access_token>")
	if err := client.InitLongPoll(0, 2); err != nil {
		log.Panic(err)
	}
	client.Log(true)
	updates, _, err := client.GetLPUpdatesChan(100, LPConfig{25, LPModeAttachments})
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("%s", update.Message.String())
		if update.Message.Text == "/start" {
			client.SendMessage(NewMessage(NewDstFromUserID(update.Message.FromID), "Hello!"))
		}

	}
}
