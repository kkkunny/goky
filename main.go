package main

import (
	"fmt"
	"github.com/kkkunny/klang/cmd"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     "kcc",
	Short:   "The compiler for the K programming language",
	Version: "v0.1",
}

func main() {
	rootCmd.AddCommand(cmd.BuildCmd(), cmd.RunCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
