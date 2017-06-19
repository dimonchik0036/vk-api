package vkapi

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
