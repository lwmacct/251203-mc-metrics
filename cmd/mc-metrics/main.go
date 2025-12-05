package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lwmacct/251203-mc-metrics/internal/command"
	"github.com/lwmacct/251203-mc-metrics/internal/command/export"
	importcmd "github.com/lwmacct/251203-mc-metrics/internal/command/import"
	"github.com/lwmacct/251203-mc-metrics/internal/command/query"
	"github.com/lwmacct/251203-mc-metrics/internal/version"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:    "mc-metrics",
		Usage:   "VictoriaMetrics 统一命令行工具",
		Version: version.GetVersion(),
		Commands: []*cli.Command{
			queryCommand(),
			exportCommand(),
			importCommand(),
			version.Command,
		},
		Flags: command.BaseFlags(),
	}
	// 动态添加 completion 命令
	app.Commands = append(app.Commands, command.NewCompletionCommand(app))

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// queryCommand 包装 query 命令
func queryCommand() *cli.Command {
	cmd := *query.Command // 复制命令
	cmd.Name = "query"
	cmd.Aliases = []string{"q"}
	cmd.Usage = "执行 MetricsQL 查询"
	return &cmd
}

// exportCommand 包装 export 命令
func exportCommand() *cli.Command {
	cmd := *export.Command
	cmd.Name = "export"
	cmd.Aliases = []string{"e"}
	cmd.Usage = "导出时序数据"
	return &cmd
}

// importCommand 包装 import 命令
func importCommand() *cli.Command {
	cmd := *importcmd.Command
	cmd.Name = "import"
	cmd.Aliases = []string{"i"}
	cmd.Usage = "导入时序数据"
	return &cmd
}
