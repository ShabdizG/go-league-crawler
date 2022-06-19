package logging

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func InitLogger(filename string) {
	// Create the log file if doesn't exist. And append to it if it already exists.
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
		// TimestampFormat: "2006-02-02 15:04:06",
	})
	mw := io.MultiWriter(os.Stdout, f)
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	} else {
		log.SetOutput(mw)
	}

}
