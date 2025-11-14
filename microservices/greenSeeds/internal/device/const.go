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

	SETSTATUS_READY = "SETSTATUS READY"

	ackTimeout = 20 * time.Second
)
