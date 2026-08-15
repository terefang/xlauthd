# xlauthd

**xlauthd — extremely lightweight LDAP authentication daemon**

## NAME

**xlauthd** — a minimal LDAP-compatible authentication daemon backed by 
htpasswd-style user and role files.

## SYNOPSIS

```text
xlauthd [OPTIONS]
```

## DESCRIPTION

**xlauthd** is an extremely lightweight authentication daemon written in Go. 
It implements the **bare minimum of the LDAP protocol** required to provide 
simple LDAP-compatible authentication without the complexity and resource 
requirements of a full LDAP directory server.

The daemon exposes a small LDAP server interface on top of these files, 
allowing applications and existing LDAP-aware software to perform authentication 
using standard LDAP operations.

xlauthd is **not intended to be a general-purpose LDAP directory server**. 
It deliberately implements only the LDAP functionality required for authentication 
and role lookup.

The design goal is simple:

> Provide just enough LDAP to authenticate users and determine their roles.

## FEATURES

* Written in Go
* Extremely small runtime footprint
* Single self-contained executable
* No LDAP database
* No external directory service required
* LDAP-compatible authentication
* htpasswd-style credential storage
* Separate user and system files
* Hot-reloadable credential files
* Minimal configuration
* Suitable for containers and minimal Linux systems
* Designed for service-to-service authentication
* No need to deploy a full OpenLDAP installation

## LDAP PROTOCOL

xlauthd implements only the subset of LDAP required for authentication.

The implementation is intentionally limited rather than attempting to 
provide a complete LDAP directory.

Supported functionality is centered around:

* LDAP bind
* User authentication
* Basic search operations required by supported clients
* User lookup
* Role/group lookup where applicable
* LDAP result codes required for successful and failed authentication
* Basic connection handling

Operations that are not required for authentication are intentionally omitted.

xlauthd should therefore be considered an **LDAP authentication endpoint**, not an LDAP directory.

## AUTHENTICATION MODEL

Authentication is backed by an htpasswd-style user file.

Each entry contains a username and password verifier and role infromation.

## USER FILE

The user file follows the traditional htpasswd-style format:

```text
username:password:role1,...,roleN
```

Blank lines and supported comment lines may be ignored.

you can use `xlathd dump` to get an example htpasswd file with all supported features.

The user file should be readable only by the account running xlauthd and system administrators:

```sh
chmod 0600 /etc/xlauthd/users.htpasswd
```

## LDAP AUTHENTICATION

A client authenticates using a normal LDAP bind.

Conceptually:

```text
LDAP client
    |
    | Bind(username, password)
    v
+-----------+
|  xlauthd  |
+-----------+
    |
    +-- system.htpasswd
    |
    +-- users.htpasswd
```

xlauthd validates the supplied credentials against the configured user file(s).

If authentication succeeds, the LDAP bind succeeds.

If authentication fails, xlauthd returns the appropriate LDAP authentication failure result.

The client does not need to know that the underlying credential store is an htpasswd file.

The `system.htpasswd` file is used for authenticating access to services without 
adding service accounts to the application's `user.htpasswd`. This keeps system, 
machine-to-machine, and service credentials separate from normal user accounts, 
avoiding unnecessary entries in the application's user store while still allowing 
services to authenticate using the standard htpasswd-style credential format.

## HOT RELOAD

Credential file support **hot reloading**.

Changes to the htpasswd file are detected without requiring xlauthd to be restarted.

This allows administrators to:

* Add users
* Remove users
* Change passwords
* Change roles
* Rotate credentials

without interrupting the LDAP service.

## CONFIGURATION

The configuration is intentionally small. xlauthd does not require 
configuration of LDAP schemas, databases, indexes, replication, 
overlays, or directory backends.

## EXAMPLE DEPLOYMENT

Create the configuration directory:

```sh
mkdir -p /etc/xlauthd
```

Create a user file:

```text
scott:$apr1$...:readonly
tiger:$6$...:users,admin
foo:$2y$...:monitoring
```

Configure xlauthd:

```ini
listen = 127.0.0.1:5389
basedn = "dc=example,dc=net"
userdb = /etc/xlauthd/users.htpasswd
systemdb = /etc/xlauthd/system.htpasswd
# tlscert = /etc/xlauthd/tls.crt.pem
# tlskey = /etc/xlauthd/tls.key.pem
# verbose = TRUE
# syslog = 127.0.0.1:514
```

Start the daemon:

```sh
xlauthd server -config /etc/xlauthd/xlauthd.conf
```

An LDAP-aware application can then be configured to authenticate against:

```text
ldap://server.example.com:5389
```

## LDAP SEARCH

TBD.

## SECURITY

xlauthd should normally be placed behind a trusted network boundary or 
used with an appropriate secure transport.

The minimal LDAP implementation does not by itself make an unencrypted 
LDAP connection secure.

Where passwords are transmitted over the network, use an appropriate TLS 
configuration or place xlauthd behind a trusted secure transport.

Credential files should have restrictive filesystem permissions:

```sh
chmod 0600 /etc/xlauthd/users.htpasswd
```

