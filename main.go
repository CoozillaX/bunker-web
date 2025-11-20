package main

import (
	"bunker-web/configs"
	"bunker-web/routers"
	"fmt"
	"os"

	"runtime/debug"
	"time"
)

func handlePanic() {
	if r := recover(); r != nil {
		// open file
		file, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		defer file.Close()
		// get stack trace
		stackTrace := debug.Stack()
		// write error log
		logMessage := fmt.Sprintf("[%s] panic: %v\n\n%s\n", time.Now().Format(time.RFC3339), r, stackTrace)
		file.WriteString(logMessage)
		// print error log
		fmt.Println("\n" + logMessage)
	}
}

// @title			BunkerWeb OpenAPI
// @description	BunkerWeb OpenAPI document, feel free to contact us if you have any questions or suggestions.
func main() {
	defer handlePanic()

	router := routers.InitRouter()

	fmt.Println("Server starts running...")

	router.Run(fmt.Sprintf(":%d", configs.HTTP_PORT))
}
