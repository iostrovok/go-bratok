package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func main() {

	cmd := exec.Command(`./test.pl`)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	e, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("STDOUT: %s\n", b)
	fmt.Printf("STDERR: %s\n", e)
}
