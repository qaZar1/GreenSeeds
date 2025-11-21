package device

import "time"

const (
	delimeter    = 0x04
	delimeterStr = "\x04"
	WAIT         = "WAIT"
	RETURN       = "RETURN"
	BUSY         = "BUSY"
	READY        = "READY"
	STAND_BY     = "STAND BY"
	ERR          = "ERR"
	END          = "END"
	BEGIN        = "BEGIN"
	ACK_BOOT     = "ACK BOOT"
	BOOT         = "BOOT"
	STATUS       = "STATUS"

	emptyLetter = "\x00"

	SETSTATUS_READY   = "SETSTATUS READY"
	ACK_SETSTATUS_READY = "ACK SETSTATUS READY"
	DEVICE_NOT_ACTIVE = "Device is not active"

	ackTimeout = 20 * time.Second
)
