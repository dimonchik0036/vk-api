package vkapi

import (
	"log"
)

func ExampleClient_SendMessage() {
	client, err := NewClientFromLogin("<username>", "<password>", ScopeMessages)
	if err != nil {
		log.Panic(err)
	}

	client.Log(true)

	if err := client.InitMyProfile(); err != nil {
		log.Panic(err)
	}

	// Sends a message to himself.
	if _, err := client.SendMessage(NewMessage(client.User.Id, "Hello!")); err != nil {
		log.Println(err)
	}
}
