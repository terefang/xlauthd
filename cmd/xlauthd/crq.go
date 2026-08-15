package main

import (
    "flag"
    "fmt"
    "os"
    "strings"

    "github.com/terefang/gommons/pkg/certs"
    "github.com/terefang/gommons/pkg/subcmd"
    "github.com/terefang/gommons/pkg/xfile"
)

func init() {
    subcmd.Register(&GenCrqCommand{})
}

type GenCrqCommand struct {
    keyfile  string
    crqfile  string
    doServer bool
    doClient bool
    useDN    string
    useDNS   string
    useBits  int
}

func (r *GenCrqCommand) Arguments(f *flag.FlagSet) {
    f.StringVar(&r.keyfile, "key", "", "key-file")
    f.StringVar(&r.crqfile, "req", "", "crq-file")
    f.StringVar(&r.useDN, "dn", "", "certificate dn")
    f.StringVar(&r.useDNS, "dns", "", "dns san")
    f.IntVar(&r.useBits, "bits", 2048, "key bits")
    f.BoolVar(&r.doServer, "server", true, "tls-server permission")
    f.BoolVar(&r.doClient, "client", true, "tls-client permission")
}

func (r GenCrqCommand) Info() (string, string) {
    return "crq", `generate certificate request+key`
}

func (r GenCrqCommand) Execute(args []string) int {

    if len(args) > 0 {
        for _, arg := range args {
            if strings.HasSuffix(arg, ".key") {
                r.keyfile = arg
            }
        }
    }

    if r.keyfile == "" {
        r.Usage()
        return -1
    }

    if !xfile.FileExists(r.keyfile) {
        certs.MakeRsaKeyFile(r.useBits, r.keyfile)
    }

    sAN := make([]string, 0)
    if r.useDNS != "" {
        _d := strings.Split(r.useDNS, ",")
        sAN = append(sAN, _d...)
    }
    err := certs.MakeRsaCrqFile(r.keyfile, r.useDN, sAN, r.doServer, r.doClient, r.crqfile)
    if err != nil {
        panic(err)
    }
    return 0
}

func (r GenCrqCommand) Usage() {
    fmt.Fprintln(os.Stderr, "usage: crq [flags]")
}
