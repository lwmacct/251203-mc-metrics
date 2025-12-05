// Package export 提供 mc-vmexport 命令
package export

import (
	"github.com/lwmacct/251203-mc-metrics/internal/command"
	"github.com/lwmacct/251203-mc-metrics/internal/version"
	"github.com/urfave/cli/v3"
)

// Command mc-vmexport 根命令
var Command = &cli.Command{
	Name:      "mc-vmexport",
	Usage:     "VictoriaMetrics 数据导出工具",
	ArgsUsage: "<match>",
	Before:    command.BeforeLoadConfig,
	Action:    actionExportJSON,
	Commands: []*cli.Command{
		jsonCommand,
		csvCommand,
		nativeCommand,
		version.Command,
	},
	Flags: exportFlags(),
}

func init() {
	Command.Commands = append(Command.Commands, command.NewCompletionCommand(Command))
}

// exportFlags 返回导出命令的 flags (基础 + 导出特定)
func exportFlags() []cli.Flag {
	return append(command.BaseFlags(),
		// 导出参数
		&cli.StringFlag{
			Name:  "start",
			Usage: "开始时间 (RFC3339 或 Unix 时间戳)",
		},
		&cli.StringFlag{
			Name:  "end",
			Usage: "结束时间 (RFC3339 或 Unix 时间戳)",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "输出文件路径 (默认: stdout)",
		},
		&cli.BoolFlag{
			Name:  "gzip",
			Usage: "启用 gzip 压缩输出",
		},
	)
}

// jsonCommand json 子命令 (显式)
var jsonCommand = &cli.Command{
	Name:      "json",
	Usage:     "导出 JSON Line 格式",
	ArgsUsage: "<match>...",
	Action:    actionExportJSON,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "max-rows-per-line",
			Usage: "每行最大样本数",
		},
	},
}

// csvCommand csv 子命令
var csvCommand = &cli.Command{
	Name:      "csv",
	Usage:     "导出 CSV 格式",
	ArgsUsage: "<match>...",
	Action:    actionExportCSV,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "csv-format",
			Usage: "CSV 列定义 (默认: __name__,__value__,__timestamp__:unix_s)",
			Value: "__name__,__value__,__timestamp__:unix_s",
		},
		&cli.BoolFlag{
			Name:  "reduce-mem-usage",
			Usage: "跳过去重以减少内存使用",
		},
	},
}

// nativeCommand native 子命令
var nativeCommand = &cli.Command{
	Name:      "native",
	Usage:     "导出 Native 二进制格式 (最高效率)",
	ArgsUsage: "<match>...",
	Action:    actionExportNative,
}
