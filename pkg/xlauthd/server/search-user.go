package server

import (
	"strings"

	"github.com/terefang/gommons/pkg/xstrings"
	"github.com/vjeantet/goldap/message"
	ldap "github.com/vjeantet/ldapserver"
)

func (xl *XlauthServerMain) HandleSearchUserBase(w ldap.ResponseWriter, m *ldap.Message) {
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
	if r.Scope() == ldap.SearchRequestScopeBaseObject {
		xl.LogInfo(m, "search/get request to %s,%s", UserOU, xl.BaseDN)
		e := ldap.NewSearchResultEntry(UserOU + "," + xl.BaseDN)
		e.AddAttribute("objectClass", "top", "organizationalUnit")
		e.AddAttribute("ou", UserRDN)
		w.Write(e)
		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
		w.Write(res)
		return
	}
	_filter := r.FilterString()
	_ifilter := strings.ToLower(_filter)
	xl.LogInfo(m, "search request to %s,%s with filter %s", UserOU, xl.BaseDN, _filter)
	if _ifilter == "(objectclass=*)" {
		xl.StreamAllAccounts(w)
		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
		w.Write(res)
		return
	}
	ofs := xstrings.IndexOf(_filter, "uid=", 0)
	if ofs > 0 {
		uid := _filter[ofs+4:]
		ofs = xstrings.IndexOf(uid, ")", 0)
		if ofs > 0 {
			uid = strings.ToLower(uid[:ofs])
			if uid == "*" {
				xl.StreamAllAccounts(w)
				res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
				w.Write(res)
				return
			}
			_, ok := xl.Accounts[uid]
			if ok {
				xl.StreamAccount(w, uid)
				res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
				w.Write(res)
				return
			}

			if xl.StreamMatchedAccounts(w, uid) > 0 {
				res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
				w.Write(res)
				return
			}
		}
	}
	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultNoSuchObject)
	w.Write(res)
}

func (xl *XlauthServerMain) StreamAllAccounts(w ldap.ResponseWriter) {
	for uid, _ := range xl.Accounts {
		xl.StreamAccount(w, uid)
	}
}

func (xl *XlauthServerMain) StreamMatchedAccounts(w ldap.ResponseWriter, match string) int {
	_ret := 0
	match = strings.ReplaceAll(match, "*", "")
	for uid, _ := range xl.Accounts {
		if match == "" {
			xl.StreamAccount(w, uid)
			_ret++
		} else if strings.Contains(strings.ToLower(uid), strings.ToLower(match)) {
			xl.StreamAccount(w, uid)
			_ret++
		}
	}
	return _ret
}

func (xl *XlauthServerMain) StreamAccount(w ldap.ResponseWriter, uid string) bool {
	uid = strings.ToLower(uid)
	cuid, ok := xl.Accounts[uid]
	if !ok {
		return false
	}
	e := ldap.NewSearchResultEntry(xl.MakeUserDn(cuid))
	e.AddAttribute("objectClass", "top", "uidObject")
	e.AddAttribute("uid", message.AttributeValue(cuid))
	_roles, ok := xl.AccountRoles[cuid]
	if ok {
		_rav := make([]message.AttributeValue, 0)
		for _, _role := range strings.Fields(_roles) {
			_role = strings.TrimSpace(_role)
			_role = strings.ToUpper(_role)
			_rav = append(_rav, message.AttributeValue(xl.MakeRoleDn(_role)))
		}
		e.AddAttribute("memberOf", _rav...)
	}
	w.Write(e)
	return true
}
