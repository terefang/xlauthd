package server

import (
	"xlauthd/pkg"

	"github.com/vjeantet/goldap/message"
	ldap "github.com/vjeantet/ldapserver"
)

func (xl *XlauthServerMain) HandleSearchDSE(w ldap.ResponseWriter, m *ldap.Message) {
	xl.LogInfo(m, "search request to DSE")
	e := ldap.NewSearchResultEntry("")
	e.AddAttribute("vendorName", "github.com/terefang")
	e.AddAttribute("vendorVersion", pkg.PkgVersion)
	e.AddAttribute("objectClass", "top", "extensibleObject")
	e.AddAttribute("supportedLDAPVersion", "3")
	e.AddAttribute("namingContexts", message.AttributeValue(xl.BaseDN),
		message.AttributeValue(UserOU+","+xl.BaseDN),
		message.AttributeValue(RoleOU+","+xl.BaseDN),
		message.AttributeValue(GroupOU+","+xl.BaseDN),
	)
	w.Write(e)
	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func (xl *XlauthServerMain) HandleSearchBase(w ldap.ResponseWriter, m *ldap.Message) {
	xl.LogInfo(m, "search request to %s", xl.BaseDN)
	e := ldap.NewSearchResultEntry(xl.BaseDN)
	e.AddAttribute("objectClass", "top", "dcObject")
	w.Write(e)
	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}
