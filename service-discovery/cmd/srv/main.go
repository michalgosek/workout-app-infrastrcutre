package main

import "log"

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	return nil
}
