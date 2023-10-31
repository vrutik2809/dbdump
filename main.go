package main

import (
	"log"

	"github.com/vrutik2809/dbdump/cmd/root"
)

func main() {
	if err := root.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
