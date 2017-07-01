package vkapi

type Client struct {
	apiClient *ApiClient
	User      *Users
	LongPoll  *LongPoll
	//	Group
	//	Wall
	//	WallComment
	//  Message
	//  Chat
	//	Note
	//	Page
	//	Board
	//	BoardComment
}

func DefaultClientFromToken(token string) (client *Client, err error) {
	client = new(Client)
	client.apiClient = DefaultApiClient()
	client.apiClient.SetAccessToken(token)
	return
}

func DefaultClientFromLogin(username string, password string, scope int64) (client *Client, err error) {
	client = new(Client)
	client.apiClient = DefaultApiClient()
	err = client.apiClient.Authenticate(DefaultApplication(username, password, scope))
	if err != nil {
		return nil, err
	}

	return
}

func (client *Client) Do(request Request) (response *Response, err *Error) {
	if client.apiClient == nil {
		return nil, NewError(ErrBadCode, "ApiClient not found")
	}

	if request.Token == "" && client.apiClient.AccessToken != nil {
		request.Token = client.apiClient.AccessToken.AccessToken
	}

	return client.apiClient.Do(request)
}
