#!/bin/bash
set -Eeuo pipefail

go test -v -coverprofile artifacts/cover.out \
      -coverpkg="$(go list \
          git.proxeus.com/core/central/main/... \
          git.proxeus.com/core/central/sys/... \
          | tr '\n' ',' | sed 's/.$//')" \
  ./main