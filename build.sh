#!/bin/bash

set -e
set -u

cd web/ui
yarn build

cd ..
PKG=web binfs public > binfs.gen.go
