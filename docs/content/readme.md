# vm-metrics

VictoriaMetrics 命令行工具集，支持 MetricsQL 查询、数据导入导出。

## 安装

```bash
go install github.com/lwmacct/251203-vm-metrics/cmd/vm-metrics@latest

# 或单独安装
go install github.com/lwmacct/251203-vm-metrics/cmd/vm-query@latest
go install github.com/lwmacct/251203-vm-metrics/cmd/vm-export@latest
go install github.com/lwmacct/251203-vm-metrics/cmd/vm-import@latest
```

## 快速开始

### 配置文件

- 配置文件参考：[config.yaml 示例](config/config.example.yaml)
- 使用 `--config <path>` 指定配置文件路径
  > 默认配置文件路径（按顺序搜索）：
  >
  > - `./config.yaml`
  > - `./config/config.yaml`
  > - `$HOME/.vm-metrics.yaml`
  > - `/etc/vm-metrics/config.yaml`

### 命令示例

```bash
# 统一命令
vm-metrics query 'up{job="prometheus"}'
vm-metrics export '{job="node"}' -o data.json
vm-metrics import data.json

# 使用别名
vm-metrics q -o json 'up'           # query
vm-metrics e '{job="node"}'         # export
vm-metrics i data.json              # import

# 独立命令 (等效)
vm-query 'up{job="prometheus"}'
vm-export '{job="node"}'
vm-import data.json
```

## 命令结构

```
vm-metrics                      # 统一入口
├── query (q)                   # MetricsQL 查询
│   ├── metrics                 # 列出所有指标
│   ├── labels                  # 列出所有标签
│   ├── label-values <label>    # 获取标签值
│   └── series <match>          # 列出时间序列
├── export (e)                  # 数据导出
│   ├── json                    # JSON Line 格式
│   ├── csv                     # CSV 格式
│   └── native                  # 原生二进制格式
├── import (i)                  # 数据导入
│   ├── json                    # JSON Line 格式
│   ├── csv                     # CSV 格式
│   ├── native                  # 原生二进制格式
│   └── prometheus              # Prometheus 格式
└── version                     # 版本信息
```

## 输出格式

```bash
# 表格 (默认)
vm-metrics query 'up'

# JSON
vm-metrics query -o json 'up'

# CSV
vm-metrics query -o csv 'up'

# ASCII 图表 (仅 range query)
vm-metrics query -o graph --range 1h 'rate(http_requests_total[5m])'
```

## 相关链接

- [VictoriaMetrics 文档](https://docs.victoriametrics.com/)
- [MetricsQL 文档](https://docs.victoriametrics.com/victoriametrics/metricsql/)

## 工具链

- [Taskfile](https://taskfile.dev) - 项目 CLI 管理
- [Pre-commit](https://pre-commit.com/) - Git 钩子管理
