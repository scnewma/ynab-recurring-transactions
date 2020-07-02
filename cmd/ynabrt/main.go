package main

import (
	"fmt"
	"os"

	"github.com/scnewma/ynabrt/cmd"
)

func main() {
	rootCmd, err := cmd.NewYNABRTCommand()
	if err != nil {
		fmt.Printf("ERROR: %v", err.Error())
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
