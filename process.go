package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var BarVar = getEnv("HOME", "/tmp")
var childPid int

func main() {

	// playWithEnvs()

	// passSomeArguments("ping", "-c", "1", "localhost")

	// go run process.go -cmd=ls -args="-la /tmp"
	// go run process.go -cmd="./run.sh" -args="true"
	cmdPtr := flag.String("cmd", "echo", "a string")
	argsPtr := flag.String("args", "hello world", "many strings")
	flag.Parse()
	// getStatus := passArguments(cmdPtr, argsPtr)
	// fmt.Println(getStatus)

	go func() {
		for {
			getStatus := passArguments(cmdPtr, argsPtr)
			if getStatus {
				os.Exit(0)
			}
		}
	}()
	select {}
}

// get 'key' environment variable if exist on HOST machine otherwise return defalutValue
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func playWithEnvs() {
	// set envs only when script is run not in HOST system
	os.Setenv("FOO", "1")
	os.Setenv("BAR", "2")

	// fmt.Println("FOO:", os.Getenv("FOO"))
	// fmt.Println("BAR:", os.Getenv("BAR"))

	// `os.Environ` list all key/value pairs in the environment.
	// This returns a slice of strings in the
	// form `KEY=value`. You can `strings.Split` them to
	// get the key and value.

	// 1 - print ALL env from HOST system + include also from 'os.Setenv'
	// for _, e := range os.Environ() {
	// 	pair := strings.Split(e, "=")
	// 	fmt.Println(pair[0], pair[1])
	// }

	// 2 - print specific keys
	// for _, e := range os.Environ() {
	// 	pair := strings.Split(e, "=")
	// 	if pair[0] == "FOO" || pair[0] == "BAR" {
	// 		fmt.Println(pair[0], pair[1])
	// 		// cmd := exec.Command("echo", pair[0])
	// 		// fmt.Println(cmd)
	// 	}
	// }
}

func passSomeArguments(cmdName string, cmdArgs ...string) {
	fmt.Println(cmdArgs)
	cmd := exec.Command(cmdName, cmdArgs...)
	fmt.Println(cmd)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out)
}

// func passArguments(wg *sync.WaitGroup, cmdName *string, cmdArgs *string) {
func passArguments(cmdName *string, cmdArgs *string) bool {

	// Start a process
	args := strings.Split(*cmdArgs, " ")
	cmd := exec.Command(*cmdName, args...)

	// goroutines stuff
	var wg sync.WaitGroup
	getStd := make(chan string)

	// display output when completed OR line by line if needed
	// stdout, _ := cmd.StdoutPipe()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("cmd.StdoutPipe(): %v", err)
	}

	// stderr, _ := cmd.StderrPipe()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("cmd.StderrPipe(): %v", err)
	}

	log.Printf("Run job and wait...")

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// get and set PID
	childPid = cmd.Process.Pid
	log.Printf("Job PID=%d \n", childPid)

	// display stdout
	scannerStdout := bufio.NewScanner(stdout)
	wg.Add(1)
	go func() {
		for scannerStdout.Scan() {
			m := scannerStdout.Text()
			getStd <- m
		}
		wg.Done()
	}()

	// display stderr
	scannerStderr := bufio.NewScanner(stderr)
	wg.Add(1)
	go func() {
		for scannerStderr.Scan() {
			m := scannerStderr.Text()
			getStd <- m
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(getStd)
	}()

	for o := range getStd {
		fmt.Println(o)
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("Job finished with: %v", err)
		return false
	}

	log.Printf("Job completed.")
	return true
}
