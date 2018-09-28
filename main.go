package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/elitecodegroovy/gmessage/server"
)

const (
	GmVersion = "1.0.0"
)

var usageStr = `
使用: gmessage [可选项]

服务器可选项:
    -a, --addr <host>                绑定主机地址 (默认: 0.0.0.0)
    -p, --port <port>                为客户端连接的端口 (默认: 6222)
    -P, --pid <file>                 存储PID的文件
    -m, --http_port <port>           http监控端口
    -ms,--https_port <port>          https监控端口
    -c, --config <file>              配置文件
    -sl,--signal <signal>[=<pid>]    发送信号给系统进程 (停止、退出、重新打开，重新加载)
        --client_advertise <string>  客户端的URL告知给其他服务器

日志可选项:
    -l, --log <file>                 文件重定向到日志输出
    -T, --logtime                    时间戳日志条目 (默认: true)
    -s, --syslog                     记录日志到系统日志或者windows事件日志中
    -r, --remote_syslog <addr>       系统日志服务地址(例如：udp://localhost:514)
    -D, --debug                      启动调试输出信息
    -V, --trace                      跟踪原始协议
    -DV                              调试和跟踪

授权可选项:
        --user <user>                连接时刻需要的用户名称
        --pass <password>            连接时刻需要的密码
        --auth <token>               连接时刻需要的授权标识符

TLS可选项:
        --tls                        启动TLS, 不必验证客户端(默认: false)
        --tlscert <file>             服务器证书文件
        --tlskey <file>              服务器证书私有密钥
        --tlsverify                  启动 TLS, 验证客户端证书
        --tlscacert <file>           验证时刻提供的客户端证书

集群可选项:
        --routes <rurl-1, rurl-2>    连接时刻的路由（Routes）
        --cluster <cluster-url>      恳求路由的集群URL链接
        --no_advertise <bool>        告知已知集群IP到客户端
        --cluster_advertise <string> 集群URL去告知其他服务器
        --connect_retries <number>   对于隐含的路由，设置重连次数


一般可选项:
    -h, --help                       显示这个消息
    -v, --version                    显示版本
        --help_tls                   TLS 帮助
`

// usage will print out the flag options for the server.
func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

// PrintServerAndExit will print our version and exit.
func PrintGMServerAndExit() {
	fmt.Printf("gmessage-server version %s\n", GmVersion)
	os.Exit(0)
}

// PrintAndDie is exported for access in other packages.
func PrintAndExit(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	os.Exit(1)
}

func main() {
	// Create a FlagSet and sets the usage
	fs := flag.NewFlagSet("gmesssage-server", flag.ExitOnError)

	// Configure the options from the flags/config file
	opts, err := server.ConfigureOptions(fs, os.Args[1:],
		PrintGMServerAndExit,
		usage,
		server.PrintTLSHelpAndExit)
	if err != nil {
		PrintAndExit(err.Error() + "\n" + usageStr)
	}

	// Create the server with appropriate options.
	s := server.New(opts)

	// Configure the logger based on the flags
	s.ConfigureLogger()

	// Start things up. Block here until done.
	if err := server.Run(s); err != nil {
		PrintAndExit(err.Error())
	}
}
