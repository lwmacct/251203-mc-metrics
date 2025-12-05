// Package query 提供 mc-vmquery 命令
package query

import (
	"time"

	"github.com/lwmacct/251203-mc-metrics/internal/command"
	"github.com/lwmacct/251203-mc-metrics/internal/version"
	"github.com/urfave/cli/v3"
)

// Command 根命令
var Command = &cli.Command{
	Name:      "mc-vmquery",
	Usage:     "VictoriaMetrics MetricsQL 查询客户端",
	ArgsUsage: "[query]",
	Before:    command.BeforeLoadConfig,
	Action:    actionQuery,
	Commands: []*cli.Command{
		metricsCommand,
		labelsCommand,
		labelValuesCommand,
		seriesCommand,
		version.Command,
	},
	Flags: queryFlags(),
}

func init() {
	// 延迟添加 completion 命令，避免循环引用
	Command.Commands = append(Command.Commands, command.NewCompletionCommand(Command))
}

// queryFlags 返回查询命令的 flags (基础 + 查询特定)
func queryFlags() []cli.Flag {
	return append(command.BaseFlags(),
		// 输出配置
		&cli.StringFlag{
			Name:    "output-format",
			Aliases: []string{"o"},
			Usage:   "输出格式: table, json, csv, graph",
			Value:   command.Defaults.Output.Format,
		},
		&cli.BoolFlag{
			Name:  "output-no-headers",
			Usage: "禁用表头输出",
			Value: command.Defaults.Output.NoHeaders,
		},
		// 查询参数
		&cli.StringFlag{
			Name:  "time",
			Usage: "查询时间点 (RFC3339 格式或 'now')",
			Value: "now",
		},
		&cli.DurationFlag{
			Name:  "range",
			Usage: "范围查询的时间跨度 (如 1h, 30m)",
		},
		&cli.DurationFlag{
			Name:  "step",
			Usage: "范围查询的步长 (如 1m, 15s)",
			Value: time.Minute,
		},
	)
}

// metricsCommand metrics 子命令
var metricsCommand = &cli.Command{
	Name:      "metrics",
	Usage:     "列出所有指标名称",
	ArgsUsage: "[match]",
	Action:    actionMetrics,
}

// labelsCommand labels 子命令
var labelsCommand = &cli.Command{
	Name:   "labels",
	Usage:  "列出所有标签名称",
	Action: actionLabels,
}

// labelValuesCommand label-values 子命令
var labelValuesCommand = &cli.Command{
	Name:      "label-values",
	Usage:     "获取指定标签的所有值",
	ArgsUsage: "<label>",
	Action:    actionLabelValues,
}

// seriesCommand series 子命令
var seriesCommand = &cli.Command{
	Name:      "series",
	Usage:     "列出匹配的时间序列",
	ArgsUsage: "<match>",
	Action:    actionSeries,
}
