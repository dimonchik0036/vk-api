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
	if _, err := client.SendMessage(NewMessage(NewDstFromUserID(client.VKUser.Me.ID), "Hello!")); err != nil {
		log.Println(err)
	}
}

func ExampleClient_SendMessage2() {
	client, err := NewClientFromLogin("<username>", "<password>", ScopeMessages)
	if err != nil {
		log.Fatal(err)
	}

	client.Log(true)
	if err := client.InitMyProfile(); err != nil {
		log.Fatal(err)
	}

	m := NewMessage(NewDstFromUserID(client.VKUser.Me.ID), "Hello!")
	m.Attachment = client.AddAttachmentPhoto("photo.jpg")
	if _, err := client.SendMessage(m); err != nil {
		log.Println(err)
	}

}
