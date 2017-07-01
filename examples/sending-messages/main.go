package main

import (
	"github.com/dimonchik0036/vk-api"
	"log"
)

func main() {
	client, err := vkapi.NewClientFromLogin("<username>", "<password>", vkapi.ScopeMessages)
	if err != nil {
		log.Fatal(err)
	}

	client.Log(true)

	if err := client.InitMyProfile([]string{}); err != nil {
		log.Fatal(err)
	}

	if _, err := client.SendMessage(vkapi.NewMessage(client.User.Id, "Hello!")); err != nil {
		log.Println(err)
	}
}
