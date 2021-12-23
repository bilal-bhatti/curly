package curly

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

var Debug = false

func init() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(0)
	// log.SetPrefix("")
	log.SetOutput(os.Stderr)
}

func Tracef(msg string, p ...interface{}) {
	if !Debug {
		return
	}

	pc, _, line, _ := runtime.Caller(1)

	caller := fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), line)
	caller = caller[strings.LastIndex(caller, ".")+1:]
	if len(p) == 0 {
		log.Printf("%14s - %s", caller, msg)
	} else {
		log.Printf("%14s - %s", caller, fmt.Sprintf(msg, p...))
	}
}