The daemon should run as a dedicated unprivileged user whenever possible.

Logs must not contain plaintext passwords, password hashes, API secrets, 
or other credential material.

## WHY NOT OPENLDAP?

If you need an LDAP directory, use a real LDAP directory server.

xlauthd exists for the cases where deploying a full directory service is unnecessary.

For example, an application may only need:

```text
LDAP bind
    |
    +-- Is the username/password valid?
    |
    +-- What roles does this user have?
```

Deploying an entire directory database and LDAP server for that use 
case may introduce unnecessary operational complexity.

xlauthd reduces the problem to:

```text
htpasswd file
       |
       v
    xlauthd
       |
       v
  LDAP clients
```

## DESIGN GOALS

### Extremely small

The implementation intentionally contains only the LDAP functionality 
required by its authentication use case.

### File-based

Users and roles are stored in a simple text file rather than a database.

### Stateless

Authentication does not require maintaining a directory database or 
persistent application state.

### Easy to deploy

The preferred deployment model is a single Go binary and a few 
configuration files.

### Compatible

Existing LDAP-aware applications can use xlauthd without implementing 
a custom authentication protocol.

### Auditable

A small protocol implementation and simple credential backend are easier 
to inspect and reason about than a full directory platform.

## LIMITATIONS

xlauthd is deliberately **not a complete LDAP implementation**.

Do not use it when you require:

* A full LDAP directory
* Directory write operations
* LDAP replication
* Active Directory compatibility
* Complex LDAP schemas
* LDAP groups with complex membership semantics
* Directory-wide ACLs
* LDAP synchronization
* High-volume directory search
* General-purpose identity management

Use xlauthd when the requirement is closer to:

> "I need my existing LDAP-aware application to authenticate against a small local credential store."

## BUILDING

xlauthd is written in Go.

Build from source:

```sh
git clone https://github.com/terefang/xlauthd.git
cd xlauthd

go build -o xlauthd .
```

Run the tests:

```sh
go test ./...
```

For a minimal Linux binary:

```sh
just build
```

## INSTALLATION

Install the binary:

```sh
install -m 0755 dist/xlauthd /usr/local/bin/xlauthd
```

Create the configuration directory:

```sh
mkdir -p /etc/xlauthd
```

Install:

```text
/etc/xlauthd/xlauthd.conf
/etc/xlauthd/users.htpasswd
```

Start xlauthd using the system's preferred service manager.

## SYSTEMD

A minimal systemd unit might look like:

```ini
[Unit]
Description=xlauthd lightweight LDAP authentication daemon
After=network.target

[Service]
Type=simple
User=xlauthd
Group=xlauthd
ExecStart=/usr/local/bin/xlauthd server -config /etc/xlauthd/xlauthd.conf
Restart=on-failure
TimeoutStopSec=5s
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

Enable the service:

```sh
systemctl enable --now xlauthd
```

## CREATING PASSWORD HASHES

xlauthd includes a built-in `crypt` command for generating password hashes suitable for use in an htpasswd-style password file.

Display the available options with:

```sh
xlauthd help crypt
```

The command supports several password-hashing algorithms, including:

* `-apr1` — Apache MD5 (`$apr1$`) format.
* `-bcrypt` — bcrypt (`$2b$`) format.
* `-crypt6` — SHA-512 crypt (`$6$`) format.
* `-argon2` — Argon2.
* `-scrypt` — scrypt.

### GENERATE A RANDOM PASSWORD HASH

If neither `-prompt` nor `-password` is specified, `xlauthd crypt` automatically 
generates a **random 32-character password** and hashes it.

For example:

```sh
xlauthd crypt -bcrypt

$2b$14$DIcZ.Po7DCEMHnG4d0JLuuF5IuCueJcrENsDSBSNdvlhY5zafqF5O
```

This is useful when creating credentials for services or other machine-to-machine 
accounts where the password does not need to be chosen by a human.

The generated password and its hash are displayed by the command. Store the generated 
password securely if it will be required by the client, and place only the resulting 
hash in the htpasswd file.

### HUMAN-FRIENDLY RANDOM PASSWORDS

The `-wordpass` option generates a random password intended for **human consumption**. 
Rather than producing a completely arbitrary string, it generates a password using a 
word-based format that is easier to read, communicate, and enter manually.

For example:

```sh
xlauthd crypt -wordpass

