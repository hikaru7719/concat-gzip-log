package main

import (
	"log"
	"os"

	"github.com/hikaru7719/concat-gzip-log/cmd"
)

func main() {
	if err := cmd.NewCommand().Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
