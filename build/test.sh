#!/bin/bash
set -Eeuo pipefail

go test -coverprofile artifacts/cover_unittests.out \
    `go list ./... | \
    grep "central/main\|central/sys" | \
    grep -v "git.proxeus.com/core/central/sys/cache" | \
    grep -v "git.proxeus.com/core/central/sys/db/storm" | \
    grep -v "git.proxeus.com/core/central/sys/eio" | \
    grep -v "central/main/handlers/assets" | \
    grep -v "central/main$"`
