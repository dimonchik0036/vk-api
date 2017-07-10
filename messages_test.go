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
	if _, err := client.SendMessage(NewMessage(NewDstFromUserID(client.User.ID), "Hello!")); err != nil {
		log.Println(err)
	}
}
