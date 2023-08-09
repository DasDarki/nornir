package log

import "log"

var isDebug = true

func Debug(msg string) {
	if isDebug {
		log.Println(msg)
	}
}

func Debugf(format string, v ...interface{}) {
	if isDebug {
		log.Printf(format, v...)
	}
}
