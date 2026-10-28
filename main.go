package main

import (
	"bunker-core/protocol/mpay/android"
	"bunker-core/protocol/mpay/windows"
	"bunker-web/configs"
	"bunker-web/routers"
	"fmt"
	"os"

	"runtime/debug"
	"sync"
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

func warmPublishedPackages() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := (&android.AndroidMpayUser{}).UsePublishedPackage(); err != nil {
			fmt.Println("android package:", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := (&windows.WindowsMpayUser{}).UsePublishedPackage(); err != nil {
			fmt.Println("windows package:", err)
		}
	}()
	wg.Wait()
}

// @title			BunkerWeb OpenAPI
// @description	BunkerWeb OpenAPI document, feel free to contact us if you have any questions or suggestions.
func main() {
	defer handlePanic()

	// Warm the packages without holding the server back. A request that needs
	// one loads it on use, waiting for this warm-up if it is still running, so
	// a slow CDN only delays those requests instead of keeping the API offline.
	fmt.Println("Loading published packages in the background...")
	go func() {
		warmPublishedPackages()
		fmt.Println("Published package warm-up finished")
	}()

	router := routers.InitRouter()

	fmt.Println("Server starts running...")

	router.Run(fmt.Sprintf(":%d", configs.HTTP_PORT))
}
