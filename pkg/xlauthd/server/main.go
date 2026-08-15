package server

import (
    "fmt"
    "log/syslog"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/terefang/gommons/pkg/xcrypt"
    ldap "github.com/vjeantet/ldapserver"
)

type XlauthServerMain struct {
    BaseDN                string
    SystemAccountStamp    int64
    SystemAccountHtpasswd string
    SystemAccounts        map[string]string
    AccountStamp          int64
    AccountHtpasswd       string
    Accounts              map[string]string
    AccountCredentials    map[string]string
    AccountRoles          map[string]string
    Roles                 map[string]string
    RoleAccounts          map[string][]string
    ListenAddress         string
    TlsCert               string
    TlsKey                string
    IsBound               bool
    IsVerbose             bool
    Syslog                *syslog.Writer
    SyslogAddress         string
    Synchronized          sync.Mutex
}

func (xl *XlauthServerMain) GetHandler() ldap.Handler {
    xl.Synchronized.Lock()
    defer xl.Synchronized.Unlock()
    if xl.SyslogAddress != "" && xl.Syslog == nil {
        xl.Syslog, _ = syslog.Dial("udp", xl.SyslogAddress, syslog.LOG_LOCAL6, "xlauthd")
    }
    routes := ldap.NewRouteMux()

    if xl.SystemAccountHtpasswd != "" {
        _stat, err := os.Stat(xl.SystemAccountHtpasswd)
        if err == nil {
            if _stat.ModTime().Unix() > xl.SystemAccountStamp {
                xl.SystemAccountStamp = _stat.ModTime().Unix()
                xl.SystemAccounts, _, err = xcrypt.ReadFromHtpasswd(xl.SystemAccountHtpasswd, false)
                if err != nil {
                    xl.LogErr(nil, "read system accounts failed: %v", err)
                    if !xl.IsVerbose {
                        return routes
                    }
                }
            }
        }
    }

    if xl.AccountHtpasswd != "" {
        _stat, err := os.Stat(xl.AccountHtpasswd)
        if err == nil {
            if _stat.ModTime().Unix() > xl.AccountStamp {
                t1 := time.Now().UnixMicro()

                xl.AccountStamp = _stat.ModTime().Unix()
                xl.AccountCredentials, xl.AccountRoles, err = xcrypt.ReadFromHtpasswd(xl.AccountHtpasswd, false)
                if err != nil {
                    xl.LogErr(nil, "read account base failed: %v", err)
                    return routes
                }
                xl.RoleAccounts = make(map[string][]string)
                xl.Roles = make(map[string]string)
                for k, v := range xl.AccountRoles {
                    roles := strings.Split(v, ",")
                    for _, role := range roles {
                        role = strings.TrimSpace(role)
                        xl.RoleAccounts[role] = append(xl.RoleAccounts[role], k)
                        xl.Roles[strings.ToUpper(role)] = role
                    }
                }
                xl.Accounts = make(map[string]string)
                for k, _ := range xl.AccountCredentials {
                    xl.Accounts[strings.ToLower(k)] = k
                }

                t2 := time.Now().UnixMicro()
                td := t2 - t1
                xl.LogInfo(nil, "read %d entries in %d us.", len(xl.Accounts), td)
            }
        }
    }

    routes.Search(xl.HandleSearchUserBase).BaseDn(UserOU + "," + xl.BaseDN)
    routes.Search(xl.HandleSearchRoleBase).BaseDn(RoleOU + "," + xl.BaseDN)
    routes.Search(xl.HandleSearchBase).BaseDn(xl.BaseDN)
    routes.Search(xl.HandleSearchDSE).BaseDn("")
    routes.Search(xl.HandleSearchGet)
    routes.Bind(xl.HandleBind)
    return routes
}

func (xl *XlauthServerMain) HandleDeadEnd(w ldap.ResponseWriter, m *ldap.Message) {
    switch m.ProtocolOpType() {
    case ldap.ApplicationBindRequest:
        res := ldap.NewBindResponse(ldap.LDAPResultInvalidCredentials)
        res.SetDiagnosticMessage("Operation not implemented by server")
        w.Write(res)
    default:
        res := ldap.NewResponse(ldap.LDAPResultUnwillingToPerform)
        res.SetDiagnosticMessage("Operation not implemented by server")
        w.Write(res)
    }
}

func (xl *XlauthServerMain) MakeUserDn(uid string) string {
    return fmt.Sprintf("uid=%s,%s,%s", uid, UserOU, xl.BaseDN)
}

func (xl *XlauthServerMain) MakeRoleDn(rid string) string {
    return fmt.Sprintf("cn=%s,%s,%s", rid, RoleOU, xl.BaseDN)
}
