<div align="center">

[![](https://github.com/Dimonchik0036/vk-api/blob/master/logo.png)]()  

# VK API? GO!
[![Build Status](https://travis-ci.org/Dimonchik0036/vk-api.svg?branch=master)](https://travis-ci.org/Dimonchik0036/vk-api)
[![GoDoc](https://godoc.org/github.com/Dimonchik0036/vk-api?status.svg)](https://godoc.org/github.com/Dimonchik0036/vk-api)
[![Language](https://img.shields.io/badge/language-Go-blue.svg)](https://github.com/golang/go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Dimonchik0036/vk-api/blob/master/LICENSE)  
  
Work with Vkontakte API for StandAlone application on The Go Programming Language.  

# Usage / Installation
</div>

## Installation
`go get -u github.com/Dimonchik0036/vk-api`

## Example  
Displays incoming messages. If this is a "/start", then a "Hello!" message will be sent.
```go
package main

import (
	"github.com/dimonchik0036/vk-api"
	"log"
)

func main() {
	//client, err := vkapi.NewClientFromLogin("<username>", "<password>", vkapi.ScopeMessages)
	client, err := vkapi.NewClientFromToken("<access_token>")
	if err != nil {
	    log.Panic(err)
	}
	
	client.Log(true)

	if err := client.InitLongPoll(0, 2); err != nil {
		log.Panic(err)
	}

	updates, _, err := client.GetLPUpdatesChan(100, vkapi.LPConfig{25, vkapi.LPModeAttachments})
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil || !update.IsNewMessage() || update.Message.Outbox(){
			continue
		}

		log.Printf("%s", update.Message.String())
		if update.Message.Text == "/start" {
			client.SendMessage(vkapi.NewMessage(vkapi.NewDstFromUserID(update.Message.FromID), "Hello!"))
		}

	}
}
```
## Technical Details 
* API version 5.67.

## Contributions
Chat me [VK](https://vk.com/dimonchik0036)/[Telegram](https://t.me/dimonchik0036) for detailed steps.
