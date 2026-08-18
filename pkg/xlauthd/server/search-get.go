package server

import (
	"strings"

	"github.com/terefang/gommons/pkg/xstrings"
	ldap "github.com/vjeantet/ldapserver"
)

func (xl *XlauthServerMain) HandleSearchGet(w ldap.ResponseWriter, m *ldap.Message) {
	xl.Synchronized.Lock()
	defer xl.Synchronized.Unlock()
	if !xl.IsBound {
		if xl.IsVerbose {
			xl.LogWarning(m, "search request without bind from %s", m.Client.GetConn().RemoteAddr().String())
		}
		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultInsufficientAccessRights)
		w.Write(res)
		return
	}

	r := m.GetSearchRequest()
	xdn := string(r.BaseObject())
	xl.LogInfo(m, "search/get request to %s", xdn)
	done := false
	if strings.HasPrefix(xdn, "uid=") {
		ofs := xstrings.IndexOf(xdn, ",", 0)
		if ofs > 0 {
			xdn = xdn[4:ofs]
		} else {
			xdn = xdn[4:]
		}
		done = xl.StreamAccount(w, xdn)
	} else if strings.HasPrefix(xdn, "cn=") {
		ofs := xstrings.IndexOf(xdn, ",", 0)
		if ofs > 0 {
			xdn = xdn[3:ofs]
		} else {
			xdn = xdn[3:]
		}

		if strings.Contains(xdn, RoleOU) {
			done = xl.StreamRole(w, xdn)
		} else if strings.Contains(xdn, GroupOU) {
			done = xl.StreamGroup(w, xdn)
		}
	}
	if done {
		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
		w.Write(res)
	} else {
		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultNoSuchObject)
		w.Write(res)
	}
}
