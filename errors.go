package vkapi

import (
	"fmt"
	"strings"
)

type ServerError int

const (
	// Full description at https://vk.com/dev/errors
	ErrZero                      ServerError = 0
	ErrUnknown                   ServerError = 1
	ErrApplicationDisabled       ServerError = 2
	ErrUnknownMethod             ServerError = 3
	ErrInvalidSignature          ServerError = 4
	ErrAuthFailed                ServerError = 5
	ErrTooManyRequests           ServerError = 6
	ErrInsufficientPermissions   ServerError = 7
	ErrInvalidRequest            ServerError = 8
	ErrTooManyOneTypeRequests    ServerError = 9
	ErrInternalServerError       ServerError = 10
	ErrAppInTestMode             ServerError = 11
	ErrCaptchaNeeded             ServerError = 14
	ErrNotAllowed                ServerError = 15
	ErrHttpsOnly                 ServerError = 16
	ErrNeedValidation            ServerError = 17
	ErrUserDeletedOrBlocked      ServerError = 18
	ErrStandaloneOnly            ServerError = 20
	ErrStandaloneOpenAPIOnly     ServerError = 21
	ErrMethodDisabled            ServerError = 23
	ErrNeedConfirmation          ServerError = 24
	ErrCommunityKeyInvalid       ServerError = 27
	ErrApplicationKeyInvalid     ServerError = 28
	ErrOneOfParametersInvalid    ServerError = 100
	ErrInvalidAPIID              ServerError = 101
	ErrInvalidAUserID            ServerError = 113
	ErrInvalidTimestamp          ServerError = 150
	ErrAlbumAccessProhibited     ServerError = 200
	ErrAudioAccessProhibited     ServerError = 201
	ErrGroupAccessProhibited     ServerError = 203
	ErrAlbumOverflow             ServerError = 300
	ErrEnableVoiceApplication    ServerError = 500
	ErrInsufficientPermissionsAd ServerError = 600
	ErrInternalServerErrorAd     ServerError = 603

	ErrInBlackList           ServerError = 900
	ErrNotAllowedToSendFirst ServerError = 901
	ErrPrivacy               ServerError = 902

	ErrBadResponseCode ServerError = -1
	ErrBadCode         ServerError = -666
)

type Errors []ExecuteError

func (e Errors) Error() string {
	var s []string
	for _, v := range e {
		s = append(s, v.Error())
	}
	return fmt.Sprintln("Execute errors:", strings.Join(s, ", "))
}

type RequestParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ExecuteError struct {
	Method  string      `json:"method"`
	Code    ServerError `json:"error_code"`
	Message string      `json:"error_msg"`
}

// Error contains standard errors.
type Error struct {
	Code       ServerError     `json:"error_code,omitempty"`
	Message    string          `json:"error_msg,omitempty"`
	Params     *[]RequestParam `json:"request_params,omitempty"`
	CaptchaSid string          `json:"captcha_sid,omitempty"`
	CaptchaImg string          `json:"captcha_img,omitempty"`
	Request    Request         `json:"-"`
}

// NewError makes *Error from our ServerError and description.
func NewError(code ServerError, description string) (err *Error) {
	err = new(Error)
	err.Code = code
	err.Message = description

	return
}

// setRequest sets Request
func (e *Error) setRequest(r Request) {
	e.Request = r
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

func (e ExecuteError) Error() string {
	return fmt.Sprintf("%s: %s (%d)", e.Method, e.Message, e.Code)
}

//Is returns true if this is an error.
func (e ServerError) Is(err error) bool {
	if error(e) == err {
		return true
	}
	if another, ok := err.(ServerError); ok {
		return another == e
	}
	if another, ok := err.(Error); ok {
		return another.Code == e
	}
	return false
}

func (s ServerError) String() string {
	return fmt.Sprintf("%d", s)
}

func (s ServerError) Error() string {
	return s.String()
}
