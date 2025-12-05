package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

// NewCompletionCommand 创建 completion 子命令
// 自动从传入的 rootCmd 生成 zsh 补全脚本
func NewCompletionCommand(rootCmd *cli.Command) *cli.Command {
	return &cli.Command{
		Name:   "completion",
		Usage:  "生成 zsh 补全脚本",
		Hidden: true, // 不在帮助中显示，也不出现在补全列表
		Description: fmt.Sprintf(`生成 zsh 补全脚本。

启用补全:

  # 确保 completions 目录在 fpath 中
  echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
  echo 'autoload -Uz compinit && compinit' >> ~/.zshrc

  # 生成补全脚本
  mkdir -p ~/.zsh/completions
  %s completion > ~/.zsh/completions/_%s

  # 重新加载 zsh
  exec zsh
`, rootCmd.Name, rootCmd.Name),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return GenerateZsh(os.Stdout, rootCmd)
		},
	}
}

// GenerateZsh 从 cli.Command 自动生成 zsh 补全脚本
func GenerateZsh(w io.Writer, cmd *cli.Command) error {
	funcName := toZshFuncName(cmd.Name)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("#compdef %s\n\n", cmd.Name))
	sb.WriteString(fmt.Sprintf("# %s zsh completion script (auto-generated)\n\n", cmd.Name))

	// 生成主函数
	generateZshFunction(&sb, cmd, funcName, true)

	// 生成子命令函数
	generateSubcommandFunctions(&sb, cmd, funcName)

	sb.WriteString(fmt.Sprintf("compdef %s %s\n", funcName, cmd.Name))

	_, err := io.WriteString(w, sb.String())
	return err
}

// generateZshFunction 生成单个命令的 zsh 补全函数
func generateZshFunction(sb *strings.Builder, cmd *cli.Command, funcName string, isRoot bool) {
	fmt.Fprintf(sb, "%s() {\n", funcName)
	sb.WriteString("    local curcontext=\"$curcontext\" state line\n")
	sb.WriteString("    typeset -A opt_args\n\n")

	// 收集 flags
	flags := collectFlags(cmd, isRoot)
	if len(flags) > 0 {
		sb.WriteString("    local -a flags\n")
		sb.WriteString("    flags=(\n")
		for _, f := range flags {
			fmt.Fprintf(sb, "        %s\n", f)
		}
		sb.WriteString("    )\n\n")
	}

	// 收集可见的子命令（只有需要展开的才处理）
	subcommands := getVisibleCommands(cmd)
	hasSubcommands := len(subcommands) > 0 && shouldExpandSubcommands(cmd)

	// 生成 _arguments 调用
	sb.WriteString("    _arguments -C \\\n")
	if len(flags) > 0 {
		sb.WriteString("        $flags \\\n")
	}
	if hasSubcommands {
		fmt.Fprintf(sb, "        '1: :%s_commands' \\\n", funcName)
		sb.WriteString("        '*::arg:->args'\n")
	} else {
		sb.WriteString("        '*:file:_files'\n")
	}

	// 生成子命令状态处理
	if hasSubcommands {
		sb.WriteString("\n    case $state in\n")
		sb.WriteString("        args)\n")
		sb.WriteString("            case $line[1] in\n")
		for _, sub := range subcommands {
			subFuncName := funcName + "_" + toZshFuncName(sub.Name)
			// 包含别名
			names := []string{sub.Name}
			names = append(names, sub.Aliases...)
			fmt.Fprintf(sb, "                %s)\n", strings.Join(names, "|"))
			fmt.Fprintf(sb, "                    %s\n", subFuncName)
			sb.WriteString("                    ;;\n")
		}
		sb.WriteString("            esac\n")
		sb.WriteString("            ;;\n")
		sb.WriteString("    esac\n")
	}

	sb.WriteString("}\n\n")
}

// generateSubcommandFunctions 递归生成所有子命令的函数
func generateSubcommandFunctions(sb *strings.Builder, cmd *cli.Command, parentFuncName string) {
	subcommands := getVisibleCommands(cmd)
	if len(subcommands) == 0 {
		return
	}

	// 生成 _commands 函数
	fmt.Fprintf(sb, "%s_commands() {\n", parentFuncName)
	sb.WriteString("    local -a commands\n")
	sb.WriteString("    commands=(\n")
	for _, sub := range subcommands {
		usage := strings.ReplaceAll(sub.Usage, "'", "'\\''")
		fmt.Fprintf(sb, "        '%s:%s'\n", sub.Name, usage)
	}
	sb.WriteString("    )\n")
	sb.WriteString("    _describe -t commands 'commands' commands\n")
	sb.WriteString("}\n\n")

	// 递归生成每个子命令的函数
	for _, sub := range subcommands {
		subFuncName := parentFuncName + "_" + toZshFuncName(sub.Name)
		generateZshFunction(sb, sub, subFuncName, false)
		// 只有需要展开的命令才递归
		if shouldExpandSubcommands(sub) {
			generateSubcommandFunctions(sb, sub, subFuncName)
		}
	}
}