Giffard5!Alar5.Voidmere6%Vermandois5_
```

This is useful when a password needs to be manually entered by an administrator or 
user while still being randomly generated.

### GENERATE A HASH INTERACTIVELY

For normal use, prompt for the password rather than placing it directly on the command line.

The command prompts for the password and prints the resulting password hash.

For example:

```text
$ xlauthd crypt -prompt -bcrypt
Enter Password: *****
Re-Enter Password: *****
$2b$14$duTYXlUcxWPkQrXsTX65hu3ktIqu9dikgI0mqj5mWHkh.cg1ygDSS
```

The generated hash can then be placed in the htpasswd file:

```text
scott:$2b$...:role1,role2
```

### GENERATE AN APR1 HASH

To generate an Apache MD5 (`$apr1$`) hash:

```sh
xlauthd crypt -prompt -apr1
```

Example:

```text
$apr1$Dw3u/NCE$4v0e8YfgVW5GdPTKAtboW1
```

### GENERATE A SHA-512 CRYPT HASH

To generate a traditional `$6$` SHA-512 crypt hash:

```sh
xlauthd crypt -prompt -crypt6
```

Example:

```text
$6$6YAk4Rhd3gdq0BXz$OqlQ51Ri2l/HfXv8i/V5u4JQDW/TYaZnz.WJCchsx2iKY6MshXMIERclaq/Lyb.YtDHZKOx5pc5yNZOTf7/UP1
```

### GENERATE A HASH FROM THE COMMAND LINE

A password can also be supplied directly using `-password`:

```sh
xlauthd crypt -password 'secret-password' -bcrypt
```

However, this is **not recommended** for real credentials because 
command-line arguments may be visible to other users through process 
inspection or shell history.

Prefer:

```sh
xlauthd crypt -prompt -bcrypt
```

### ADDING A USER

Once a hash has been generated, add the username and hash to the htpasswd file using the format:

```text
username:password-hash:role1,...,roleN
```

For example:

```text
scott:$2b$...:role1,role2
tiger:$6$...:role99
foo:$apr1$...:role66
```

Each user occupies one line.

### RECOMMENDED WORKFLOW

For a new user, generate the hash interactively:

```sh
xlauthd crypt -prompt -bcrypt
```

Then add the resulting hash to:

```text
/etc/xlauthd/users.htpasswd
```

For example:

```text
scott:$2b$12$...
```

After saving the file, xlauthd will detect the updated password file during 
its normal hot-reload cycle, so no daemon restart is required.

### SECURITY

Use `-prompt` whenever possible. Avoid passing production passwords through 
`-password`, shell scripts, environment variables, or command-line arguments.

The htpasswd file contains password verifiers and should be protected accordingly:

```sh
chmod 0600 /etc/xlauthd/users.htpasswd
```

The password itself should never be stored in the file; only the generated hash should be present.

## CONFIGURING TLS

xlauthd can provide TLS-secured connections by configuring a certificate and private key:

```ini
tlscert = /etc/xlauthd/tls.crt
tlskey = /etc/xlauthd/tls.key
```

The certificate and key should be stored with appropriate filesystem permissions. 
The private key should only be readable by the account running xlauthd.

### CREATE A CERTIFICATE REQUEST

Use the built-in `crq` command to generate a private key and certificate signing request (CSR):

```sh
xlauthd crq -bits 2048 \
    -key tls.key \
    -req tls.req \
    -dns server.lan \
    -dn 'cn=server.lan,dc=acme'
```

This creates:

* `tls.key` — the private key in PEM format.
* `tls.req` — the certificate signing request in PEM format.

Keep the private key secure and do not send it to the certificate authority.

### SIGN THE REQUEST

Copy the generated `tls.req` file to your certificate authority (CA) and have it signed.

The CA should return a signed certificate in PEM format, for example:

```text
tls.crt
```

Copy the signed certificate and private key to the locations configured above:

```text
/etc/xlauthd/tls.crt
/etc/xlauthd/tls.key
```

The certificate should contain the appropriate DNS name(s) used by clients to connect to xlauthd.

### VERIFYING THE CERTIFICATE

After the certificate has been signed by your certificate authority, verify that 
it is valid and contains the expected subject and DNS names before deploying it to xlauthd.

#### VERIFY WITH OPENSSL

Display the certificate details:

```sh
openssl x509 -in tls.crt -text -noout
```

#### VERIFY WITH CERTTOOL

If using GnuTLS, the equivalent certificate inspection can be performed with `certtool`:

```sh
certtool -i < tls.crt
```

#### FINAL CHECK

Before installing the certificate, verify that:

1. The certificate is signed by the expected CA.
2. The certificate is currently within its validity period.
3. The expected DNS name is present in the **Subject Alternative Name** extension.
4. The certificate's public key matches the generated private key.
5. The certificate is intended for server authentication.

Once these checks pass, install the certificate and key in the locations configured 
by `tlscert` and `tlskey`.

## PROJECT STATUS

xlauthd is intentionally focused on a narrow problem: providing enough LDAP 
functionality for authentication while keeping the implementation extremely small.

The protocol surface should remain small by design.

New features should be considered carefully against the project's primary objective:

> **If a feature does not contribute directly to lightweight LDAP authentication, it probably does not belong in xlauthd.**

## LICENSE

See `LICENSE` for the license covering xlauthd.

xlauthd may use third-party Go dependencies distributed under licenses different 
from the project's own license. Those dependencies remain subject to their respective 
licenses and attribution requirements.

Users redistributing xlauthd should review the licenses of dependencies included in 
their build and preserve any notices, copyright statements, and license texts required 
by those licenses.

## NAME ORIGIN

**xlauthd** is intended to convey **"extremely lightweight authentication daemon"**.

The name reflects the project's purpose: a small authentication daemon providing 
just enough LDAP protocol support to integrate simple htpasswd-style credentials 
with existing LDAP-aware software.
