package vkapi

import (
	"fmt"
	"strings"
)

const (
	// Full description at https://vk.com/dev/errors
	ErrZero = iota
	ErrUnknown
	ErrApplicationDisabled
	ErrUnknownMethod
	ErrInvalidSignature
	ErrAuthFailed
	ErrTooManyRequests
	ErrInsufficientPermissions
	ErrInvalidRequest
	ErrTooManyOneTypeRequests
	ErrInternalServerError
	ErrAppInTestMode
	ErrCaptchaNeeded             = 14
	ErrNotAllowed                = 15
	ErrHttpsOnly                 = 16
	ErrNeedValidation            = 17
	ErrUserDeletedOrBlocked      = 18
	ErrStandaloneOnly            = 20
	ErrStandaloneOpenAPIOnly     = 21
	ErrMethodDisabled            = 23
	ErrNeedConfirmation          = 24
	ErrCommunityKeyInvalid       = 27
	ErrApplicationKeyInvalid     = 28
	ErrOneOfParametersInvalid    = 100
	ErrInvalidAPIID              = 101
	ErrInvalidAUserID            = 113
	ErrInvalidTimestamp          = 150
	ErrAlbumAccessProhibited     = 200
	ErrAudioAccessProhibited     = 201
	ErrGroupAccessProhibited     = 203
	ErrAlbumOverflow             = 300
	ErrEnableVoiceApplication    = 500
	ErrInsufficientPermissionsAd = 600
	ErrInternalServerErrorAd     = 603

	ErrBadResponseCode = -1
)

type Errors []ExecuteError

func (e Errors) Error() string {
	var s []string
	for _, v := range e {
		s = append(s, v.Error())
	}
	return fmt.Sprintln("Execute errors:", strings.Join(s, ", "))
}

type ServerError int

type RequestParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ExecuteError struct {
	Method  string      `json:"method"`
	Code    ServerError `json:"error_code"`
	Message string      `json:"error_msg"`
}

type Error struct {
	Code    ServerError    `json:"error_code,omitempty"`
	Message string         `json:"error_msg,omitempty"`
	Params  []RequestParam `json:"request_params,omitempty"`
	Request Request        `json:"-"`
}

func (e *Error) setRequest(r Request) {
	e.Request = r
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

func (e ExecuteError) Error() string {
	return fmt.Sprintf("%s: %s (%d)", e.Method, e.Message, e.Code)
}

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

func IsServerError(err error) bool {
	if _, ok := err.(Error); ok {
		return true
	}
	return false
}

func GetServerError(err error) Error {
	if s, ok := err.(Error); ok {
		return s
	}
	panic("not a server error")
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

func (s ServerError) String() string {
	return string(s)
}

func (s ServerError) Error() string {
	return s.String()
}
