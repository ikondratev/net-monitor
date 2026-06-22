package main

import (
	"log"
	"os"

	"github.com/ikondratev/net-monitor/lib/cli"
)

func main() {
	if err := cli.Run(os.Args[1:], os.Stdout); err != nil {
		log.Fatal(err)
	}
}
