package server

import (
    "fmt"
    "log"

    ldap "github.com/vjeantet/ldapserver"
)

func (xl *XlauthServerMain) LogEmerg(p *ldap.Message, f string, m ...any) error {
    if p != nil {
        f = fmt.Sprintf("%-16s %s", p.Client.GetConn().RemoteAddr().String(), f)
    }
    if xl.Syslog == nil {
        log.SetPrefix("EMERG  ")
        if len(m) == 0 {
            log.Println(f)
        } else {
            log.Printf(f, m...)
        }
        return nil
    } else {
        if len(m) == 0 {
            return xl.Syslog.Emerg(f)
        } else {
            return xl.Syslog.Emerg(fmt.Sprintf(f, m...))
        }
    }
}
func (xl *XlauthServerMain) LogAlert(p *ldap.Message, f string, m ...any) error {
    if p != nil {
        f = fmt.Sprintf("%-16s %s", p.Client.GetConn().RemoteAddr().String(), f)
    }
    if xl.Syslog == nil {
        log.SetPrefix("ALERT  ")
        if len(m) == 0 {
            log.Println(f)
        } else {
            log.Printf(f, m...)
        }
        return nil
    } else {
        if len(m) == 0 {
            return xl.Syslog.Alert(f)
        } else {
            return xl.Syslog.Alert(fmt.Sprintf(f, m...))
        }
    }
}

func (xl *XlauthServerMain) LogCrit(p *ldap.Message, f string, m ...any) error {
    if p != nil {
        f = fmt.Sprintf("%-16s %s", p.Client.GetConn().RemoteAddr().String(), f)
    }
    if xl.Syslog == nil {
        log.SetPrefix("CRIT   ")
        if len(m) == 0 {
            log.Println(f)
        } else {
            log.Printf(f, m...)
        }
        return nil
    } else {
        if len(m) == 0 {
            return xl.Syslog.Crit(f)
        } else {
            return xl.Syslog.Crit(fmt.Sprintf(f, m...))
        }
    }
}

func (xl *XlauthServerMain) LogErr(p *ldap.Message, f string, m ...any) error {
    if p != nil {
        f = fmt.Sprintf("%-16s %s", p.Client.GetConn().RemoteAddr().String(), f)
    }
    if xl.Syslog == nil {
        log.SetPrefix("ERROR  ")
        if len(m) == 0 {
            log.Println(f)
        } else {
            log.Printf(f, m...)
        }
        return nil
    } else {
        if len(m) == 0 {
            return xl.Syslog.Err(f)
        } else {
            return xl.Syslog.Err(fmt.Sprintf(f, m...))
        }
    }
}

func (xl *XlauthServerMain) LogWarning(p *ldap.Message, f string, m ...any) error {
    if p != nil {
        f = fmt.Sprintf("%-16s %s", p.Client.GetConn().RemoteAddr().String(), f)
    }
    if xl.Syslog == nil {
        log.SetPrefix("WARN   ")
        if len(m) == 0 {
            log.Println(f)
        } else {
            log.Printf(f, m...)
        }
        return nil
    } else {
        if len(m) == 0 {
            return xl.Syslog.Warning(f)
        } else {
            return xl.Syslog.Warning(fmt.Sprintf(f, m...))
        }
    }
}

func (xl *XlauthServerMain) LogNotice(p *ldap.Message, f string, m ...any) error {
    if p != nil {
        f = fmt.Sprintf("%-16s %s", p.Client.GetConn().RemoteAddr().String(), f)
    }
    if xl.Syslog == nil {
        log.SetPrefix("NOTICE ")
        if len(m) == 0 {
            log.Println(f)
        } else {
            log.Printf(f, m...)
        }
        return nil
    } else {
        if len(m) == 0 {
            return xl.Syslog.Notice(f)
        } else {
            return xl.Syslog.Notice(fmt.Sprintf(f, m...))
        }
    }
}

func (xl *XlauthServerMain) LogInfo(p *ldap.Message, f string, m ...any) error {
    if p != nil {
        f = fmt.Sprintf("%-16s %s", p.Client.GetConn().RemoteAddr().String(), f)
    }
    if xl.Syslog == nil {
        log.SetPrefix("INFO   ")
        if len(m) == 0 {
            log.Println(f)
        } else {
            log.Printf(f, m...)
        }
        return nil
    } else {
        if len(m) == 0 {
            return xl.Syslog.Info(f)
        } else {
            return xl.Syslog.Info(fmt.Sprintf(f, m...))
        }
    }
}

func (xl *XlauthServerMain) LogDebug(p *ldap.Message, f string, m ...any) error {
    if p != nil {
        f = fmt.Sprintf("%-16s %s", p.Client.GetConn().RemoteAddr().String(), f)
    }
    if xl.Syslog == nil {
        log.SetPrefix("DEBUG  ")
        if len(m) == 0 {
            log.Println(f)
        } else {
            log.Printf(f, m...)
        }
        return nil
    } else {
        if len(m) == 0 {
            return xl.Syslog.Debug(f)
        } else {
            return xl.Syslog.Debug(fmt.Sprintf(f, m...))
        }
    }
}
