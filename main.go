package main

import (
	"bank-system-go/cmd"
)

//go:generate wire ./...

func main() {
	cmd.Execute()
}