// collectFlags 收集命令的 flags，转换为 zsh 格式
func collectFlags(cmd *cli.Command, includeGlobal bool) []string {
	var flags []string
	seen := make(map[string]bool)

	// 收集当前命令的 flags
	for _, f := range cmd.Flags {
		zshFlag := flagToZsh(f)
		if zshFlag != "" && !seen[zshFlag] {
			flags = append(flags, zshFlag)
			seen[zshFlag] = true
		}
	}

	// 如果是子命令，也收集父命令的 flags（通过 root 传递）
	if includeGlobal {
		// help flag
		flags = append(flags, "'(- *)'{-h,--help}'[显示帮助信息]'")
	}

	return flags
}

// flagToZsh 将 cli.Flag 转换为 zsh 补全格式
func flagToZsh(f cli.Flag) string {
	names := f.Names()
	if len(names) == 0 {
		return ""
	}

	// 获取 flag 的描述和其他属性
	usage := ""
	takesValue := false
	valueType := ""

	switch flag := f.(type) {
	case *cli.StringFlag:
		usage = flag.Usage
		takesValue = true
		valueType = getValueCompletion(flag.Name, flag.Usage)
	case *cli.BoolFlag:
		usage = flag.Usage
		takesValue = false
	case *cli.IntFlag:
		usage = flag.Usage
		takesValue = true
		valueType = ":number:"
	case *cli.DurationFlag:
		usage = flag.Usage
		takesValue = true
		valueType = ":duration:"
	case *cli.StringSliceFlag:
		usage = flag.Usage
		takesValue = true
		valueType = ":value:"
	default:
		// 其他类型，尝试获取基本信息
		if nf, ok := f.(interface{ GetUsage() string }); ok {
			usage = nf.GetUsage()
		}
	}

	usage = strings.ReplaceAll(usage, "'", "'\\''")
	usage = strings.ReplaceAll(usage, "[", "(")
	usage = strings.ReplaceAll(usage, "]", ")")

	// 构建 zsh flag 字符串
	if len(names) == 1 {
		name := names[0]
		if len(name) == 1 {
			// 短选项
			if takesValue {
				return fmt.Sprintf("'-%s[%s]%s'", name, usage, valueType)
			}
			return fmt.Sprintf("'-%s[%s]'", name, usage)
		}
		// 长选项
		if takesValue {
			return fmt.Sprintf("'--%s[%s]%s'", name, usage, valueType)
		}
		return fmt.Sprintf("'--%s[%s]'", name, usage)
	}

	// 有别名的情况（如 -c, --config）
	var short, long string
	for _, n := range names {
		if len(n) == 1 {
			short = "-" + n
		} else {
			long = "--" + n
		}
	}

	if short != "" && long != "" {
		if takesValue {
			return fmt.Sprintf("'(%s %s)'{%s,%s}'[%s]%s'", short, long, short, long, usage, valueType)
		}
		return fmt.Sprintf("'(%s %s)'{%s,%s}'[%s]'", short, long, short, long, usage)
	}

	// fallback
	name := names[0]
	prefix := "--"
	if len(name) == 1 {
		prefix = "-"
	}
	if takesValue {
		return fmt.Sprintf("'%s%s[%s]%s'", prefix, name, usage, valueType)
	}
	return fmt.Sprintf("'%s%s[%s]'", prefix, name, usage)
}

// getValueCompletion 根据 flag 名称和描述推断补全类型
// 设计原则：从 Usage 描述推断，不硬编码业务值
func getValueCompletion(name, usage string) string {
	nameLower := strings.ToLower(name)
	usageLower := strings.ToLower(usage)

	// 1. 优先从 Usage 解析枚举值（如 "类型: a, b, c" 或 "format: json, csv"）
	if values := parseEnumFromUsage(usage); len(values) > 0 {
		return fmt.Sprintf(":value:(%s)", strings.Join(values, " "))
	}

	// 2. URL 类型（从 name 推断）
	if strings.Contains(nameLower, "url") {
		return ":url:"
	}

	// 3. 文件路径类型（从 name 或 usage 推断）
	if isFilePath(nameLower, usageLower) {
		return ":file:_files"
	}

	// 4. 数字类型
	if strings.Contains(usageLower, "number") ||
		strings.Contains(usageLower, "数量") ||
		strings.Contains(usageLower, "个数") {
		return ":number:"
	}

	return ":value:"
}

