#! /usr/bin/env bash
set -e

exec mcdev-each-change go test {{.Pkg}}
