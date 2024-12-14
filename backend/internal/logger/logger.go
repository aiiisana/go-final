package logging

import (
	"log"
	"os"
)

func InitLogger() (*os.File, error) {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)

	return file, nil
}

func LogInfo(message string) {
	log.Printf("[INFO] %s\n", message)
}

func LogError(message string, err error) {
	log.Printf("[ERROR] %s: %v\n", message, err)
}

func LogWarning(message string) {
	log.Printf("[WARN] %s\n", message)
}
