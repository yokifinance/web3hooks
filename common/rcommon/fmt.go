package rcommon

import (
	"fmt"
	"log"
	"time"
)

var includeDatetime int

func init() {
	includeDatetime = GetParamIntOrDefault("DATETIME_IN_CONSOLE", 1)
	if includeDatetime == 1 {
		log.SetFlags(log.Flags() | log.LstdFlags) // set flag
	} else {
		log.SetFlags(log.Flags() & (log.LstdFlags ^ -1)) // clear flag
	}
}

func Println(msg string, a ...any) {
	msg = fmt.Sprintf(msg, a...)
	if includeDatetime == 1 {
		msg = time.Now().Format("2006/01/02 15:04:05.000000") + " " + msg
	}
	fmt.Println(msg)
}

// returns true if we need to print more detailed log messages (usually for debug)
func IsExtendedLog() bool {
	return GetParamIntOrDefault("EXTENDED_LOGS", 0) != 0
}
