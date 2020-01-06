package log

import "log"

func Init() {
	log.SetPrefix("[G] ")
	log.SetFlags(log.Ldate | log.Ltime)
}
