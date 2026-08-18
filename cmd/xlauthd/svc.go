package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"xlauthd/pkg"
	"xlauthd/pkg/xlauthd/server"

	"github.com/terefang/gommons/pkg/subcmd"
	"github.com/terefang/gommons/pkg/xini"
	ldap "github.com/vjeantet/ldapserver"
)

type xLogger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type nullLogger struct{}

func (n nullLogger) Fatal(v ...interface{}) {}

func (n nullLogger) Fatalf(format string, v ...interface{}) {}

func (n nullLogger) Fatalln(v ...interface{}) {}

func (n nullLogger) Panic(v ...interface{}) {}

func (n nullLogger) Panicf(format string, v ...interface{}) {}

func (n nullLogger) Panicln(v ...interface{}) {}

func (n nullLogger) Print(v ...interface{}) {}

func (n nullLogger) Printf(format string, v ...interface{}) {}

func (n nullLogger) Println(v ...interface{}) {}

func init() {
	subcmd.Register(&SvcCommand{})
}

type SvcCommand struct {
	Listen   string `ini:"listen"`
	BaseDN   string `ini:"basedn"`
	UserDB   string `ini:"userdb"`
	SystemDB string `ini:"systemdb"`
	GroupDB  string `ini:"groupdb"`
	TlsCert  string `ini:"tlscert"`
	TlsKey   string `ini:"tlskey"`
	Verbose  bool   `ini:"verbose"`
	Syslog   string `ini:"syslog"`
	config   string
}

func (r *SvcCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&r.Listen, "listen", "127.0.0.1:5389", "listen address and port")
	f.StringVar(&r.BaseDN, "basedn", "dc=eample,dc=net", "path prefix")
	f.StringVar(&r.UserDB, "userdb", "", "user password file (.htpasswd style)")
	f.StringVar(&r.GroupDB, "groupdb", "", "group file (.htgroup style)")
	f.StringVar(&r.SystemDB, "systemdb", "", "system password file (.htpasswd style)")
	f.StringVar(&r.TlsCert, "tlscert", "", "tls certificate (.pem)")
	f.StringVar(&r.TlsKey, "tlskey", "", "tls key (.pem)")
	f.BoolVar(&r.Verbose, "verbose", false, "verbose output")
	f.StringVar(&r.Syslog, "syslog", "", "syslog address and port")
	f.StringVar(&r.config, "config", "", "read config parameters from file")
}

func (r SvcCommand) Info() (string, string) {
	return "server", `runs the server`
}

func (r *SvcCommand) Execute(args []string) int {
	fmt.Println(Banner)
	fmt.Println(pkg.PkgVersion)

	svc := &server.XlauthServerMain{}
	if r.config != "" {
		_ini, err := xini.NewIniConfig(r.config)
		if err != nil {
			panic(err)
		}
		_ini.Unmarshal(xini.GLOBAL_SECTION, r)
	}

	svc.SystemAccountHtpasswd = r.SystemDB
	svc.AccountHtpasswd = r.UserDB
	svc.GroupHtgroup = r.GroupDB
	svc.BaseDN = r.BaseDN
	svc.IsVerbose = r.Verbose
	svc.SyslogAddress = r.Syslog

	//ldap logger
	ldap.Logger = &nullLogger{}

	log.Println("setting basedn", svc.BaseDN)
	log.Println("setting userdn", server.UserOU, svc.BaseDN)
	log.Println("setting roledn", server.RoleOU, svc.BaseDN)
	log.Println("setting groupdn", server.GroupOU, svc.BaseDN)
	//Create a new LDAP Server
	server := ldap.NewServerWithHandlerSource(svc)

	if r.Verbose {
		log.Println("start listen on", r.Listen)
	}
	if r.TlsCert != "" && r.TlsKey != "" {
		//SSL/TLS
		if r.Verbose {
			log.Println("enabled tls")
		}
		secureConn := func(s *ldap.Server) {
			config, _ := r.GetTLSconfig()
			s.Listener = tls.NewListener(s.Listener, config)
		}
		go server.ListenAndServe(r.Listen, secureConn)
	} else {
		go server.ListenAndServe(r.Listen)
	}

	// When CTRL+C, SIGINT and SIGTERM signal occurs
	// Then stop server gracefully
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	close(ch)

	server.Stop()
	return 0
}

func (r *SvcCommand) GetTLSconfig() (*tls.Config, error) {
	hostCert, err := os.ReadFile(r.TlsCert)
	if err != nil {
		log.Println(err.Error())
		return &tls.Config{}, err
	}
	hostKey, err := os.ReadFile(r.TlsKey)
	if err != nil {
		log.Println(err.Error())
		return &tls.Config{}, err
	}
	cert, err := tls.X509KeyPair(hostCert, hostKey)
	if err != nil {
		log.Println(err.Error())
		return &tls.Config{}, err
	}

	return &tls.Config{
		MinVersion: tls.VersionTLS10,
		MaxVersion: tls.VersionTLS13,
		GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			return &cert, nil
		},
		//ServerName:   cert.Leaf.Subject.CommonName,
	}, nil
}
