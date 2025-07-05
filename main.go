package main

import (
	"go-burrokuchen/cmd"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
}

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
