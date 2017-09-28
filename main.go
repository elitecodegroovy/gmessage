package main

import (
	"github.com/go-svc/svc"
	"syscall"
	"log"
	"path/filepath"
	"os"
	"flag"
	"github.com/elitecodegroovy/gmessage/message"
	"time"
	"math/rand"
	"fmt"
	"strings"
	"net/url"
	"net"
	"nats-io/gnatsd/server"
)


var tipsMsg = `
Usage: gnatsd [options]

Server Options:
    -a, --addr <host>                Bind to host address (default: 0.0.0.0)
    -p, --client_port <port>         Use port for clients (default: 5000)
    -P, --pid <file>                 File to store PID
    -m, --m_http_port <port>         Use port for http monitoring
    -ms,--m_https_port <port>        Use port for https monitoring
    -c, --config <file>              Configuration file
    -sl,--signal <signal>[=<pid>]    Send signal to gmessage process (stop, quit, reopen, reload)

Logging Options:
    -l, --log <file>                 File to redirect log output
    -T, --logtime                    Timestamp log entries (default: true)
    -s, --syslog                     Log to syslog or windows event log
    -r, --remote_syslog <addr>       Syslog server addr (udp://localhost:514)
    -D, --debug                      Enable debugging output
    -V, --trace                      Trace the raw protocol
    -DV                              Debug and trace

Authorization Options:
        --user <user>                User required for connections
        --pass <password>            Password required for connections
        --auth <token>               Authorization token required for connections

TLS Options:
        --tls                        Enable TLS, do not verify clients (default: false)
        --tlscert <file>             Server certificate file
        --tlskey <file>              Private key for server certificate
        --tlsverify                  Enable TLS, verify client certificates
        --tlscacert <file>           Client certificate CA for verification

Cluster Options:
        --routes <rurl-1, rurl-2>    Routes to solicit and connect
        --cluster <cluster-url>      Cluster URL for solicited routes
        --no_advertise <bool>        Advertise known cluster IPs to clients
        --connect_retries <number>   For implicit routes, number of connect retries


Common Options:
    -h, --help                       Show this message
    -v, --version                    Show version
        --help_tls                   TLS help
`


