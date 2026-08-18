package server

import (
	"strings"

	"github.com/terefang/gommons/pkg/xmatch"
	"github.com/terefang/gommons/pkg/xstrings"
	"github.com/vjeantet/goldap/message"
	ldap "github.com/vjeantet/ldapserver"
)

func (xl *XlauthServerMain) HandleSearchGroupBase(w ldap.ResponseWriter, m *ldap.Message) {
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
		xl.LogInfo(m, "search/get request to %s,%s", GroupOU, xl.BaseDN)
		e := ldap.NewSearchResultEntry(GroupOU + "," + xl.BaseDN)
		e.AddAttribute("objectClass", "top", "organizationalUnit")
		e.AddAttribute("ou", GroupRDN)
		w.Write(e)
		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
		w.Write(res)
		return
	}
	_filter := r.FilterString()
	_ifilter := strings.ToLower(_filter)
	xl.LogInfo(m, "search request to %s,%s with filter %s", GroupOU, xl.BaseDN, _filter)
	if _ifilter == "(objectclass=*)" {
		xl.StreamAllGroups(w)
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
				xl.StreamAllGroups(w)
				res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
				w.Write(res)
				return
			}
			_, ok := xl.GroupAccounts[rid]
			if ok {
				xl.StreamGroup(w, rid)
				res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
				w.Write(res)
				return
			}

			if xl.StreamMatchedGroups(w, rid) > 0 {
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
		}
		if uid == "*" {
			xl.StreamAllGroups(w)
			res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
			w.Write(res)
			return
		}
		ofs = xstrings.IndexOf(uid, ",", 0)
		if ofs > 0 {
			uid = uid[:ofs]
		}
		uid = strings.ToLower(uid)
		if xl.StreamMatchedGroupsByUser(w, uid) > 0 {
			res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
			w.Write(res)
			return
		}
	}

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultNoSuchObject)
	w.Write(res)
}

func (xl *XlauthServerMain) StreamAllGroups(w ldap.ResponseWriter) {
	for rid, _ := range xl.GroupAccounts {
		xl.StreamGroup(w, rid)
	}
}

func (xl *XlauthServerMain) StreamMatchedGroupsByUser(w ldap.ResponseWriter, uid string) int {
	_ret := 0
	uid = strings.ReplaceAll(uid, "*", "")
	roles, err := xmatch.MatchSimpleContainsVforK(xl.GroupAccounts, uid)
	if err == nil {
		for _, role := range roles {
			xl.StreamGroup(w, role)
			_ret++
		}
	}
	return _ret
}

func (xl *XlauthServerMain) StreamMatchedGroups(w ldap.ResponseWriter, match string) int {
	_ret := 0
	match = strings.ReplaceAll(match, "*", "")
	for rid, _ := range xl.GroupAccounts {
		if match == "" {
			xl.StreamGroup(w, rid)
			_ret++
		} else if strings.Contains(strings.ToLower(rid), strings.ToLower(match)) {
			xl.StreamGroup(w, rid)
			_ret++
		}
	}
	return _ret
}

func (xl *XlauthServerMain) StreamGroup(w ldap.ResponseWriter, rid string) bool {
	rid = strings.ToLower(rid)
	_uids, ok := xl.GroupAccounts[rid]
	if !ok {
		return false
	}
	e := ldap.NewSearchResultEntry(xl.MakeGroupDn(rid))
	e.AddAttribute("objectClass", "top", "cnObject")
	e.AddAttribute("cn", message.AttributeValue(strings.ToUpper(rid)))
	_rav := make([]message.AttributeValue, 0)
	for _, _uid := range _uids {
		_rav = append(_rav, message.AttributeValue(xl.MakeUserDn(_uid)))
	}
	e.AddAttribute("uniqueMember", _rav...)
	w.Write(e)
	return true
}
