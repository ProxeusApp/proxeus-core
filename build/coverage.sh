#!/bin/bash
set -Eeuo pipefail


go test -v -timeout=24h -coverprofile artifacts/cover.out \
      -coverpkg="$(go list \
          git.proxeus.com/core/central/main/... \
          git.proxeus.com/core/central/sys/... \
          | grep -v /assets | tr '\n' ',' | sed 's/.$//')" \
  ./main