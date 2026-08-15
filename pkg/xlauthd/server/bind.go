package server

import (
    "strings"

    "github.com/terefang/gommons/pkg/xcrypt"
    "github.com/terefang/gommons/pkg/xstrings"
    ldap "github.com/vjeantet/ldapserver"
)

func (xl *XlauthServerMain) HandleBind(w ldap.ResponseWriter, m *ldap.Message) {
    xl.Synchronized.Lock()
    defer xl.Synchronized.Unlock()
    r := m.GetBindRequest()
    res := ldap.NewBindResponse(ldap.LDAPResultInvalidCredentials)
    if r.AuthenticationChoice() == "simple" {
        udn := string(r.Name())
        upw := string(r.AuthenticationSimple())
        xl.LogInfo(m, "bind request for %s", udn)
        if (len(udn) == 0) && (xl.SystemAccounts == nil || len(xl.SystemAccounts) == 0) {
            if xl.IsVerbose {
                xl.IsBound = true
                xl.LogWarning(m, "anonymous bind with no system accounts")
            }
            res.SetResultCode(ldap.LDAPResultSuccess)
            w.Write(res)
            return
        }

        // check user bind
        if strings.HasPrefix(udn, "uid=") {
            udn = udn[4:]
            ofs := xstrings.IndexOf(udn, ",", 0)
            if ofs > 0 {
                udn = udn[:ofs]
            }
            cudn, ok := xl.Accounts[strings.ToLower(udn)]
            if ok {
                _spw := xl.AccountCredentials[cudn]
                _ok, err := xcrypt.ValidateCryptedCredential(upw, _spw)
                if err == nil && _ok {
                    xl.IsBound = true
                    if xl.IsVerbose {
                        xl.LogInfo(m, "bind ok for user %s", udn)
                    }
                    res.SetResultCode(ldap.LDAPResultSuccess)
                    w.Write(res)
                    return
                }
            }
            xl.LogWarning(m, "invalid credentials for user %s", udn)
            res.SetResultCode(ldap.LDAPResultInvalidCredentials)
            res.SetDiagnosticMessage("invalid credentials")
            w.Write(res)
            return
        }

        //check system bind
        _spw, ok := xl.SystemAccounts[udn]
        if ok {
            _ok, err := xcrypt.ValidateCryptedCredential(upw, _spw)
            if err == nil && _ok {
                xl.IsBound = true
                if xl.IsVerbose {
                    xl.LogInfo(m, "bind ok for user %s", udn)
                }
                res.SetResultCode(ldap.LDAPResultSuccess)
                w.Write(res)
                return
            } else {
                xl.LogWarning(m, "invalid credentials for user %s", udn)
                res.SetResultCode(ldap.LDAPResultInvalidCredentials)
                res.SetDiagnosticMessage("invalid credentials")
                w.Write(res)
                return
            }
        }
        xl.LogWarning(m, "bind failed for user %s", udn)
        res.SetResultCode(ldap.LDAPResultInvalidCredentials)
        res.SetDiagnosticMessage("invalid credentials")
    } else {
        res.SetResultCode(ldap.LDAPResultUnwillingToPerform)
        res.SetDiagnosticMessage("Authentication choice not supported")
    }
    w.Write(res)
}
