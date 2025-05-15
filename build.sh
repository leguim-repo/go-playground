#!/bin/bash
for d in cmd/*/; do go build -o "bin/$(basename "$d")" "./$d"; done
