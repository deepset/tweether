package logs

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {

	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file : ", err)
	}

	InfoLogger = log.New(file, "Info : ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "Error : ", log.Ldate|log.Ltime|log.Lshortfile)
}
