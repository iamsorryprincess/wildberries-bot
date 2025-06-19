#!/bin/sh

srcdir=cmd
bindir=bin
remotehost=test-spb-servak

env GOOS=linux GOARCH=arm64 go build -o $bindir/api $srcdir/api/main.go
scp $bindir/api $remotehost:api
rm $bindir/api
