package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func main() {
	var Version bool
	// 创建rootCmd主命令，并定义Run执行函数(注意，此处是定义Run函数而非直接执行该函数)。也可以通过rootCmd.AddCommand方法添加子命令。
	var rootCmd = &cobra.Command{
		Use:   "root [sub]",
		Short: "root command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside rootCmd Rum with args: %v\n", args)
			if Version {
				fmt.Printf("Version:1.0\n")
			}
		},
	}

	// 为命令添加命令行参数(Flag)。
	flags := rootCmd.Flags()
	flags.BoolVarP(&Version, "version", "v", false, "Print version information and quit")

	// 执行rootCmd命令调用的函数，rootCmd.Execute会在内部回调Run执行函数。
	_ = rootCmd.Execute()
}
