// Package importcmd 提供 mc-vmimport 命令
// 注意：包名使用 importcmd 因为 import 是 Go 保留字
package importcmd

import (
	"github.com/lwmacct/251203-mc-metrics/internal/command"
	"github.com/lwmacct/251203-mc-metrics/internal/version"
	"github.com/urfave/cli/v3"
)

// Command mc-vmimport 根命令
var Command = &cli.Command{
	Name:      "mc-vmimport",
	Usage:     "VictoriaMetrics 数据导入工具",
	ArgsUsage: "[file]",
	Before:    command.BeforeLoadConfig,
	Action:    actionImportJSON,
	Commands: []*cli.Command{
		jsonCommand,
		csvCommand,
		nativeCommand,
		prometheusCommand,
		version.Command,
	},
	Flags: importFlags(),
}

func init() {
	Command.Commands = append(Command.Commands, command.NewCompletionCommand(Command))
}

// importFlags 返回导入命令的 flags (基础 + 导入特定)
func importFlags() []cli.Flag {
	return append(command.BaseFlags(),
		// 导入参数
		&cli.StringFlag{
			Name:    "input",
			Aliases: []string{"i"},
			Usage:   "输入文件路径 (默认: stdin)",
		},
		&cli.BoolFlag{
			Name:  "gzip",
			Usage: "输入为 gzip 压缩格式",
		},
	)
}

// jsonCommand json 子命令
var jsonCommand = &cli.Command{
	Name:      "json",
	Usage:     "导入 JSON Line 格式",
	ArgsUsage: "[file]",
	Action:    actionImportJSON,
}

// csvCommand csv 子命令
var csvCommand = &cli.Command{
	Name:      "csv",
	Usage:     "导入 CSV 格式",
	ArgsUsage: "[file]",
	Action:    actionImportCSV,
}

// nativeCommand native 子命令
var nativeCommand = &cli.Command{
	Name:      "native",
	Usage:     "导入 Native 二进制格式",
	ArgsUsage: "[file]",
	Action:    actionImportNative,
}

// prometheusCommand prometheus 子命令
var prometheusCommand = &cli.Command{
	Name:      "prometheus",
	Usage:     "导入 Prometheus exposition 格式",
	ArgsUsage: "[file]",
	Action:    actionImportPrometheus,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "job",
			Usage: "Pushgateway job 标签",
		},
		&cli.StringFlag{
			Name:  "instance",
			Usage: "Pushgateway instance 标签",
		},
	},
}
