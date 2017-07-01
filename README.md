# Golang vk-api
[![Build Status](https://travis-ci.org/Dimonchik0036/vk-api.svg?branch=master)](https://travis-ci.org/Dimonchik0036/vk-api)
[![GoDoc](https://godoc.org/github.com/Dimonchik0036/vk-api?status.svg)](https://godoc.org/github.com/Dimonchik0036/vk-api)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Dimonchik0036/vk-api/new/master/LICENSE)  
  
Golang package for VKontakte API.  

Current VK API version 5.65.

# Installation
`go get -u github.com/Dimonchik0036/vk-api`

# Example  
Displays incoming messages. If this is a "/start", then a "Hello!" message will be sent.
```go
package main

import (
	"github.com/dimonchik0036/vk-api"
	"log"
)

func main() {
	client, _ := vkapi.NewClientFromToken("<access_token>")
	/*client, err := vkapi.NewClientFromLogin("<username>", "<password>", vkapi.ScopeMessages)
	if err != nil {
	    log.Panic(err)
	}*/
	
	if err := client.InitLongPoll(0, 2); err != nil {
		log.Panic(err)
	}

	updates, _, err := client.GetLPUpdatesChan(100, vkapi.LPConfig{25, vkapi.LPModeAttachments})
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("%d writes:[%s]", update.Message.FromID, update.Message.Text)
		if update.Message.Text == "/start" {
			client.SendMessage(vkapi.NewMessage(update.Message.FromID, "Hello!"))
		}

	}
}
```
