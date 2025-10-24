package main

import (
	"gochat/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
