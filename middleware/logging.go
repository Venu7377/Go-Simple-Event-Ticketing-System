package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logFilePath := "./logs/app.log"
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("___New Log_____")
	start := time.Now()
	log.Printf("[Method : %s] %s\n", r.Method, r.URL.Path)

	l.handler.ServeHTTP(w, r)
	log.Printf("Duration : %v", time.Since(start))
}

func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}
