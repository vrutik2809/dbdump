package main

import (
	"fmt"
	"os"

	"github.com/vrutik2809/dbdump/cmd/root"
)

func main() {
	if err := root.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
