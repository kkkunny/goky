//go:build !test

package main

import (
	"fmt"
	"github.com/kkkunny/Sim/cmd"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     "sim",
	Short:   "The compiler for the Sim programming language",
	Version: "v0.1",
}

func main() {
	rootCmd.AddCommand(cmd.BuildCmd(), cmd.RunCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
