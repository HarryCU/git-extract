package log

import (
	extract "github.com/HarryCU/git-extract"
	"log"
)

func Debug(format string, value ...interface{}) {
	if extract.Config.Debug {
		log.Printf(format, value...)
	}
}
