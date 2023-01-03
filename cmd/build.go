package cmd

import (
	"errors"
	stlos "github.com/kkkunny/stl/os"
	"github.com/kkkunny/stl/util"
	"github.com/spf13/cobra"
	"os"
)

type buildConfig struct {
	// Backend      string       // 后端类型
	Target       stlos.Path   // 目标地址
	Release      bool         // release模式
	Output       stlos.Path   // 输出地址
	LLVM         bool         // 输出llvm ir
	Asm          bool         // 输出汇编
	Shared       bool         // 输出共享库
	Object       bool         // 输出对象文件
	Linkages     []stlos.Path // 链接
	Libraries    []string     // 链接库
	LibraryPaths []string     // 链接库地址
}

func BuildCmd() *cobra.Command {
	var conf buildConfig
	cmd := &cobra.Command{
		Use:   "build",
		Short: "compiler a k source file",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}
			target := stlos.Path(args[0])
			if !target.IsExist() {
				return errors.New("expect a goky source file path")
			}
			target, err := target.GetAbsolute()
			if err != nil {
				return err
			}
			conf.Target = target
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			util.Must(build(conf))
			return nil
		},
	}
	// model
	cmd.Flags().BoolVarP(&conf.Release, "release", "r", false, "with release model")
	// output path
	cmd.Flags().StringVarP((*string)(&conf.Output), "output", "o", "", "output path")
	// output type
	cmd.Flags().BoolVar(&conf.LLVM, "llvm", false, "output llvm ir file")
	cmd.Flags().BoolVar(&conf.Asm, "asm", false, "output asm file")
	cmd.Flags().BoolVar(&conf.Shared, "shared", false, "output shared library file")
	cmd.Flags().BoolVar(&conf.Object, "object", false, "output object file")
	// lib
	cmd.Flags().StringSliceVarP(&conf.Libraries, "lib", "l", nil, "linkage extern library")
	cmd.Flags().StringSliceVarP(&conf.LibraryPaths, "lib_path", "L", nil, "library path")
	// backend
	// cmd.Flags().StringVar(&conf.Backend, "backend", "llvm", "type of the compile backend")
	return cmd
}

func build(conf buildConfig) error {
	// 输出地址
	if conf.Output == "" {
		if !conf.Target.IsDir() {
			switch {
			case conf.LLVM:
				conf.Output = conf.Target.WithExtension("ll")
			case conf.Asm:
				conf.Output = conf.Target.WithExtension("s")
			case conf.Shared:
				conf.Output = conf.Target.GetParent().Join("lib" + conf.Target.GetBase().WithExtension("so"))
			case conf.Object:
				conf.Output = conf.Target.WithExtension("o")
			default:
				conf.Output = conf.Target.WithExtension("out")
			}
		} else {
			switch {
			case conf.LLVM:
				conf.Output = conf.Target.Join(conf.Target.GetBase().WithExtension("ll"))
			case conf.Asm:
				conf.Output = conf.Target.Join(conf.Target.GetBase().WithExtension("s"))
			case conf.Shared:
				conf.Output = conf.Target.Join("lib" + conf.Target.GetBase().WithExtension("so"))
			case conf.Object:
				conf.Output = conf.Target.Join(conf.Target.GetBase().WithExtension("o"))
			default:
				conf.Output = conf.Target.Join(conf.Target.GetBase().WithExtension("out"))
			}
		}
	}

	// llvm ir
	if conf.LLVM {
		_, _, err := outputLLVM(&conf, conf.Target, conf.Output, true)
		return err
	}
	module, targetMachine, err := outputLLVM(&conf, conf.Target, "", false)
	if err != nil {
		return err
	}

	// 汇编
	var asmPath stlos.Path
	if conf.Asm {
		asmPath = conf.Output
		_, err = outputAsm(module, targetMachine, conf.Output)
		return err
	} else {
		asmPath, err = outputAsm(module, targetMachine, "")
		if err != nil {
			return err
		}
	}
	defer os.Remove(asmPath.String())

	// 链接
	var objectPath stlos.Path
	if conf.Object {
		objectPath = conf.Output
		_, err = outputObject(asmPath, conf.Output, conf.Linkages)
		return err
	} else {
		objectPath, err = outputObject(asmPath, "", conf.Linkages)
		if err != nil {
			return err
		}
	}
	defer os.Remove(objectPath.String())

	// 动态库
	if conf.Shared {
		_, err = outputSharedFile(objectPath, conf.Output, conf.Libraries, conf.LibraryPaths)
		return err
	}

	// 可执行文件
	_, err = outputExecutableFile(objectPath, conf.Output, conf.Libraries, conf.LibraryPaths)
	return err
}
