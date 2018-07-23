#!/bin/bash

set -e
set -u

cd web/ui
yarn build

cd ..
rm public/*.map
PKG=web binfs public > binfs.gen.go
