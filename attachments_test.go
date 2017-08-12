package vkapi

import (
	"io/ioutil"
	"os"
)

func ExampleClient_SendDoc() {
	var id int64
	client, _ := NewClientFromToken("<access token>")
	client.Log(true)
	//#1
	client.SendDoc(NewDstFromUserID(id), "logo", "logo.png")
	//#2
	file, _ := os.Open("logo.png")
	client.SendDoc(NewDstFromUserID(id), "logo", FileReader{Reader: file, Size: -1, Name: file.Name()})
	//#3
	data, _ := ioutil.ReadFile("logo.png")
	client.SendDoc(NewDstFromUserID(id), "logo", FileBytes{Bytes: data, Name: "logo.png"})
}
