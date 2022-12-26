package cmd

import (
	"errors"
	stlos "github.com/kkkunny/stl/os"
	"github.com/kkkunny/stl/util"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

func RunCmd() *cobra.Command {
	var conf buildConfig
	cmd := &cobra.Command{
		Use:   "run",
		Short: "compiler and then run a k source file",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}
			// conf.Backend = "asm"
			target := stlos.Path(args[0])
			if !target.IsExist() {
				return errors.New("expect a k source file path")
			}
			target, err := target.GetAbsolute()
			if err != nil {
				return err
			}
			conf.Target = target
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			util.Must(run(conf, args[1:]))
			return nil
		},
	}
	return cmd
}

func run(conf buildConfig, args []string) error {
	if err := build(conf); err != nil {
		return err
	}
	var binary stlos.Path
	if !conf.Target.IsDir() {
		binary = conf.Target.WithExtension("out")
	} else {
		binary = conf.Target.Join(conf.Target.GetBase().WithExtension("out"))
	}
	defer os.Remove(binary.String())

	cmd := exec.Command(binary.String(), args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			_ = os.Remove(binary.String())
			os.Exit(exit.ExitCode())
		} else {
			return err
		}
	}
	return nil
}
