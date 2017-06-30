package vkapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

type Response struct {
	Errors   Errors `json:"execute_errors,omitempty"`
	Error    Error  `json:"error,omitempty"`
	Response Raw    `json:"response,omitempty"`
}

type Request struct {
	Method string     `json:"method"`
	Token  string     `json:"token"`
	Values url.Values `json:"values"`
}

type Raw []byte

func (r Raw) Bytes() []byte {
	return []byte(r)
}

func (r Raw) String() string {
	return bytes.NewBuffer(r).String()
}

func (m Raw) MarshalJSON() ([]byte, error) {
	log.Println("marshalling to", m)
	return m, nil
}

func (m *Raw) UnmarshalJSON(data []byte) error {
	*m = data
	return nil
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (api *ApiClient) Do(request Request) (response *Response, error *Error) {
	for i, v := range api.Values() {
		if request.Values.Get(i) == "" {
			request.Values.Add(i, v[0])
		}
	}

	req := request.HTTP()
	start := time.Now()
	log.Println("DO", request.Method)

	var res *http.Response
	var err interface {
		Error() string
	}

	for attempt := 1; attempt < 5; attempt++ {
		res, err = api.httpClient.Do(req)
		if err == nil {
			break
		}

		log.Println("HTTP attempt", err, attempt)
		time.Sleep(time.Second * 3)
	}

	if err != nil {
		log.Println("HTTP fatal", err)
		return nil, NewError(ErrBadCode, "HTTP fatal "+err.Error())
	}

	log.Println("HTTP", res.Status, time.Now().Sub(start))
	if res.StatusCode != http.StatusOK {
		return nil, NewError(ErrBadResponseCode, res.Status)
	}

	return Process(res.Body)
}

func (r Request) HTTP() (req *http.Request) {
	values := r.Values

	if r.Token != "" {
		values.Set(paramToken, r.Token)
	}

	u := ApiUrl()
	u.Path = path.Join(u.Path, r.Method)
	u.RawQuery = values.Encode()

	log.Print("URL: ", u.String())
	req, err := http.NewRequest(defaultMethod, u.String(), nil)

	must(err)

	return req
}

func (r Request) JS() string {
	args := make(map[string]string)
	for k := range r.Values {
		args[k] = r.Values.Get(k)
	}
	js := new(bytes.Buffer)
	encoder := json.NewEncoder(js)

	must(encoder.Encode(args))

	jsString := js.String()
	jsString = strings.TrimSpace(jsString)

	return fmt.Sprintf("API.%s(%s)", r.Method, jsString)
}

type vkResponseProcessor struct {
	input io.Reader
}

func (r Response) To(v interface{}) error {
	return json.Unmarshal(r.Response.Bytes(), v)
}

func (r Response) ServerError() error {
	if r.Errors != nil {
		return r.Errors
	}
	if r.Error.Code == ErrZero {
		return nil
	}
	return r.Error
}

func (d vkResponseProcessor) To(response *Response) *Error {
	if rc, ok := d.input.(io.ReadCloser); ok {
		defer rc.Close()
	}

	decoder := json.NewDecoder(d.input)

	if err := decoder.Decode(response); err != nil {
		return NewError(ErrBadCode, err.Error())
	}

	return &response.Error
}

func Process(input io.Reader) (response *Response, err *Error) {
	response = new(Response)
	return response, vkResponseProcessor{input}.To(response)
}

type Encoder struct {
	response *Response
	err      error
}

func (e Encoder) To(v interface{}) error {
	if e.err != nil {
		return e.err
	}
	if e.response == nil {
		return errors.New("wtf")
	}
	return e.response.To(v)
}

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