func msgFlagSet(opts *message.Options) *flag.FlagSet {
	flagSet := flag.NewFlagSet("nsqd", flag.ExitOnError)

	flagSet.IntVar(&opts.Port, "port", 0, "Port to listen on.")
	flagSet.IntVar(&opts.Port, "p", 0, "Port to listen on.")
	flagSet.StringVar(&opts.Host, "addr", "", "Network host to listen on.")
	flagSet.StringVar(&opts.Host, "a", "", "Network host to listen on.")
	flagSet.StringVar(&opts.Host, "net", "", "Network host to listen on.")
	flagSet.BoolVar(&opts.Debug, "D", false, "Enable Debug logging.")
	flagSet.BoolVar(&opts.Debug, "debug", false, "Enable Debug logging.")
	flagSet.BoolVar(&opts.Trace, "V", false, "Enable Trace logging.")
	flagSet.BoolVar(&opts.Trace, "trace", false, "Enable Trace logging.")
	//flagSet.BoolVar(&debugAndTrace, "DV", false, "Enable Debug and Trace logging.")

	flagSet.BoolVar(&opts.Logtime, "T", true, "Timestamp log entries.")
	flagSet.BoolVar(&opts.Logtime, "logtime", true, "Timestamp log entries.")
	flagSet.StringVar(&opts.Username, "user", "", "Username required for connection.")
	flagSet.StringVar(&opts.Password, "pass", "", "Password required for connection.")
	flagSet.StringVar(&opts.Authorization, "auth", "", "Authorization token required for connection.")
	flagSet.IntVar(&opts.HTTPPort, "m", 0, "HTTP Port for /varz, /connz endpoints.")
	flagSet.IntVar(&opts.HTTPPort, "http_port", 0, "HTTP Port for /varz, /connz endpoints.")
	flagSet.IntVar(&opts.HTTPSPort, "ms", 0, "HTTPS Port for /varz, /connz endpoints.")
	flagSet.IntVar(&opts.HTTPSPort, "https_port", 0, "HTTPS Port for /varz, /connz endpoints.")
	flagSet.String("configFile", "", "Configuration file.")

	//flagSet.StringVar(&configFile, "config", "", "Configuration file.")
	flagSet.String("sl", "", "Send signal to g-message process (stop, quit, reopen, reload)")
	flagSet.String("signal", "", "Send signal to g-message process (stop, quit, reopen, reload)")
	flagSet.StringVar(&opts.PidFile, "P", "", "File to store process pid.")
	flagSet.StringVar(&opts.PidFile, "pid", "", "File to store process pid.")
	flagSet.StringVar(&opts.LogFile, "l", "", "File to store logging output.")
	flagSet.StringVar(&opts.LogFile, "log", "", "File to store logging output.")
	flagSet.BoolVar(&opts.Syslog, "s", false, "Enable syslog as log method.")
	flagSet.BoolVar(&opts.Syslog, "syslog", false, "Enable syslog as log method..")
	flagSet.StringVar(&opts.RemoteSyslog, "r", "", "Syslog server addr (udp://localhost:514).")

	flagSet.StringVar(&opts.RemoteSyslog, "remote_syslog", "", "Syslog server addr (udp://localhost:514).")
	flagSet.Bool( "version", false, "Print version information.")
	flagSet.Bool("v", false, "Print version information.")

	flagSet.Bool("help", false, "show usage information.")
	flagSet.IntVar(&opts.ProfPort, "profile", 0, "Profiling HTTP port")
	flagSet.StringVar(&opts.RoutesStr, "routes", "", "Routes to actively solicit a connection.")
	flagSet.StringVar(&opts.Cluster.ListenStr, "cluster", "", "Cluster url from which members can solicit routes.")
	flagSet.StringVar(&opts.Cluster.ListenStr, "cluster_listen", "", "Cluster url from which members can solicit routes.")
	flagSet.BoolVar(&opts.Cluster.NoAdvertise, "no_advertise", false, "Advertise known cluster IPs to clients.")
	flagSet.IntVar(&opts.Cluster.ConnectRetries, "connect_retries", 0, "For implicit routes, number of connect retries")
	//flagSet.BoolVar(&showTLSHelp, "help_tls", false, "TLS help.")

	flagSet.BoolVar(&opts.TLS, "tls", false, "Enable TLS.")
	flagSet.BoolVar(&opts.TLSVerify, "tlsverify", false, "Enable TLS with client verification.")
	flagSet.StringVar(&opts.TLSCert, "tlscert", "", "Server certificate file.")
	flagSet.StringVar(&opts.TLSKey, "tlskey", "", "Private key for server certificate.")
	flagSet.StringVar(&opts.TLSCaCert, "tlscacert", "", "Client certificate CA for verification.")

	return flagSet
}

func showUsageInfoNExit(){
	fmt.Printf("%s\n", tipsMsg)
	os.Exit(0)
}


type program struct {
	msg *message.Server
}



func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}

// init starter
func (p *program) Init(env svc.Environment) error {
	log.Printf("is win service? %v\n", env.IsWindowsService())
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *program) Stop() error {
	// The Stop method is invoked by stopping the Windows service, or by pressing Ctrl+C on the console.
	// This method may block, but it's a good idea to finish quickly or your process may be killed by
	// Windows during a shutdown/reboot. As a general rule you shouldn't rely on graceful shutdown.



	return nil
}


//process input command signal (stop, quit, reopen, reload)
func processSignal(signal string) {
	var (
		pid           string
		commandAndPid = strings.Split(signal, "=")
	)
	if l := len(commandAndPid); l == 2 {
		pid = commandAndPid[1]
	} else if l > 2 {
		showUsageInfoNExit()
	}
	if err := message.ProcessSignal(message.Command(commandAndPid[0]), pid); err != nil {
		message.PrintNExit(err.Error())
	}
	os.Exit(0)
}


func configureTLS(opts *message.Options) {
	// If no trigger flags, ignore the others
	if !opts.TLS && !opts.TLSVerify {
		return
	}
	if opts.TLSCert == "" {
		message.PrintNExit("TLS Server certificate must be present and valid.")
	}
	if opts.TLSKey == "" {
		message.PrintNExit("TLS Server private key must be present and valid.")
	}

	tc := server.TLSConfigOpts{}
	tc.CertFile = opts.TLSCert
	tc.KeyFile = opts.TLSKey
	tc.CaFile = opts.TLSCaCert

	if opts.TLSVerify {
		tc.Verify = true
	}
	var err error
	if opts.TLSConfig, err = server.GenTLSConfig(&tc); err != nil {
		message.PrintNExit(err.Error())
	}
}

