package process

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
)

func PassArguments(cmdName *string, cmdArgs *string) bool {

	args := strings.Split(*cmdArgs, " ")
	cmd := exec.Command(*cmdName, args...)

	var wg sync.WaitGroup
	getStd := make(chan string)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("cmd.StdoutPipe(): %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("cmd.StderrPipe(): %v", err)
	}

	log.Printf("Run job and wait...")

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	childPid := cmd.Process.Pid
	log.Printf("Job PID=%d \n", childPid)

	scannerStdout := bufio.NewScanner(stdout)
	wg.Add(1)
	go func() {
		for scannerStdout.Scan() {
			m := scannerStdout.Text()
			getStd <- m
		}
		wg.Done()
	}()

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
