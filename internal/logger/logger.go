package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var infoLogger *log.Logger
var errorLogger *log.Logger

func Init() {
	os.MkdirAll("logs", os.ModePerm)
	filename := fmt.Sprintf("logs/gateway_%s.log", time.Now().Format("20060102"))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error abriendo archivo de log:", err)
	}

	infoLogger = log.New(file, "[INFO] ", log.Ldate|log.Ltime)
	errorLogger = log.New(file, "[ERROR] ", log.Ldate|log.Ltime)
}

func Info(msg string) {
	fmt.Println("[INFO]", msg)
	if infoLogger != nil {
		infoLogger.Println(msg)
	}
}

func Error(msg string) {
	fmt.Println("[ERROR]", msg)
	if errorLogger != nil {
		errorLogger.Println(msg)
	}
}

func Request(method, path, clientIP string, status int, latency time.Duration) {
	msg := fmt.Sprintf("%s %s | IP: %s | Status: %d | Latency: %s",
		method, path, clientIP, status, latency)
	Info(msg)
}
