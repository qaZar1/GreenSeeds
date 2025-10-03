package service

import "errors"

var (
	ErrInvalidBody             = errors.New("Invalid body")
	ErrInvalidUnmarshal        = errors.New("Invalid unmarshal")
	ErrValidateStruct          = errors.New("Invalid validate struct")
	ErrInvalidGeneratePassHash = errors.New("Invalid generate password hash")
	ErrUserAlreadyExists       = errors.New("User already exists")
	ErrUserNotFound            = errors.New("User not found")
	ErrInvalidGetDataFromDB    = errors.New("Invalid get data from DB")

	ErrInvalidUsernameOrPassword = errors.New("Invalid username or password")
	ErrInvalidAddUser            = errors.New("Invalid add user")
	ErrInvalidGenerateToken      = errors.New("Invalid generate token")
	ErrInvalidUpdateData         = errors.New("Invalid update data")

	ErrInvalidGetMessages      = errors.New("Invalid get messages")
	ErrInvalidRelocateMessages = errors.New("Invalid relocate messages")
	ErrInvalidRemoveMessage    = errors.New("Invalid remove message")
	ErrInvalidResendMessage    = errors.New("Invalid resend message")

	ErrChannelAlreadyExists = errors.New("Channel already exists")
	ErrChannelNotFound      = errors.New("Channel not found")
	ErrInvalidAddChannel    = errors.New("Invalid add channel")
	ErrInvalidUpdateChannel = errors.New("Invalid update channel")
	ErrInvalidRemoveChannel = errors.New("Invalid remove channel")

	ErrSystemAlreadyExists = errors.New("System already exists")
	ErrSystemNotFound      = errors.New("System not found")
	ErrInvalidAddSystem    = errors.New("Invalid add system")
	ErrInvalidUpdateSystem = errors.New("Invalid update system")
	ErrInvalidRemoveSystem = errors.New("Invalid remove system")

	ErrInvalidRemoveResend   = errors.New("Invalid remove resend")
	ErrInvalidUpdateResend   = errors.New("Invalid update resend")
	ErrInvalidAddResend      = errors.New("Invalid add resend")
	ErrInvalidResendNotFound = errors.New("Resend not found")

	ErrRouteAlreadyExists = errors.New("Route already exists")
	ErrRouteNotFound      = errors.New("Route not found")
	ErrInvalidAddRoute    = errors.New("Invalid add route")
	ErrInvalidUpdateRoute = errors.New("Invalid update route")
	ErrInvalidRemoveRoute = errors.New("Invalid remove route")

	ErrInvalidRemoveBinding   = errors.New("Invalid remove binding")
	ErrInvalidUpdateBinding   = errors.New("Invalid update binding")
	ErrInvalidAddBinding      = errors.New("Invalid add binding")
	ErrInvalidBindingNotFound = errors.New("Binding not found")
	ErrInvalidSyncBindings    = errors.New("Invalid sync bindings")
	ErrNotFound               = errors.New("Not found")

	ErrMappingAlreadyExists = errors.New("Mapping already exists")
	ErrMappingNotFound      = errors.New("Mapping not found")
	ErrInvalidSyncMappings  = errors.New("Invalid sync mappings")
	ErrInvalidAddMapping    = errors.New("Invalid add mapping")
	ErrInvalidUpdateMapping = errors.New("Invalid update mapping")
	ErrInvalidRemoveMapping = errors.New("Invalid remove mapping")

	ErrInvalidAddLicense    = errors.New("Invalid add license")
	ErrInvalidUpdateLicense = errors.New("Invalid update license")
	ErrInvalidRemoveLicense = errors.New("Invalid remove license")
	ErrLicenseAlreadyExists = errors.New("License already exists")
	ErrLicenseNotFound      = errors.New("License not found")

	ErrInvalidInterval = errors.New("Invalid interval")
)
