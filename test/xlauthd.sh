#!/usr/bin/env bash

XDIR=$(cd $(dirname $0) && pwd)

exec $XDIR/xlauthd \
  -basedn 'dc=example,dc=net' \
  -listen ":5389" \
  -userdb $XDIR/user.htpasswd \
  -systemdb $XDIR/system.htpasswd \
  -tlscert $XDIR/tls.crt.pem \
  -tlskey $XDIR/tls.key.pem