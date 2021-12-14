# GoFrame Polaris Plugin

[![Go Reference](https://pkg.go.dev/badge/github.com/gogf/polaris.svg)](https://pkg.go.dev/github.com/gogf/polaris)
[![GoFrame Polaris CI](https://github.com/gogf/polaris/actions/workflows/go.yml/badge.svg)](https://github.com/gogf/polaris/actions/workflows/go.yml)
[![Go Report](https://goreportcard.com/badge/github.com/gogf/polaris?v=1)](https://goreportcard.com/report/github.com/gogf/polaris)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/gogf/polaris)
[![License](https://img.shields.io/github/license/gogf/polaris.svg?style=flat)](https://github.com/gogf/polaris)

English | [简体中文](README_ZH.MD)

Use Polaris mesh as service registration management and heartbeat reporting


## Installation
```
go get -u -v github.com/gogf/polaris
```
suggested using `go.mod`:
```
require github.com/gogf/polaris latest
```

## Limitation
```
golang version >= 1.16
```

Use yaml configuration file by default

Assuming `config.yaml` configuration

```yaml
server:
  Address: :8199
  ServerRoot: ./
  ServerAgent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36
  LogPath: ./logs/server
  OpenApiPath: /api.json
  SwaggerPath: /swagger
logger:
  Path: ./logs
  Level: all
  Stdout: true
database:
  link: mysql:root:root@tcp(127.0.0.1:3306)/test
  debug: true
  logger:
    Path: ./logs/sql
    Level: all
    Stdout: true

exam:
  env: app_prod

app:
  env: local
  environment: develop
  version: 1.0.0
  jaegerEndpoint: http://tracing-analysis-dc-bj.aliyuncs.com

polaris:
    # system configuration
    config:
      # Whether the heartbeat is reported
      isHeartbeat: 1
      loggerPath: ./logs/polaris
      backupPath: ./logs/polaris
    server:
      namespace: default
      service: houseme
      version: 1.0.0
      port: 8199
      ttl: 2

global:
  consumer:
    localCache:
      persistDir: $HOME/logs/backup
  serverConnector:
    connectTimeout: 1000ms
    addresses:
    - 192.168.100.222:8091
#     - 192.168.100.222:8090
```

## Example

[config.toml](example/config.yaml)

[main.go](example/main.go)

## Reference example

```go
import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"

	"github.com/gogf/polaris"
)

func main() {
	ctx := context.Background()
	adapter, err := gcfg.NewAdapterFile("config.yaml")
	if err != nil {
		g.Log().Fatal(ctx, "boot init g cfg.NewAdapterFile error:", err)
	}
	g.Cfg().SetAdapter(adapter)

	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello, polaris!")
	})
	s.Plugin(polaris.GfPolarisPlugin{
		Listener: func(config string) {
			g.Log().Print(ctx, "Polaris register success")
		},
	})
	s.Run()
}

```

## License

`GoFrame Polaris` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.