func configureClusterOpts(opts *server.Options) error {
	// If we don't have cluster defined in the configuration
	// file and no cluster listen string override, but we do
	// have a routes override, we need to report misconfiguration.
	if opts.Cluster.ListenStr == "" && opts.Cluster.Host == "" &&
		opts.Cluster.Port == 0 {
		if opts.RoutesStr != "" {
			message.PrintNExit("Solicited routes require cluster capabilities, e.g. --cluster.")
		}
		return nil
	}

	// If cluster flag override, process it
	if opts.Cluster.ListenStr != "" {
		clusterURL, err := url.Parse(opts.Cluster.ListenStr)
		if err != nil {
			return err
		}
		h, p, err := net.SplitHostPort(clusterURL.Host)
		if err != nil {
			return err
		}
		opts.Cluster.Host = h
		_, err = fmt.Sscan(p, &opts.Cluster.Port)
		if err != nil {
			return err
		}

		if clusterURL.User != nil {
			pass, hasPassword := clusterURL.User.Password()
			if !hasPassword {
				return fmt.Errorf("Expected cluster password to be set.")
			}
			opts.Cluster.Password = pass

			user := clusterURL.User.Username()
			opts.Cluster.Username = user
		} else {
			// Since we override from flag and there is no user/pwd, make
			// sure we clear what we may have gotten from config file.
			opts.Cluster.Username = ""
			opts.Cluster.Password = ""
		}
	}

	// If we have routes but no config file, fill in here.
	if opts.RoutesStr != "" && opts.Routes == nil {
		opts.Routes = server.RoutesFromStr(opts.RoutesStr)
	}

	return nil
}

func (p *program) Start() error {
	// The Start method must not block, or Windows may assume your service failed
	// to start. Launch a Goroutine here to do something interesting/blocking.
	opts := &message.Options{}
	flagSet := msgFlagSet(opts)
	flagSet.Parse(os.Args[1:])

	rand.Seed(time.Now().UTC().UnixNano())
	// Process args looking for non-flag options,
	// 'version' and 'help' only for now
	showVersion, showHelp, err := message.ProcessHelperCommands(flagSet)

	if err != nil {
		message.PrintNExit(err.Error() + tipsMsg)
	}else if showVersion {
		fmt.Println(message.VersionInfo("g-message"))
		os.Exit(0)
	}else if showHelp {
		showUsageInfoNExit()
	}
	//if showTLSHelp {
	//	server.PrintTLSHelpAndDie()
	//}
	if flagSet.Lookup("Debug").Value.(flag.Getter).Get().(bool) {
		opts.Trace = true
	}
	//signal
	signal := flagSet.Lookup("signal").Value.(flag.Getter).Get().(string)
	sl := flagSet.Lookup("sl").Value.(flag.Getter).Get().(string)
	if signal != "" {
		sl = signal
	}

	if  sl != "" {
		processSignal(sl)
	}
	configFile := flagSet.Lookup("configFile").Value.(flag.Getter).Get().(string)
	if configFile != "" {
		fileOpts, err := message.ProcessConfigFile(configFile)
		if err != nil {
			message.PrintNExit(err.Error())
		}
		opts = message.MergeOptions(fileOpts, opts)
	}

	// Remove any host/ip that points to itself in Route
	newroutes, err := message.RemoveSelfReference(opts.Cluster.Port, opts.Routes)
	if err != nil {
		message.PrintNExit(err.Error())
	}
	opts.Routes = newroutes

	// Configure TLS based on any present flags
	configureTLS(opts)

	// Configure cluster opts if explicitly set via flags.
	err = configureClusterOpts(opts)
	if err != nil {
		message.PrintNExit(err.Error())
	}

	// Create the server with appropriate options.
	msg := message.New(opts)

	// Configure the logger based on the flags
	msg.ConfigureLogger()

	msg.Start()
	p.msg = msg
	return nil
}