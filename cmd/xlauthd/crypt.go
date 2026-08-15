package main

import (
    "flag"
    "fmt"

    "github.com/skip2/go-qrcode"
    "github.com/terefang/gommons/pkg/xcrypt"

    "github.com/terefang/gommons/pkg/subcmd"
    "github.com/terefang/gommons/pkg/xtui"
)

type CryptCommand struct {
    doCrypt6   bool
    doApr1     bool
    doPrompt   bool
    password   string
    doTotp     bool
    doQrcode   bool
    doBcrypt   bool
    doWordpass bool
    doArgon2   bool
    doScrypt   bool
}

func (r *CryptCommand) Arguments(f *flag.FlagSet) {
    f.BoolVar(&r.doCrypt6, "crypt6", false, "crypt $6$ password")
    f.BoolVar(&r.doBcrypt, "bcrypt", false, "crypt $2b$ password")
    f.BoolVar(&r.doApr1, "apr1", false, "crypt $apr1$ password")
    f.BoolVar(&r.doScrypt, "scrypt", false, "crypt scrypt password")
    f.BoolVar(&r.doArgon2, "argon2", false, "crypt argon2 password")
    f.BoolVar(&r.doTotp, "totp", false, "generate totp fob")
    f.BoolVar(&r.doQrcode, "qrcode", false, "print qrcode for fob")
    f.StringVar(&r.password, "password", "", "dont generate password but use the one given")
    f.BoolVar(&r.doWordpass, "wordpass", false, "use wordpass algorithm")
    f.BoolVar(&r.doPrompt, "prompt", false, "prompt for a given password, more secure")
}

func (r CryptCommand) Info() (string, string) {
    return "crypt", "crypt a password"
}

func (r CryptCommand) Execute(args []string) int {

    if r.doPrompt {
        _pass, _err := xtui.ReadSecretVerifyString("Enter Password: ", "Re-Enter Password: ")
        if _err != nil {
            panic(_err)
        }
        r.password = _pass
    }

    if r.password == "" {
        if r.doWordpass {
            r.password = xcrypt.GenerateWordPass(32)
        } else {
            r.password = xcrypt.GeneratePasswordWithSym(xcrypt.PasswordSymbolSetExtensive, 32)
        }
        fmt.Printf("%s\n", r.password)
    }

    if r.doTotp {
        r.password = xcrypt.GeneratePasswordWithKdfSymLevel(r.password, xcrypt.PasswordSymbolSetExtensive, 32, 10)
        fmt.Println(xcrypt.TotpCredential(r.password))
        if r.doQrcode {
            _uri, err := xcrypt.TotpCredentialUrl(r.password)
            if err != nil {
                panic(err)
            }
            _qr, err := qrcode.New(_uri, qrcode.Low)
            if err != nil {
                panic(err)
            }
            fmt.Println(_qr.ToString(false))
        }
    } else if r.doScrypt {
        fmt.Println(xcrypt.GenerateScrypt(r.password))
    } else if r.doArgon2 {
        fmt.Println(xcrypt.GenerateArgon2(r.password))
    } else if r.doBcrypt {
        fmt.Println(xcrypt.GenerateBcrypt(r.password))
    } else if r.doCrypt6 {
        fmt.Println(xcrypt.GenerateSha512Crypt(r.password))
    } else if r.doApr1 {
        fmt.Println(xcrypt.CryptApr1Credential(r.password))
    } else {
        fmt.Println("{type7}" + xcrypt.Type7_encrypt(r.password))
    }
    return 0
}

func init() {
    subcmd.Register(&CryptCommand{})
}
