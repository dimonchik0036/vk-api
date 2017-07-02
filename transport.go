package vkapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	defaultHTTPTimeout        = 60 * time.Second
	defaultRequestTimeout     = 120 * time.Second
	defaultKeepAliveInterval  = 60 * time.Second
	defaultHTTPHeadersTimeout = defaultRequestTimeout
)

// Response is response from API server.
type Response struct {
	Errors   *Errors `json:"execute_errors,omitempty"`
	Error    *Error  `json:"error,omitempty"`
	Response Raw     `json:"response,omitempty"`
}

// Request contains data for the request to the API server.
type Request struct {
	Method string     `json:"method"`
	Token  string     `json:"token"`
	Values url.Values `json:"values"`
}

// NewRequest creates a new Request instance.
func NewRequest(method string, token string, values url.Values) (req Request) {
	req.Method = method
	req.Token = token
	req.Values = values
	return
}

// Raw similar to the jsonRAW.
type Raw []byte

func (r Raw) Bytes() []byte {
	return []byte(r)
}

func (r Raw) String() string {
	return bytes.NewBuffer(r).String()
}

func (r Raw) MarshalJSON() ([]byte, error) {
	return r, nil
}

func (r *Raw) UnmarshalJSON(data []byte) error {
	*r = data
	return nil
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Do sends a request to a specific endpoint with our request
// and returns response.
func (api *APIClient) Do(request Request) (response *Response, error *Error) {
	if request.Values == nil {
		request.Values = url.Values{}
	}

	for i, v := range api.Values() {
		if request.Values.Get(i) == "" {
			request.Values.Add(i, v[0])
		}
	}

	req := request.HTTP()

	api.logPrintf("Request: %s", request.JS())

	var res *http.Response
	var err interface {
		Error() string
	}

	for attempt := 1; attempt < 5; attempt++ {
		res, err = api.httpClient.Do(req)
		if err == nil {
			break
		}

		api.logPrintf("HTTP attempt %s %d", err, attempt)

		time.Sleep(time.Second * 3)
	}

	if err != nil {
		api.logPrintf("HTTP fatal %s", err)

		return nil, NewError(ErrBadCode, "HTTP fatal "+err.Error())
	}

	if res.StatusCode != http.StatusOK {
		return nil, NewError(ErrBadResponseCode, res.Status)
	}

	if response, error = Process(res.Body); error != nil {
		api.logPrintf("Response error: %s", error.Error())
	} else {
		api.logPrintf("Response: %s", response.Response.String())
	}

	return
}

// HTTP translates the Request in *http.Request.
func (r Request) HTTP() (req *http.Request) {
	values := r.Values

	if r.Token != "" {
		values.Set(paramToken, r.Token)
	}

	u := ApiURL()
	u.Path = path.Join(u.Path, r.Method)
	u.RawQuery = values.Encode()

	req, err := http.NewRequest(defaultMethod, u.String(), nil)
	if err != nil {
		panic(err)
	}

	return req
}

// JS translates the Request in string.
func (r Request) JS() string {
	args := make(map[string]string)
	for k := range r.Values {
		args[k] = r.Values.Get(k)
	}

	js := new(bytes.Buffer)
	if err := json.NewEncoder(js).Encode(args); err != nil {
		panic(err)
	}

	jsString := js.String()
	jsString = strings.TrimSpace(jsString)

	return fmt.Sprintf("%s(%s)", r.Method, jsString)
}

// vkResponseProcessor stores the Reader.
type vkResponseProcessor struct {
	input io.Reader
}

// To unmarshal the Response data.
func (r Response) To(v interface{}) error {
	return json.Unmarshal(r.Response.Bytes(), v)
}

// ServerError checks the Response to the error.
func (r Response) ServerError() error {
	if r.Errors != nil {
		return r.Errors
	}
	if r.Error != nil {
		return nil
	}
	return r.Error
}

// To unmarshal the Response from vkResponseProcessor data.
func (d vkResponseProcessor) To(response *Response) *Error {
	if rc, ok := d.input.(io.ReadCloser); ok {
		defer rc.Close()
	}

	if err := json.NewDecoder(d.input).Decode(response); err != nil {
		return NewError(ErrBadCode, err.Error())
	}

	return response.Error
}

// Process processes the input data
// and returns a *Response in case of success.
func Process(input io.Reader) (response *Response, err *Error) {
	response = new(Response)
	return response, vkResponseProcessor{input}.To(response)
}

// Encoder is the structure for processing the responses.
type Encoder struct {
	response *Response
	err      error
}

// To unmarshal the Encoder data.
func (e Encoder) To(v interface{}) error {
	if e.err != nil {
		return e.err
	}

	if e.response == nil {
		return errors.New("Unexpected error.")
	}

	return e.response.To(v)
}

// Encode makes Encoder from input data.
func Encode(input io.Reader) Encoder {
	res, err := Process(input)
	return Encoder{res, err}
}

func defaultHTTPClient() HTTPClient {
	client := &http.Client{
		Timeout: defaultRequestTimeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   defaultHTTPTimeout,
				KeepAlive: defaultKeepAliveInterval,
			}).Dial,
			TLSHandshakeTimeout:   defaultHTTPTimeout,
			ResponseHeaderTimeout: defaultHTTPHeadersTimeout,
		},
	}
	return client
}
