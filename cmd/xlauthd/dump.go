package main

import (
	"fmt"

	"github.com/terefang/gommons/pkg/subcmd"
)

type DumpHtpCommand struct {
	subcmd.NoFlags
}

func init() {
	subcmd.Register(&DumpHtpCommand{})
}

func (r DumpHtpCommand) Info() (string, string) {
	return "dump", "print sample config files"
}

func (r DumpHtpCommand) Execute(args []string) int {
	fmt.Println(`
## --- SAMPLE HTPASSWD
/ --- comment
; --- comment
scott:$1$abcdefgh$G//4keteveJp0qb8z2DxG/:role1,role2
# --- comment
% --- comment
tiger:$apr1$VomI8RgV$JR59eyvesgXIeOQjOvvnP1:role0
! --- comment
$ --- comment
foobar:{TOTP}IZATU4DVPB7EMUJMO4STSTZJGE7HKNSIFM4HOODFMNGFWMSTG43A:role98
END
- everything below /END/ is ignored


## --- SAMPLE HTGROUP
/ --- comment
; --- comment
group1:scott tiger
# --- comment
% --- comment
group2:tiger
! --- comment
$ --- comment
group99:foo bar foobar baz
END
- everything below /END/ is ignored


## --- SAMPLE SERVER CONFIG
listen = 0.0.0.0:3890
users = /etc/xlauthd/users.htpasswd
listen = 127.0.0.1:5389
basedn = "dc=example,dc=net"
userdb = /path/to/users.htpasswd
systemdb = /path/to/system.htpasswd
tlscert = /path/to/tls.crt.pem
tlskey = /path/to/tls.key.pem
# verbose = TRUE
# syslog = 127.0.0.1:514
`)
	return 0
}
