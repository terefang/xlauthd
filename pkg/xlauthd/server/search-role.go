package server

import (
	"strings"

	"github.com/terefang/gommons/pkg/xstrings"
	"github.com/vjeantet/goldap/message"
	ldap "github.com/vjeantet/ldapserver"
)

func (xl *XlauthServerMain) HandleSearchRoleBase(w ldap.ResponseWriter, m *ldap.Message) {
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
		xl.LogInfo(m, "search/get request to %s,%s", RoleOU, xl.BaseDN)
		e := ldap.NewSearchResultEntry(RoleOU + "," + xl.BaseDN)
		e.AddAttribute("objectClass", "top", "organizationalUnit")
		e.AddAttribute("ou", RoleRDN)
		w.Write(e)
		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
		w.Write(res)
		return
	}
	_filter := r.FilterString()
	_ifilter := strings.ToLower(_filter)
	xl.LogInfo(m, "search request to %s,%s with filter %s", RoleOU, xl.BaseDN, _filter)
	if _ifilter == "(objectclass=*)" {
		xl.StreamAllRoles(w)
		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
		w.Write(res)
		return
	}
	ofs := xstrings.IndexOf(_ifilter, "cn=", 0)
	if ofs > 0 {
		rid := _filter[ofs+3:]
		ofs = xstrings.IndexOf(rid, ")", 0)
		if ofs > 0 {
			rid = rid[:ofs]
			if rid == "*" {
				xl.StreamAllRoles(w)
				res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
				w.Write(res)
				return
			}
			_, ok := xl.RoleAccounts[rid]
			if ok {
				xl.StreamRole(w, rid)
				return
			}

			if xl.StreamMatchedRoles(w, rid) > 0 {
				res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
				w.Write(res)
				return
			}
		}
	}
	ofs = xstrings.IndexOf(_ifilter, "=uid=", 0)
	if ofs > 0 {
		uid := _filter[ofs+5:]
		ofs = xstrings.IndexOf(uid, ")", 0)
		if ofs > 0 {
			uid = uid[:ofs]
			if uid == "*" {
				xl.StreamAllRoles(w)
				res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
				w.Write(res)
				return
			}
			ofs = xstrings.IndexOf(uid, ",", 0)
			if ofs > 0 {
				uid = uid[:ofs]
			}
			uid = strings.ToLower(uid)
			cuid, ok := xl.Accounts[uid]
			if ok {
				roles, ok := xl.AccountRoles[cuid]
				if ok {
					for _, role := range strings.Split(roles, ",") {
						xl.StreamRole(w, role)
					}
					res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
					w.Write(res)
					return
				}
			}
		}
	}

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultNoSuchObject)
	w.Write(res)
}

func (xl *XlauthServerMain) StreamAllRoles(w ldap.ResponseWriter) {
	for rid, _ := range xl.RoleAccounts {
		xl.StreamRole(w, rid)
	}
}

func (xl *XlauthServerMain) StreamMatchedRoles(w ldap.ResponseWriter, match string) int {
	_ret := 0
	match = strings.ReplaceAll(match, "*", "")
	for rid, _ := range xl.RoleAccounts {
		if match == "" {
			xl.StreamRole(w, rid)
			_ret++
		} else if strings.Contains(strings.ToLower(rid), strings.ToLower(match)) {
			xl.StreamRole(w, rid)
			_ret++
		}
	}
	return _ret
}

func (xl *XlauthServerMain) StreamRole(w ldap.ResponseWriter, rid string) bool {
	rid = strings.ToUpper(rid)
	crid, ok := xl.Roles[rid]
	if !ok {
		return false
	}
	_uids, ok := xl.RoleAccounts[crid]
	if !ok {
		return false
	}
	e := ldap.NewSearchResultEntry(xl.MakeRoleDn(rid))
	e.AddAttribute("objectClass", "top", "cnObject")
	e.AddAttribute("cn", message.AttributeValue(rid))
	_rav := make([]message.AttributeValue, 0)
	for _, _uid := range _uids {
		_rav = append(_rav, message.AttributeValue(xl.MakeUserDn(_uid)))
	}
	e.AddAttribute("uniqueMember", _rav...)
	w.Write(e)
	return true
}
