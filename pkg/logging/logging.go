package logging

import (
	"log"
	"os"
)

var (
	CommonLog *log.Logger
	ErrorLog  *log.Logger
)

func Init() *log.Logger {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)

	err := os.MkdirAll("logs", 0755)

	if err != nil || os.IsExist(err) {
		panic("can't create log dir. no configured logging to files")
	} else {
		openLogfile, err := os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		CommonLog = log.New(openLogfile, "Common Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
		ErrorLog = log.New(openLogfile, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
		return log.New(openLogfile, "Default Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	}

}
