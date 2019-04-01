package main

import (
	"flag"
	"os"

	"./pkg/process"
)

func main() {

	// go run main.go -cmd=ls -args="-la /tmp"
	// go run main.go -cmd=ping -args="-c 10 127.0.0.1"
	cmdPtr := flag.String("cmd", "echo", "a string")
	argsPtr := flag.String("args", "hello world", "many strings")
	flag.Parse()

	go func() {
		for {
			getStatus := process.PassArguments(cmdPtr, argsPtr)
			if getStatus {
				os.Exit(0)
			}
		}
	}()
	select {}
}