// parseEnumFromUsage 从 Usage 描述中解析枚举值
// 支持格式：
//   - "类型: a, b, c"
//   - "format: json, csv, xml"
//   - "模式 (a/b/c)"
//   - "type (a|b|c)"
func parseEnumFromUsage(usage string) []string {
	// 模式1: "xxx: a, b, c" 或 "xxx：a, b, c"（中英文冒号）
	if idx := strings.IndexAny(usage, ":："); idx != -1 {
		rest := strings.TrimSpace(usage[idx+1:])
		// 去掉括号内容（如果有的话，可能是补充说明）
		if parenIdx := strings.IndexAny(rest, "(（"); parenIdx != -1 {
			rest = strings.TrimSpace(rest[:parenIdx])
		}
		// 按逗号分割
		if strings.Contains(rest, ",") || strings.Contains(rest, "，") {
			parts := strings.FieldsFunc(rest, func(r rune) bool {
				return r == ',' || r == '，' || r == ' '
			})
			var values []string
			for _, p := range parts {
				p = strings.TrimSpace(p)
				// 只保留简单的值（无空格、非空）
				if p != "" && !strings.Contains(p, " ") && len(p) < 20 {
					values = append(values, p)
				}
			}
			if len(values) >= 2 {
				return values
			}
		}
	}

	// 模式2: "(a/b/c)" 或 "(a|b|c)"
	if start := strings.IndexAny(usage, "(（"); start != -1 {
		if end := strings.IndexAny(usage[start:], ")）"); end != -1 {
			inner := usage[start+1 : start+end]
			// 检查是否是枚举格式
			if strings.ContainsAny(inner, "/|") && !strings.Contains(inner, " ") {
				parts := strings.FieldsFunc(inner, func(r rune) bool {
					return r == '/' || r == '|'
				})
				if len(parts) >= 2 {
					var values []string
					for _, p := range parts {
						p = strings.TrimSpace(p)
						if p != "" && len(p) < 20 {
							values = append(values, p)
						}
					}
					return values
				}
			}
		}
	}

	return nil
}

// isFilePath 判断是否是文件路径类型
// 从 flag 名称和 usage 描述推断
func isFilePath(nameLower, usageLower string) bool {
	// 从 name 推断
	fileNamePatterns := []string{
		"file", "path", "config", "input", "output",
		"cert", "key", "ca", // 证书相关
	}
	for _, pattern := range fileNamePatterns {
		if strings.Contains(nameLower, pattern) {
			// 排除 "prefix"、"format" 等误判
			if strings.Contains(nameLower, "prefix") ||
				strings.Contains(nameLower, "format") {
				return false
			}
			return true
		}
	}

	// 从 usage 推断（中英文）
	fileUsagePatterns := []string{
		"file", "path", "文件", "路径", "证书",
	}
	for _, pattern := range fileUsagePatterns {
		if strings.Contains(usageLower, pattern) {
			return true
		}
	}

	return false
}

// getVisibleCommands 获取可见的子命令（排除 hidden 和特殊命令）
func getVisibleCommands(cmd *cli.Command) []*cli.Command {
	var visible []*cli.Command
	for _, sub := range cmd.Commands {
		// 跳过隐藏命令
		if sub.Hidden {
			continue
		}
		// 跳过 help、completion 等不需要在补全中显示的命令
		if sub.Name == "help" || sub.Name == "completion" {
			continue
		}
		visible = append(visible, sub)
	}
	return visible
}

// shouldExpandSubcommands 判断是否需要展开子命令的补全
// version 等终端命令不需要展开其子命令
func shouldExpandSubcommands(cmd *cli.Command) bool {
	// version 命令的子命令（short、json）不需要在补全中展开
	if cmd.Name == "version" {
		return false
	}
	return true
}

// toZshFuncName 将命令名转换为合法的 zsh 函数名
func toZshFuncName(name string) string {
	// 替换 - 为 _，添加前缀 _
	return "_" + strings.ReplaceAll(name, "-", "_")
}
