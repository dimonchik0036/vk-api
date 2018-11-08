package vkapi

func ExampleClient_PostWall() {
	client, _ := NewClientFromToken("<access token>")
	client.Log(true)

	id, err := client.PostWall(PostConfig{
		Message: "test",
	})

	if err != nil {
		panic(err.Error())
	}
	println(id)
}
