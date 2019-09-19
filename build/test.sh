#!/bin/bash
set -Eeuo pipefail

go test \
    `go list ./... | \
    grep "central/main\|central/lib" | \
    grep -v "central/main/handlers/assets" | \
    grep -v "central/main$"`
