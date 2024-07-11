#!/bin/bash

git tag "$1"
git push origin "$1"
GOPROXY=proxy.golang.org go list -m "github.com/nothing2512/mailer@$1"