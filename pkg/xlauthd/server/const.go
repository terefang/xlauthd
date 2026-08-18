package server

import (
	"fmt"
	"strings"
)

const UserOU = "ou=users"
const RoleOU = "ou=roles"
const GroupOU = "ou=groups"

const UserRDN = "users"
const RoleRDN = "roles"
const GroupRDN = "groups"

func (xl *XlauthServerMain) MakeUserDn(uid string) string {
	return fmt.Sprintf("uid=%s,%s,%s", strings.ToLower(uid), UserOU, xl.BaseDN)
}

func (xl *XlauthServerMain) MakeRoleDn(rid string) string {
	return fmt.Sprintf("cn=%s,%s,%s", strings.ToUpper(rid), RoleOU, xl.BaseDN)
}

func (xl *XlauthServerMain) MakeGroupDn(rid string) string {
	return fmt.Sprintf("cn=%s,%s,%s", strings.ToUpper(rid), GroupOU, xl.BaseDN)
}
