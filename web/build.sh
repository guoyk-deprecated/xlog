#!/bin/bash

PKG=web binfs views public > binfs.gen.go
gofmt -s -w binfs.gen.go
