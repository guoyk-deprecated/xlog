#!/bin/bash

PKG=web binfs views > binfs.gen.go
gofmt -s -w binfs.gen.go
